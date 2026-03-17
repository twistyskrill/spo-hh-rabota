package admin

import (
	"go-api/internal/middleware"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, logger *slog.Logger, r chi.Router) {
	admin := chi.NewRouter()

	// Защита: требуется аутентификация + роль администратора
	admin.Use(middleware.AuthMiddleware(logger))
	admin.Use(middleware.AdminMiddleware(db, logger))

	// Управление пользователями
	admin.Get("/users", GetUsersHandler(db, logger))                       // GET /admin/users - список пользователей
	admin.Get("/users/{userID}", GetUserHandler(db, logger))               // GET /admin/users/123 - пользователь по ID
	admin.Delete("/users/{userID}", DeleteUserHandler(db, logger))         // DELETE /admin/users/123 - удалить пользователя
	admin.Patch("/users/{userID}/role", UpdateUserRoleHandler(db, logger)) // PATCH /admin/users/123/role - изменить роль

	// Модерация объявлений
	admin.Get("/ads", GetAllAdsHandler(db, logger))                  // GET /admin/ads - все объявления (?status=pending|approved|rejected)
	admin.Delete("/ads/{adID}", DeleteAdHandler(db, logger))         // DELETE /admin/ads/123 - удалить объявление
	admin.Patch("/ads/{adID}/approve", ApproveAdHandler(db, logger)) // PATCH /admin/ads/123/approve - одобрить объявление
	admin.Patch("/ads/{adID}/reject", RejectAdHandler(db, logger))   // PATCH /admin/ads/123/reject - отклонить объявление

	// Модерация откликов
	admin.Get("/responses", GetAllResponsesHandler(db, logger))                // GET /admin/responses - все отклики
	admin.Delete("/responses/{responseID}", DeleteResponseHandler(db, logger)) // DELETE /admin/responses/123 - удалить отклик

	// Модерация профилей мастеров
	admin.Get("/workers", GetPendingWorkersHandler(db, logger))                  // GET /admin/workers - профили мастеров (?status=pending|approved|rejected)
	admin.Patch("/workers/{workerID}/approve", ApproveWorkerHandler(db, logger)) // PATCH /admin/workers/123/approve - одобрить профиль
	admin.Patch("/workers/{workerID}/reject", RejectWorkerHandler(db, logger))   // PATCH /admin/workers/123/reject - отклонить профиль

	// Статистика
	admin.Get("/stats", GetStatsHandler(db, logger)) // GET /admin/stats - общая статистика

	// Черный список
	admin.Get("/blacklist", GetBlacklistHandler(db, logger))                   // GET /admin/blacklist
	admin.Post("/blacklist", AddToBlacklistHandler(db, logger))                // POST /admin/blacklist
	admin.Delete("/blacklist/{email}", RemoveFromBlacklistHandler(db, logger)) // DELETE /admin/blacklist/email@example.com

	// Управление справочниками
	// Категории
	admin.Post("/categories", CreateCategoryHandler(db, logger))                // POST /admin/categories - создать категорию
	admin.Patch("/categories/{categoryID}", UpdateCategoryHandler(db, logger))  // PATCH /admin/categories/123 - обновить категорию
	admin.Delete("/categories/{categoryID}", DeleteCategoryHandler(db, logger)) // DELETE /admin/categories/123 - удалить категорию

	// Единицы цены
	admin.Post("/price-units", CreatePriceUnitHandler(db, logger))                 // POST /admin/price-units - создать единицу цены
	admin.Patch("/price-units/{priceUnitID}", UpdatePriceUnitHandler(db, logger))  // PATCH /admin/price-units/123 - обновить единицу цены
	admin.Delete("/price-units/{priceUnitID}", DeletePriceUnitHandler(db, logger)) // DELETE /admin/price-units/123 - удалить единицу цены

	r.Mount("/admin", admin)
}
