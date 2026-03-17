package admin

import (
	"encoding/json"
	"go-api/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// ======================================================================
// КАТЕГОРИИ
// ======================================================================

// CreateCategoryHandler - создание новой категории
func CreateCategoryHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Ошибка декодирования запроса", "error", err)
			return
		}

		if req.Name == "" {
			http.Error(w, "Category name is required", http.StatusBadRequest)
			return
		}

		category := models.Category{
			Name: req.Name,
		}

		if err := db.Create(&category).Error; err != nil {
			http.Error(w, "Failed to create category", http.StatusInternalServerError)
			logger.Error("Ошибка создания категории", "error", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Category created successfully",
			"id":      category.ID,
			"name":    category.Name,
		})

		logger.Info("Категория создана", "id", category.ID, "name", category.Name)
	}
}

// UpdateCategoryHandler - обновление категории
func UpdateCategoryHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		categoryIDStr := chi.URLParam(r, "categoryID")
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Ошибка декодирования запроса", "error", err)
			return
		}

		if req.Name == "" {
			http.Error(w, "Category name is required", http.StatusBadRequest)
			return
		}

		var category models.Category
		if err := db.First(&category, categoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Category not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			logger.Error("Ошибка поиска категории", "error", err)
			return
		}

		category.Name = req.Name

		if err := db.Save(&category).Error; err != nil {
			http.Error(w, "Failed to update category", http.StatusInternalServerError)
			logger.Error("Ошибка обновления категории", "error", err)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Category updated successfully",
			"id":      category.ID,
			"name":    category.Name,
		})

		logger.Info("Категория обновлена", "id", category.ID, "name", category.Name)
	}
}

// DeleteCategoryHandler - удаление категории
func DeleteCategoryHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		categoryIDStr := chi.URLParam(r, "categoryID")
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		var category models.Category
		if err := db.First(&category, categoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Category not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			logger.Error("Ошибка поиска категории", "error", err)
			return
		}

		// Проверяем, есть ли объявления с этой категорией
		var adsCount int64
		db.Model(&models.Ad{}).Where("category_id = ?", categoryID).Count(&adsCount)
		if adsCount > 0 {
			http.Error(w, "Cannot delete category: it is used in ads", http.StatusConflict)
			return
		}

		// Проверяем, есть ли рабочие с этой категорией
		var workerCategoriesCount int64
		db.Model(&models.WorkerCategory{}).Where("category_id = ?", categoryID).Count(&workerCategoriesCount)
		if workerCategoriesCount > 0 {
			http.Error(w, "Cannot delete category: it is used by workers", http.StatusConflict)
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			http.Error(w, "Failed to delete category", http.StatusInternalServerError)
			logger.Error("Ошибка удаления категории", "error", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Category deleted successfully",
		})

		logger.Info("Категория удалена", "id", categoryID)
	}
}

// ======================================================================
// ЕДИНИЦЫ ЦЕНЫ
// ======================================================================

// CreatePriceUnitHandler - создание новой единицы цены
func CreatePriceUnitHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Ошибка декодирования запроса", "error", err)
			return
		}

		if req.Name == "" {
			http.Error(w, "Price unit name is required", http.StatusBadRequest)
			return
		}

		priceUnit := models.PriceUnit{
			Name: req.Name,
		}

		if err := db.Create(&priceUnit).Error; err != nil {
			http.Error(w, "Failed to create price unit", http.StatusInternalServerError)
			logger.Error("Ошибка создания единицы цены", "error", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Price unit created successfully",
			"id":      priceUnit.ID,
			"name":    priceUnit.Name,
		})

		logger.Info("Единица цены создана", "id", priceUnit.ID, "name", priceUnit.Name)
	}
}

// UpdatePriceUnitHandler - обновление единицы цены
func UpdatePriceUnitHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		priceUnitIDStr := chi.URLParam(r, "priceUnitID")
		priceUnitID, err := strconv.ParseUint(priceUnitIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid price unit ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Ошибка декодирования запроса", "error", err)
			return
		}

		if req.Name == "" {
			http.Error(w, "Price unit name is required", http.StatusBadRequest)
			return
		}

		var priceUnit models.PriceUnit
		if err := db.First(&priceUnit, priceUnitID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Price unit not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			logger.Error("Ошибка поиска единицы цены", "error", err)
			return
		}

		priceUnit.Name = req.Name

		if err := db.Save(&priceUnit).Error; err != nil {
			http.Error(w, "Failed to update price unit", http.StatusInternalServerError)
			logger.Error("Ошибка обновления единицы цены", "error", err)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Price unit updated successfully",
			"id":      priceUnit.ID,
			"name":    priceUnit.Name,
		})

		logger.Info("Единица цены обновлена", "id", priceUnit.ID, "name", priceUnit.Name)
	}
}

// DeletePriceUnitHandler - удаление единицы цены
func DeletePriceUnitHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		priceUnitIDStr := chi.URLParam(r, "priceUnitID")
		priceUnitID, err := strconv.ParseUint(priceUnitIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid price unit ID", http.StatusBadRequest)
			return
		}

		var priceUnit models.PriceUnit
		if err := db.First(&priceUnit, priceUnitID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Price unit not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			logger.Error("Ошибка поиска единицы цены", "error", err)
			return
		}

		// Проверяем, есть ли объявления с этой единицей цены
		var adsCount int64
		db.Model(&models.Ad{}).Where("price_unit_id = ?", priceUnitID).Count(&adsCount)
		if adsCount > 0 {
			http.Error(w, "Cannot delete price unit: it is used in ads", http.StatusConflict)
			return
		}

		if err := db.Delete(&priceUnit).Error; err != nil {
			http.Error(w, "Failed to delete price unit", http.StatusInternalServerError)
			logger.Error("Ошибка удаления единицы цены", "error", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Price unit deleted successfully",
		})

		logger.Info("Единица цены удалена", "id", priceUnitID)
	}
}
