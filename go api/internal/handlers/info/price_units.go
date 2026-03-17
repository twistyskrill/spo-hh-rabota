package info

import (
	"encoding/json"
	"go-api/internal/models"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func PriceUnitsHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var price_units []models.PriceUnit

		if err := db.Select("id, name").Find(&price_units).Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			logger.Error("Ошибка парсинга цен из бд", "error", err)
			return
		}

		type PriceUnit struct {
			Name string `json:"name"`
			ID   uint   `json:"id"`
		}

		result := make([]PriceUnit, len(price_units))
		for i, p := range price_units {
			result[i] = PriceUnit{ID: p.ID, Name: p.Name}
		}

		json.NewEncoder(w).Encode(result)
	}
}
