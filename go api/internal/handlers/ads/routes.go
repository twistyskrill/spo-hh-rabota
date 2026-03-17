package ads

import (
	"go-api/internal/middleware"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {
	public := chi.NewRouter()
	protected := chi.NewRouter()
	master := chi.NewRouter()

	//  ПУБЛИЧНЫЕ (мастера смотрят без токена)
	public.Get("/", PublicAdsHandler(db, logger))       // GET /ads - список всех
	public.Get("/{adID}", PublicAdsHandler(db, logger)) // GET /ads/123 - конкретное объявление

	//  ЗАЩИЩЁННЫЕ (клиент управляет своими объявлениями)
	protected.Use(middleware.AuthMiddleware(logger))
	protected.Get("/", ProtectedAdsHandler(db, logger))          // GET /my-ads - мои объявления
	protected.Get("/{adID}", ProtectedAdsHandler(db, logger))    // GET /my-ads/123 - моё объявление
	protected.Post("/", ProtectedAdsHandler(db, logger))         // POST /my-ads - создать
	protected.Patch("/{adID}", ProtectedAdsHandler(db, logger))  // PATCH /my-ads/123 - обновить
	protected.Delete("/{adID}", ProtectedAdsHandler(db, logger)) // DELETE /my-ads/123 - удалить

	// МАСТЕРА (управление откликами)
	master.Use(middleware.AuthMiddleware(logger))
	master.Get("/", MasterResponsesHandler(db, logger))                // GET /responses - мои отклики
	master.Post("/", MasterResponsesHandler(db, logger))               // POST /responses - создать отклик
	master.Delete("/{responseID}", MasterResponsesHandler(db, logger)) // DELETE /responses/123 - удалить отклик

	r.Mount("/ads", public)       // /ads → публичные объявления (для всех)
	r.Mount("/my-ads", protected) // /my-ads → личный кабинет клиента
	r.Mount("/responses", master) // /responses → отклики мастера
}
