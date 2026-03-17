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

// GetAllAdsHandler - получить все объявления (с фильтрами)
func GetAllAdsHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
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

		type AdInfo struct {
			ID             uint      `json:"id"`
			Title          string    `json:"title"`
			Price          float64   `json:"price"`
			Location       string    `json:"location"`
			CreatedAt      time.Time `json:"created_at"`
			CategoryName   string    `json:"category_name"`
			PriceUnitName  string    `json:"price_unit_name"`
			UserID         uint      `json:"user_id"`
			UserName       string    `json:"user_name"`
			UserEmail      string    `json:"user_email"`
			ResponsesCount int64     `json:"responses_count"`
			Status         string    `json:"status"`
		}

		var ads []AdInfo
		query := db.Table("ads a").
			Select("a.id, a.title, a.price, a.location, a.created_at, a.status, " +
				"c.name as category_name, pu.name as price_unit_name, " +
				"u.id as user_id, u.name as user_name, u.email as user_email, " +
				"COUNT(r.id) as responses_count").
			Joins("JOIN categories c ON a.category_id = c.id").
			Joins("JOIN price_units pu ON a.price_unit_id = pu.id").
			Joins("JOIN users u ON a.user_id = u.id").
			Joins("LEFT JOIN responses r ON a.id = r.ad_id AND r.deleted_at IS NULL").
			Where("a.deleted_at IS NULL").
			Group("a.id, a.title, a.price, a.location, a.created_at, a.status, c.name, pu.name, u.id, u.name, u.email").
			Order("a.created_at DESC").
			Limit(limit).
			Offset(offset)

		// Фильтры
		if category := r.URL.Query().Get("category"); category != "" {
			query = query.Where("c.name ILIKE ?", "%"+category+"%")
		}
		if userID := r.URL.Query().Get("user_id"); userID != "" {
			query = query.Where("a.user_id = ?", userID)
		}
		if status := r.URL.Query().Get("status"); status != "" {
			query = query.Where("a.status = ?", status)
		}

		var total int64
		db.Model(&models.Ad{}).Count(&total)

		if err := query.Scan(&ads).Error; err != nil {
			logger.Error("failed to get ads", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"ads":    ads,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		})
	}
}

// DeleteAdHandler - удалить объявление
func DeleteAdHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		adIDStr := chi.URLParam(r, "adID")
		adID, err := strconv.ParseUint(adIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid ad id"}`, http.StatusBadRequest)
			return
		}

		var ad models.Ad
		if err := db.First(&ad, uint(adID)).Error; err != nil {
			http.Error(w, `{"error": "ad not found"}`, http.StatusNotFound)
			return
		}

		if err := db.Delete(&ad).Error; err != nil {
			logger.Error("failed to delete ad", "error", err)
			http.Error(w, `{"error": "failed to delete ad"}`, http.StatusInternalServerError)
			return
		}

		logger.Info("ad deleted by admin", "ad_id", adID)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "ad deleted successfully",
		})
	}
}

// GetAllResponsesHandler - получить все отклики (с фильтрами)
func GetAllResponsesHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
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

		type ResponseInfo struct {
			ID            uint      `json:"id"`
			AdID          uint      `json:"ad_id"`
			AdTitle       string    `json:"ad_title"`
			WorkerID      uint      `json:"worker_id"`
			WorkerName    string    `json:"worker_name"`
			WorkerEmail   string    `json:"worker_email"`
			Message       string    `json:"message"`
			ProposedPrice *float64  `json:"proposed_price"`
			Status        string    `json:"status"`
			CreatedAt     time.Time `json:"created_at"`
		}

		var responses []ResponseInfo
		query := db.Table("responses r").
			Select("r.id, r.ad_id, r.worker_id, r.message, r.proposed_price, r.status, r.created_at, " +
				"a.title as ad_title, " +
				"u.name as worker_name, u.email as worker_email").
			Joins("JOIN ads a ON r.ad_id = a.id").
			Joins("JOIN users u ON r.worker_id = u.id").
			Where("r.deleted_at IS NULL").
			Order("r.created_at DESC").
			Limit(limit).
			Offset(offset)

		// Фильтры
		if status := r.URL.Query().Get("status"); status != "" {
			query = query.Where("r.status = ?", status)
		}
		if workerID := r.URL.Query().Get("worker_id"); workerID != "" {
			query = query.Where("r.worker_id = ?", workerID)
		}

		var total int64
		db.Model(&models.Response{}).Count(&total)

		if err := query.Scan(&responses).Error; err != nil {
			logger.Error("failed to get responses", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"responses": responses,
			"total":     total,
			"limit":     limit,
			"offset":    offset,
		})
	}
}

// DeleteResponseHandler - удалить отклик
func DeleteResponseHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		responseIDStr := chi.URLParam(r, "responseID")
		responseID, err := strconv.ParseUint(responseIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid response id"}`, http.StatusBadRequest)
			return
		}

		var response models.Response
		if err := db.First(&response, uint(responseID)).Error; err != nil {
			http.Error(w, `{"error": "response not found"}`, http.StatusNotFound)
			return
		}

		if err := db.Delete(&response).Error; err != nil {
			logger.Error("failed to delete response", "error", err)
			http.Error(w, `{"error": "failed to delete response"}`, http.StatusInternalServerError)
			return
		}

		logger.Info("response deleted by admin", "response_id", responseID)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "response deleted successfully",
		})
	}
}

// GetStatsHandler - получить статистику платформы
func GetStatsHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		type Stats struct {
			TotalUsers     int64            `json:"total_users"`
			TotalAds       int64            `json:"total_ads"`
			TotalResponses int64            `json:"total_responses"`
			TotalWorkers   int64            `json:"total_workers"`
			UsersByRole    map[string]int64 `json:"users_by_role"`
			AdsByCategory  map[string]int64 `json:"ads_by_category"`
			RecentActivity struct {
				NewUsersToday     int64 `json:"new_users_today"`
				NewAdsToday       int64 `json:"new_ads_today"`
				NewResponsesToday int64 `json:"new_responses_today"`
			} `json:"recent_activity"`
		}

		var stats Stats

		// Подсчет пользователей
		db.Model(&models.User{}).Count(&stats.TotalUsers)
		db.Model(&models.Ad{}).Count(&stats.TotalAds)
		db.Model(&models.Response{}).Count(&stats.TotalResponses)
		db.Model(&models.WorkerProfile{}).Where("have_worker_profile = ?", true).Count(&stats.TotalWorkers)

		// Пользователи по ролям
		stats.UsersByRole = make(map[string]int64)
		type RoleCount struct {
			RoleName string
			Count    int64
		}
		var roleCounts []RoleCount
		db.Table("users u").
			Select("r.role_name, COUNT(u.id) as count").
			Joins("JOIN roles r ON u.role_id = r.id").
			Where("u.deleted_at IS NULL").
			Group("r.role_name").
			Scan(&roleCounts)
		for _, rc := range roleCounts {
			stats.UsersByRole[rc.RoleName] = rc.Count
		}

		// Объявления по категориям
		stats.AdsByCategory = make(map[string]int64)
		type CategoryCount struct {
			CategoryName string
			Count        int64
		}
		var catCounts []CategoryCount
		db.Table("ads a").
			Select("c.name as category_name, COUNT(a.id) as count").
			Joins("JOIN categories c ON a.category_id = c.id").
			Where("a.deleted_at IS NULL").
			Group("c.name").
			Scan(&catCounts)
		for _, cc := range catCounts {
			stats.AdsByCategory[cc.CategoryName] = cc.Count
		}

		// Активность за сегодня
		today := time.Now().Truncate(24 * time.Hour)
		db.Model(&models.User{}).Where("created_at >= ?", today).Count(&stats.RecentActivity.NewUsersToday)
		db.Model(&models.Ad{}).Where("created_at >= ?", today).Count(&stats.RecentActivity.NewAdsToday)
		db.Model(&models.Response{}).Where("created_at >= ?", today).Count(&stats.RecentActivity.NewResponsesToday)

		json.NewEncoder(w).Encode(stats)
	}
}

// GetBlacklistHandler - получить черный список
func GetBlacklistHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var blacklist []models.BlackList
		if err := db.Find(&blacklist).Error; err != nil {
			logger.Error("failed to get blacklist", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"blacklist": blacklist,
			"total":     len(blacklist),
		})
	}
}

// AddToBlacklistHandler - добавить email в черный список
func AddToBlacklistHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		type BlacklistRequest struct {
			Email string `json:"email"`
		}

		var req BlacklistRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
			return
		}

		if req.Email == "" {
			http.Error(w, `{"error": "email is required"}`, http.StatusBadRequest)
			return
		}

		blacklistEntry := models.BlackList{Email: req.Email}
		if err := db.Create(&blacklistEntry).Error; err != nil {
			logger.Error("failed to add to blacklist", "error", err)
			http.Error(w, `{"error": "failed to add to blacklist"}`, http.StatusInternalServerError)
			return
		}

		logger.Info("email added to blacklist by admin", "email", req.Email)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "email added to blacklist",
			"email":   req.Email,
		})
	}
}

// RemoveFromBlacklistHandler - удалить email из черного списка
func RemoveFromBlacklistHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		email := chi.URLParam(r, "email")
		if email == "" {
			http.Error(w, `{"error": "email is required"}`, http.StatusBadRequest)
			return
		}

		var blacklistEntry models.BlackList
		if err := db.Where("email = ?", email).First(&blacklistEntry).Error; err != nil {
			http.Error(w, `{"error": "email not found in blacklist"}`, http.StatusNotFound)
			return
		}

		if err := db.Delete(&blacklistEntry).Error; err != nil {
			logger.Error("failed to remove from blacklist", "error", err)
			http.Error(w, `{"error": "failed to remove from blacklist"}`, http.StatusInternalServerError)
			return
		}

		logger.Info("email removed from blacklist by admin", "email", email)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "email removed from blacklist",
			"email":   email,
		})
	}
}

// ======================================================================
// МОДЕРАЦИЯ — ОДОБРЕНИЕ / ОТКЛОНЕНИЕ
// ======================================================================

// ApproveAdHandler - одобрить объявление
func ApproveAdHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		adIDStr := chi.URLParam(r, "adID")
		adID, err := strconv.ParseUint(adIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid ad id"}`, http.StatusBadRequest)
			return
		}

		result := db.Model(&models.Ad{}).Where("id = ?", adID).Update("status", "approved")
		if result.Error != nil {
			logger.Error("failed to approve ad", "error", result.Error)
			http.Error(w, `{"error": "failed to approve ad"}`, http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			http.Error(w, `{"error": "ad not found"}`, http.StatusNotFound)
			return
		}

		logger.Info("ad approved by admin", "ad_id", adID)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "ad approved successfully",
			"ad_id":   adID,
			"status":  "approved",
		})
	}
}

// RejectAdHandler - отклонить объявление
func RejectAdHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		adIDStr := chi.URLParam(r, "adID")
		adID, err := strconv.ParseUint(adIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid ad id"}`, http.StatusBadRequest)
			return
		}

		result := db.Model(&models.Ad{}).Where("id = ?", adID).Update("status", "rejected")
		if result.Error != nil {
			logger.Error("failed to reject ad", "error", result.Error)
			http.Error(w, `{"error": "failed to reject ad"}`, http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			http.Error(w, `{"error": "ad not found"}`, http.StatusNotFound)
			return
		}

		logger.Info("ad rejected by admin", "ad_id", adID)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "ad rejected successfully",
			"ad_id":   adID,
			"status":  "rejected",
		})
	}
}

// ApproveWorkerHandler - одобрить профиль мастера
func ApproveWorkerHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		workerIDStr := chi.URLParam(r, "workerID")
		workerID, err := strconv.ParseUint(workerIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid worker id"}`, http.StatusBadRequest)
			return
		}

		result := db.Model(&models.WorkerProfile{}).Where("user_id = ?", workerID).Update("status", "approved")
		if result.Error != nil {
			logger.Error("failed to approve worker", "error", result.Error)
			http.Error(w, `{"error": "failed to approve worker profile"}`, http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			http.Error(w, `{"error": "worker profile not found"}`, http.StatusNotFound)
			return
		}

		logger.Info("worker profile approved by admin", "worker_id", workerID)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "worker profile approved successfully",
			"worker_id": workerID,
			"status":    "approved",
		})
	}
}

// RejectWorkerHandler - отклонить профиль мастера
func RejectWorkerHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		workerIDStr := chi.URLParam(r, "workerID")
		workerID, err := strconv.ParseUint(workerIDStr, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "invalid worker id"}`, http.StatusBadRequest)
			return
		}

		result := db.Model(&models.WorkerProfile{}).Where("user_id = ?", workerID).Update("status", "rejected")
		if result.Error != nil {
			logger.Error("failed to reject worker", "error", result.Error)
			http.Error(w, `{"error": "failed to reject worker profile"}`, http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			http.Error(w, `{"error": "worker profile not found"}`, http.StatusNotFound)
			return
		}

		logger.Info("worker profile rejected by admin", "worker_id", workerID)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "worker profile rejected successfully",
			"worker_id": workerID,
			"status":    "rejected",
		})
	}
}

// GetPendingWorkersHandler - список профилей мастеров на модерации
func GetPendingWorkersHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := r.URL.Query().Get("status")
		if status == "" {
			status = "pending"
		}

		limit := 10
		offset := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			limit, _ = strconv.Atoi(l)
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			offset, _ = strconv.Atoi(o)
		}

		type WorkerInfo struct {
			UserID      uint    `json:"user_id"`
			Name        string  `json:"name"`
			Email       string  `json:"email"`
			Phone       string  `json:"phone"`
			ExpYears    *int    `json:"exp_years"`
			Description *string `json:"description"`
			Location    string  `json:"location"`
			Schedule    string  `json:"schedule"`
			Status      string  `json:"status"`
		}

		var workers []WorkerInfo
		err := db.Table("users u").
			Select("u.id as user_id, u.name, u.email, u.phone, wp.exp_years, wp.description, wp.location, wp.schedule, wp.status").
			Joins("JOIN worker_profiles wp ON u.id = wp.user_id").
			Where("wp.have_worker_profile = ? AND wp.status = ?", true, status).
			Order("wp.updated_at DESC").
			Limit(limit).
			Offset(offset).
			Scan(&workers).Error
		if err != nil {
			logger.Error("failed to get workers for moderation", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		var total int64
		db.Model(&models.WorkerProfile{}).Where("have_worker_profile = ? AND status = ?", true, status).Count(&total)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"workers": workers,
			"total":   total,
			"limit":   limit,
			"offset":  offset,
			"status":  status,
		})
	}
}
