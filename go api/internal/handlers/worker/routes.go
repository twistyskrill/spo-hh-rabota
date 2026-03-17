package worker

import (
	"log/slog"
	"net/http"

	"go-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {
	r.Route("/handyman", func(r chi.Router) {
		r.Get("/", AllWorkersHandler(db, logger))
		r.Get("/{id}", WorkerHandler(db, logger))
	})

	// Маршруты для управления категориями конкретного мастера (требуют аутентификации)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(logger))
		r.Route("/handyman/categories", func(r chi.Router) {
			r.Method(http.MethodGet, "/", CategoryHandler(db, logger))
			r.Method(http.MethodPost, "/", CategoryHandler(db, logger))
			r.Method(http.MethodDelete, "/", CategoryHandler(db, logger))
		})
	})
}
