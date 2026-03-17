package auth

import (
	"go-api/internal/middleware"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", LoginHandler(db, logger))
		r.Post("/register", RegisterHandler(db, logger))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(logger))
		r.Get("/profile", ProfileHandler(db, logger))
		r.Patch("/profile", ProfileHandler(db, logger))
	})
}
