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

// PublicAdsHandler - публичный доступ (только GET)
func PublicAdsHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			getAdsPublic(db, logger, w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// ProtectedAdsHandler - защищённый доступ (POST/PATCH/DELETE/GET для личных объявлений)
func ProtectedAdsHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		switch r.Method {
		case http.MethodGet:
			getAdsProtected(db, logger, w, r, userID)
		case http.MethodPost:
			createAd(db, logger, w, r, userID)
		case http.MethodPatch:
			updateAd(db, logger, w, r, userID)
		case http.MethodDelete:
			deleteAd(db, logger, w, r, userID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// GET для публичного доступа
func getAdsPublic(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	adIDStr := chi.URLParam(r, "adID")

	if adIDStr != "" {
		id, _ := strconv.ParseUint(adIDStr, 10, 32)
		getAdByIDPublic(db, logger, w, uint(id))
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
	getAdsList(db, logger, w, r, limit, offset)
}

// GET для защищённого доступа (личные объявления)
func getAdsProtected(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, userID uint) {
	adIDStr := chi.URLParam(r, "adID")

	if adIDStr != "" {
		id, _ := strconv.ParseUint(adIDStr, 10, 32)
		getAdByID(db, logger, w, uint(id), userID)
		return
	}
	// Список личных объявлений пользователя
	limit := 10
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}
	getMyAdsList(db, logger, w, userID, limit, offset)
}

// Вспомогательные GET функции

// Объявление по id (публичное)
func getAdByIDPublic(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, adID uint) {
	var ad models.Ad

	if err := db.Preload("Category").Preload("PriceUnit").Preload("User").
		Where("id = ? AND status = ?", adID, "approved").
		First(&ad).Error; err != nil {
		http.Error(w, `{"error": "ad not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ad)
}

// Объявление по id (для владельца)
func getAdByID(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, adID uint, ownerID uint) {
	var ad models.Ad

	if err := db.Preload("Category").Preload("PriceUnit").Preload("User").
		Where("id = ? AND user_id = ?", adID, ownerID).
		First(&ad).Error; err != nil {
		http.Error(w, `{"error": "ad not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ad)
}

// Список объявлений (публичный)
func getAdsList(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, limit, offset int) {
	type AdList struct {
		ID            uint      `json:"id"`
		Title         string    `json:"title"`
		Price         float64   `json:"price"`
		Location      string    `json:"location"`
		Schedule      string    `json:"schedule"`
		CreatedAt     time.Time `json:"created_at"`
		CategoryID    uint      `json:"category_id"`
		CategoryName  string    `json:"category_name"`
		PriceUnitID   uint      `json:"price_unit_id"`
		PriceUnitName string    `json:"price_unit_name"`
		UserID        uint      `json:"user_id"`
		UserName      string    `json:"user_name"`
		UserPhone     string    `json:"user_phone"`
	}

	var ads []AdList
	query := db.Table("ads a").
		Select("a.id, a.title, a.price, a.location, a.schedule, a.created_at, "+
			"c.id as category_id, c.name as category_name, "+
			"pu.id as price_unit_id, pu.name as price_unit_name, "+
			"u.id as user_id, u.name as user_name, u.phone as user_phone").
		Joins("JOIN categories c ON a.category_id = c.id").
		Joins("JOIN price_units pu ON a.price_unit_id = pu.id").
		Joins("JOIN users u ON a.user_id = u.id").
		Where("a.status = ?", "approved").
		Order("a.created_at DESC").
		Limit(limit).
		Offset(offset)

	// Фильтры
	if category := r.URL.Query().Get("category"); category != "" {
		query = query.Where("c.name ILIKE ?", "%"+category+"%")
	}
	if location := r.URL.Query().Get("location"); location != "" {
		query = query.Where("a.location ILIKE ?", "%"+location+"%")
	}

	var total int64
	db.Model(&models.Ad{}).Where("status = ?", "approved").Count(&total)

	if err := query.Scan(&ads).Error; err != nil {
		logger.Error("failed to get ads list", "error", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	type Response struct {
		Ads    []AdList `json:"ads"`
		Total  int64    `json:"total"`
		Limit  int      `json:"limit"`
		Offset int      `json:"offset"`
	}

	json.NewEncoder(w).Encode(Response{
		Ads:    ads,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// Список личных объявлений пользователя
func getMyAdsList(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, userID uint, limit, offset int) {
	type AdList struct {
		ID            uint      `json:"id"`
		Title         string    `json:"title"`
		Price         float64   `json:"price"`
		Location      string    `json:"location"`
		Schedule      string    `json:"schedule"`
		CreatedAt     time.Time `json:"created_at"`
		CategoryID    uint      `json:"category_id"`
		CategoryName  string    `json:"category_name"`
		PriceUnitID   uint      `json:"price_unit_id"`
		PriceUnitName string    `json:"price_unit_name"`
		Status        string    `json:"status"` // владелец видит статус модерации
	}

	var ads []AdList
	query := db.Table("ads a").
		Select("a.id, a.title, a.price, a.location, a.schedule, a.created_at, "+
			"c.id as category_id, c.name as category_name, "+
			"pu.id as price_unit_id, pu.name as price_unit_name, "+
			"a.status").
		Joins("JOIN categories c ON a.category_id = c.id").
		Joins("JOIN price_units pu ON a.price_unit_id = pu.id").
		Where("a.user_id = ?", userID).
		Order("a.created_at DESC").
		Limit(limit).
		Offset(offset)

	var total int64
	db.Model(&models.Ad{}).Where("user_id = ?", userID).Count(&total)

	if err := query.Scan(&ads).Error; err != nil {
		logger.Error("failed to get my ads list", "error", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	type Response struct {
		Ads    []AdList `json:"ads"`
		Total  int64    `json:"total"`
		Limit  int      `json:"limit"`
		Offset int      `json:"offset"`
	}

	json.NewEncoder(w).Encode(Response{
		Ads:    ads,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func createAd(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, userID uint) {
	type CreateAdRequest struct {
		Title       string  `json:"title"`
		Price       float64 `json:"price"`
		CategoryID  uint    `json:"category_id"`
		PriceUnitID uint    `json:"price_unit_id"`
		Location    string  `json:"location"`
		Schedule    string  `json:"schedule"`
	}

	var req CreateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode request", "error", err)
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Title == "" || req.Price <= 0 || req.CategoryID == 0 || req.PriceUnitID == 0 {
		http.Error(w, `{"error": "title, price, category_id and price_unit_id are required"}`, http.StatusBadRequest)
		return
	}

	// Проверяем существование категории и единицы измерения
	var category models.Category
	if err := db.First(&category, req.CategoryID).Error; err != nil {
		http.Error(w, `{"error": "category not found"}`, http.StatusNotFound)
		return
	}

	var priceUnit models.PriceUnit
	if err := db.First(&priceUnit, req.PriceUnitID).Error; err != nil {
		http.Error(w, `{"error": "price unit not found"}`, http.StatusNotFound)
		return
	}

	ad := models.Ad{
		Title:       req.Title,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		PriceUnitID: req.PriceUnitID,
		UserID:      userID,
		Location:    req.Location,
		Schedule:    req.Schedule,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&ad).Error; err != nil {
		logger.Error("failed to create ad", "error", err)
		http.Error(w, `{"error": "failed to create ad"}`, http.StatusInternalServerError)
		return
	}

	// Загружаем связанные данные для ответа
	db.Preload("Category").Preload("PriceUnit").Preload("User").First(&ad, ad.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ad)
}

func updateAd(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, userID uint) {
	adIDStr := chi.URLParam(r, "adID")
	if adIDStr == "" {
		http.Error(w, `{"error": "ad id is required"}`, http.StatusBadRequest)
		return
	}

	adID, err := strconv.ParseUint(adIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error": "invalid ad id"}`, http.StatusBadRequest)
		return
	}

	type UpdateAdRequest struct {
		Title       *string  `json:"title"`
		Price       *float64 `json:"price"`
		CategoryID  *uint    `json:"category_id"`
		PriceUnitID *uint    `json:"price_unit_id"`
		Location    *string  `json:"location"`
		Schedule    *string  `json:"schedule"`
	}

	var req UpdateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode request", "error", err)
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Находим объявление и проверяем владельца
	var ad models.Ad
	if err := db.Where("id = ? AND user_id = ?", uint(adID), userID).First(&ad).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error": "ad not found or access denied"}`, http.StatusNotFound)
		} else {
			logger.Error("failed to find ad", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Обновляем только переданные поля
	updates := make(map[string]interface{})
	if req.Title != nil {
		if *req.Title == "" {
			http.Error(w, `{"error": "title cannot be empty"}`, http.StatusBadRequest)
			return
		}
		updates["title"] = *req.Title
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			http.Error(w, `{"error": "price must be greater than 0"}`, http.StatusBadRequest)
			return
		}
		updates["price"] = *req.Price
	}
	if req.CategoryID != nil {
		var category models.Category
		if err := db.First(&category, *req.CategoryID).Error; err != nil {
			http.Error(w, `{"error": "category not found"}`, http.StatusNotFound)
			return
		}
		updates["category_id"] = *req.CategoryID
	}
	if req.PriceUnitID != nil {
		var priceUnit models.PriceUnit
		if err := db.First(&priceUnit, *req.PriceUnitID).Error; err != nil {
			http.Error(w, `{"error": "price unit not found"}`, http.StatusNotFound)
			return
		}
		updates["price_unit_id"] = *req.PriceUnitID
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Schedule != nil {
		updates["schedule"] = *req.Schedule
	}

	if len(updates) == 0 {
		http.Error(w, `{"error": "no fields to update"}`, http.StatusBadRequest)
		return
	}

	if err := db.Model(&ad).Updates(updates).Error; err != nil {
		logger.Error("failed to update ad", "error", err)
		http.Error(w, `{"error": "failed to update ad"}`, http.StatusInternalServerError)
		return
	}

	// Загружаем обновленное объявление со связанными данными
	db.Preload("Category").Preload("PriceUnit").Preload("User").First(&ad, ad.ID)

	json.NewEncoder(w).Encode(ad)
}

func deleteAd(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, userID uint) {
	adIDStr := chi.URLParam(r, "adID")
	if adIDStr == "" {
		http.Error(w, `{"error": "ad id is required"}`, http.StatusBadRequest)
		return
	}

	adID, err := strconv.ParseUint(adIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error": "invalid ad id"}`, http.StatusBadRequest)
		return
	}

	// Находим объявление и проверяем владельца
	var ad models.Ad
	if err := db.Where("id = ? AND user_id = ?", uint(adID), userID).First(&ad).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error": "ad not found or access denied"}`, http.StatusNotFound)
		} else {
			logger.Error("failed to find ad", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Мягкое удаление (GORM автоматически использует soft delete)
	if err := db.Delete(&ad).Error; err != nil {
		logger.Error("failed to delete ad", "error", err)
		http.Error(w, `{"error": "failed to delete ad"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "ad deleted successfully"})
}
