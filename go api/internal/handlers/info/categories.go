package info

import (
	"encoding/json"
	"go-api/internal/models"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func CategoriesHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var categories []models.Category

		if err := db.Select("id, name").Find(&categories).Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			logger.Error("Ошибка парсинга категорий из бд", "error", err)
			return
		}

		type CategoryName struct {
			Name string `json:"name"`
			ID   uint   `json:"id"`
		}

		names := make([]CategoryName, len(categories))
		for i, cat := range categories {
			names[i] = CategoryName{ID: cat.ID, Name: cat.Name}
		}
		json.NewEncoder(w).Encode(names)
	}
}
