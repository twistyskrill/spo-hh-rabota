package ads

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

// MasterResponsesHandler - управление откликами мастера
func MasterResponsesHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodGet:
			getMyResponses(db, logger, w, r, userID)
		case http.MethodPost:
			createResponse(db, logger, w, r, userID)
		case http.MethodDelete:
			deleteResponse(db, logger, w, r, userID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// createResponse - создать отклик на объявление
func createResponse(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, userID uint) {
	// Проверяем, что у пользователя есть профиль мастера
	var workerProfile models.WorkerProfile
	if err := db.Where("user_id = ? AND have_worker_profile = ?", userID, true).First(&workerProfile).Error; err != nil {
		http.Error(w, `{"error": "worker profile not found"}`, http.StatusForbidden)
		return
	}

	type CreateResponseRequest struct {
		AdID          uint     `json:"ad_id"`
		Message       string   `json:"message"`
		ProposedPrice *float64 `json:"proposed_price"`
	}

	var req CreateResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode request", "error", err)
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.AdID == 0 {
		http.Error(w, `{"error": "ad_id is required"}`, http.StatusBadRequest)
		return
	}

	// Проверяем существование объявления
	var ad models.Ad
	if err := db.Preload("Category").First(&ad, req.AdID).Error; err != nil {
		http.Error(w, `{"error": "ad not found"}`, http.StatusNotFound)
		return
	}

	// Проверяем, что объявление не принадлежит мастеру
	if ad.UserID == userID {
		http.Error(w, `{"error": "cannot respond to own ad"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, что категория объявления входит в категории мастера
	var workerCategory models.WorkerCategory
	if err := db.Where("worker_id = ? AND category_id = ?", userID, ad.CategoryID).First(&workerCategory).Error; err != nil {
		http.Error(w, `{"error": "ad category does not match worker categories"}`, http.StatusForbidden)
		return
	}

	// Проверяем, что мастер еще не откликнулся на это объявление
	var existingResponse models.Response
	if err := db.Where("ad_id = ? AND worker_id = ?", req.AdID, userID).First(&existingResponse).Error; err == nil {
		http.Error(w, `{"error": "response already exists"}`, http.StatusConflict)
		return
	}

	// Создаем отклик
	response := models.Response{
		AdID:          req.AdID,
		WorkerID:      userID,
		Message:       req.Message,
		ProposedPrice: req.ProposedPrice,
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	if err := db.Create(&response).Error; err != nil {
		logger.Error("failed to create response", "error", err)
		http.Error(w, `{"error": "failed to create response"}`, http.StatusInternalServerError)
		return
	}

	// Загружаем связанные данные для ответа
	db.Preload("Ad").Preload("Ad.Category").Preload("Ad.User").First(&response, response.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// getMyResponses - получить список откликов мастера
func getMyResponses(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, userID uint) {
	// Проверяем, что у пользователя есть профиль мастера
	var workerProfile models.WorkerProfile
	if err := db.Where("user_id = ? AND have_worker_profile = ?", userID, true).First(&workerProfile).Error; err != nil {
		http.Error(w, `{"error": "worker profile not found"}`, http.StatusForbidden)
		return
	}

	limit := 10
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	type ResponseList struct {
		ID            uint      `json:"id"`
		AdID          uint      `json:"ad_id"`
		AdTitle       string    `json:"ad_title"`
		Message       string    `json:"message"`
		ProposedPrice *float64  `json:"proposed_price"`
		Status        string    `json:"status"`
		CreatedAt     time.Time `json:"created_at"`
		ClientName    string    `json:"client_name"`
		ClientPhone   string    `json:"client_phone"`
	}

	var responses []ResponseList
	query := db.Table("responses r").
		Select("r.id, r.ad_id, r.message, r.proposed_price, r.status, r.created_at, "+
			"a.title as ad_title, "+
			"u.name as client_name, u.phone as client_phone").
		Joins("JOIN ads a ON r.ad_id = a.id AND a.deleted_at IS NULL").
		Joins("JOIN users u ON a.user_id = u.id AND u.deleted_at IS NULL").
		Where("r.worker_id = ? AND r.deleted_at IS NULL", userID).
		Order("r.created_at DESC").
		Limit(limit).
		Offset(offset)

	// Фильтр по статусу
	if status := r.URL.Query().Get("status"); status != "" {
		query = query.Where("r.status = ?", status)
	}

	var total int64
	db.Model(&models.Response{}).Where("worker_id = ?", userID).Count(&total)

	if err := query.Scan(&responses).Error; err != nil {
		logger.Error("failed to get my responses", "error", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	type Response struct {
		Responses []ResponseList `json:"responses"`
		Total     int64          `json:"total"`
		Limit     int            `json:"limit"`
		Offset    int            `json:"offset"`
	}

	json.NewEncoder(w).Encode(Response{
		Responses: responses,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// deleteResponse - удалить (отменить) отклик
func deleteResponse(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, userID uint) {
	responseIDStr := chi.URLParam(r, "responseID")
	if responseIDStr == "" {
		http.Error(w, `{"error": "response id is required"}`, http.StatusBadRequest)
		return
	}

	responseID, err := strconv.ParseUint(responseIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error": "invalid response id"}`, http.StatusBadRequest)
		return
	}

	// Находим отклик и проверяем владельца
	var response models.Response
	if err := db.Where("id = ? AND worker_id = ?", uint(responseID), userID).First(&response).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error": "response not found or access denied"}`, http.StatusNotFound)
		} else {
			logger.Error("failed to find response", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Мягкое удаление
	if err := db.Delete(&response).Error; err != nil {
		logger.Error("failed to delete response", "error", err)
		http.Error(w, `{"error": "failed to delete response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "response deleted successfully"})
}
