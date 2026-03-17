package worker

import (
	"encoding/json"
	"fmt"
	"go-api/internal/models"
	"go-api/internal/storage"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func CategoryHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		worker, err := storage.WorkerByUserID(db, userID)
		if err != nil || worker == nil || worker.ID == 0 {
			http.Error(w, `{"error": "worker not found"}`, http.StatusForbidden)
			return
		}

		switch r.Method {
		case http.MethodGet:
			getMyCategories(w, worker)
		case http.MethodPost:
			addMyCategoriesByName(db, logger, w, r, worker.ID) // worker.ID == user_id == worker_id
		case http.MethodDelete:
			deleteMyCategoriesByName(db, logger, w, r, worker.ID)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func getMyCategories(w http.ResponseWriter, worker *storage.WorkerResponse) {
	json.NewEncoder(w).Encode(worker.Categories)
}

func addMyCategoriesByName(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, workerID uint) {
	type CategoryReq struct {
		CategoryNames []string `json:"category_names"`
	}
	var req CategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Гарантируем наличие WorkerProfile (для старых пользователей)
	var wp models.WorkerProfile
	if err := tx.FirstOrCreate(&wp, models.WorkerProfile{UserID: workerID}).Error; err != nil {
		tx.Rollback()
		logger.Error("failed to ensure worker profile", "error", err)
		http.Error(w, `{"error": "failed to add category"}`, http.StatusInternalServerError)
		return
	}

	for _, name := range req.CategoryNames {
		var category models.Category
		if err := tx.Where("name ILIKE ?", name).First(&category).Error; err != nil {
			tx.Rollback()
			logger.Error("category not found", "name", name)
			http.Error(w, fmt.Sprintf(`{"error": "category '%s' not found"}`, name), http.StatusBadRequest)
			return
		}

		wc := models.WorkerCategory{WorkerID: workerID, CategoryID: category.ID}
		if err := tx.Where("worker_id = ? AND category_id = ?", workerID, category.ID).
			FirstOrCreate(&wc).Error; err != nil {
			tx.Rollback()
			logger.Error("failed to add category", "error", err)
			http.Error(w, `{"error": "failed to add category"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("failed to commit categories add tx", "error", err)
		http.Error(w, `{"error": "failed to add category"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем обновлённый список категорий
	updatedWorker, err := storage.WorkerByUserID(db, workerID)
	if err != nil || updatedWorker == nil {
		http.Error(w, `{"error": "failed to load categories"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "categories added",
		"categories": updatedWorker.Categories,
	})
}

func deleteMyCategoriesByName(db *gorm.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request, workerID uint) {
	type CategoryReq struct {
		CategoryNames []string `json:"category_names"`
	}
	var req CategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, name := range req.CategoryNames {
		var category models.Category
		if err := tx.Where("name ILIKE ?", name).First(&category).Error; err != nil {
			tx.Rollback()
			http.Error(w, fmt.Sprintf(`{"error": "category '%s' not found"}`, name), http.StatusBadRequest)
			return
		}

		result := tx.Where("worker_id = ? AND category_id = ?", workerID, category.ID).
			Delete(&models.WorkerCategory{})
		if result.Error != nil {
			tx.Rollback()
			logger.Error("failed to delete category", "error", result.Error)
			http.Error(w, `{"error": "failed to delete category"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("failed to commit categories delete tx", "error", err)
		http.Error(w, `{"error": "failed to delete category"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем обновлённый список категорий
	updatedWorker, err := storage.WorkerByUserID(db, workerID)
	if err != nil || updatedWorker == nil {
		http.Error(w, `{"error": "failed to load categories"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "categories deleted",
		"categories": updatedWorker.Categories,
	})
}
