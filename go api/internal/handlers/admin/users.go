package admin

import (
	"encoding/json"
	"go-api/internal/models"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// GetUsersHandler - получить список всех пользователей
func GetUsersHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		limit := 10
		offset := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			limit, _ = strconv.Atoi(l)
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			offset, _ = strconv.Atoi(o)
		}

		type UserInfo struct {
			ID                uint      `json:"id"`
			Email             string    `json:"email"`
			Name              string    `json:"name"`
			Phone             string    `json:"phone"`
			RoleID            uint      `json:"role_id"`
			RoleName          string    `json:"role_name"`
			CreatedAt         time.Time `json:"created_at"`
			HaveWorkerProfile bool      `json:"have_worker_profile"`
			AdsCount          int64     `json:"ads_count"`
			ResponsesCount    int64     `json:"responses_count"`
		}

		var users []UserInfo
		query := db.Table("users u").
			Select("u.id, u.email, u.name, u.phone, u.role_id, u.created_at, " +
				"r.role_name, " +
				"COALESCE(wp.have_worker_profile, false) as have_worker_profile, " +
				"COUNT(DISTINCT a.id) as ads_count, " +
				"COUNT(DISTINCT resp.id) as responses_count").
			Joins("JOIN roles r ON u.role_id = r.id").
			Joins("LEFT JOIN worker_profiles wp ON u.id = wp.user_id").
			Joins("LEFT JOIN ads a ON u.id = a.user_id AND a.deleted_at IS NULL").
			Joins("LEFT JOIN responses resp ON u.id = resp.worker_id AND resp.deleted_at IS NULL").
			Where("u.deleted_at IS NULL").
			Group("u.id, u.email, u.name, u.phone, u.role_id, u.created_at, r.role_name, wp.have_worker_profile").
			Order("u.created_at DESC").
			Limit(limit).
			Offset(offset)

		// Фильтры
		if role := r.URL.Query().Get("role"); role != "" {
			query = query.Where("r.role_name = ?", role)
		}
		if search := r.URL.Query().Get("search"); search != "" {
			query = query.Where("u.email ILIKE ? OR u.name ILIKE ?", "%"+search+"%", "%"+search+"%")
		}

		var total int64
		db.Model(&models.User{}).Count(&total)

		if err := query.Scan(&users).Error; err != nil {
			logger.Error("failed to get users", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"users":  users,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		})
	}
}

// GetUserHandler - получить пользователя по ID
func GetUserHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userIDStr := chi.URLParam(r, "userID")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Preload("Role").
			Preload("WorkerProfile").
			Preload("Ads").
			First(&user, uint(userID)).Error; err != nil {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
			return
		}

		// Подсчет откликов
		var responsesCount int64
		db.Model(&models.Response{}).Where("worker_id = ?", userID).Count(&responsesCount)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":              user.ID,
			"email":           user.Email,
			"name":            user.Name,
			"phone":           user.Phone,
			"role":            user.Role.RoleName,
			"role_id":         user.RoleID,
			"created_at":      user.CreatedAt,
			"worker_profile":  user.WorkerProfile,
			"ads_count":       len(user.Ads),
			"responses_count": responsesCount,
		})
	}
}

// DeleteUserHandler - удалить пользователя (soft delete)
func DeleteUserHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userIDStr := chi.URLParam(r, "userID")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
			return
		}

		// Проверяем, что пользователь существует
		var user models.User
		if err := db.First(&user, uint(userID)).Error; err != nil {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
			return
		}

		// Мягкое удаление
		if err := db.Delete(&user).Error; err != nil {
			logger.Error("failed to delete user", "error", err)
			http.Error(w, `{"error": "failed to delete user"}`, http.StatusInternalServerError)
			return
		}

		logger.Info("user deleted by admin", "user_id", userID)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user deleted successfully",
		})
	}
}

// UpdateUserRoleHandler - изменить роль пользователя
func UpdateUserRoleHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userIDStr := chi.URLParam(r, "userID")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
			return
		}

		type RoleUpdateRequest struct {
			RoleName string `json:"role_name"`
		}

		var req RoleUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
			return
		}

		if req.RoleName == "" {
			http.Error(w, `{"error": "role_name is required"}`, http.StatusBadRequest)
			return
		}

		// Находим роль по имени
		var role models.Role
		if err := db.Where("role_name = ?", req.RoleName).First(&role).Error; err != nil {
			http.Error(w, `{"error": "role not found"}`, http.StatusNotFound)
			return
		}

		// Проверяем, что пользователь существует
		var user models.User
		if err := db.First(&user, uint(userID)).Error; err != nil {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
			return
		}

		// Обновляем роль
		if err := db.Model(&user).Update("role_id", role.ID).Error; err != nil {
			logger.Error("failed to update user role", "error", err)
			http.Error(w, `{"error": "failed to update role"}`, http.StatusInternalServerError)
			return
		}

		logger.Info("user role updated by admin", "user_id", userID, "new_role", req.RoleName)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "role updated successfully",
			"user_id": userID,
			"role":    req.RoleName,
		})
	}
}
