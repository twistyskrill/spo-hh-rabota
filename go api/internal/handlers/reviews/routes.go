package reviews

import (
	"go-api/internal/middleware"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {
	// Public: list reviews for a worker
	r.Get("/handyman/{id}/reviews", GetWorkerReviewsHandler(db, logger))

	// Protected: create review
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(logger))
		r.Post("/reviews", CreateReviewHandler(db, logger))
	})
}
