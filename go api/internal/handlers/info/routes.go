package info

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {

	r.Route("/info", func(r chi.Router) {
		r.Get("/categories", CategoriesHandler(db, logger))
		r.Get("/price_units", PriceUnitsHandler(db, logger))
	})

}
