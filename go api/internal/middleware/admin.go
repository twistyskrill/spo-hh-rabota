package middleware

import (
	"go-api/internal/models"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

// AdminMiddleware - проверяет, что пользователь имеет роль администратора
// Должен использоваться ПОСЛЕ AuthMiddleware
func AdminMiddleware(db *gorm.DB, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем user_id из контекста (должен быть установлен в AuthMiddleware)
			userID, ok := r.Context().Value("user_id").(uint)
			if !ok {
				logger.Error("user_id not found in context")
				http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// Загружаем пользователя с ролью
			var user models.User
			if err := db.Preload("Role").First(&user, userID).Error; err != nil {
				logger.Error("failed to load user", "error", err)
				http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
				return
			}

			// Проверяем роль (например, "admin" или "administrator")
			if user.Role.RoleName != "admin" && user.Role.RoleName != "administrator" {
				logger.Warn("access denied: not an admin", "user_id", userID, "role", user.Role.RoleName)
				http.Error(w, `{"error": "access denied: admin role required"}`, http.StatusForbidden)
				return
			}

			// Пользователь - администратор, пропускаем дальше
			logger.Info("admin access granted", "user_id", userID, "email", user.Email)
			next.ServeHTTP(w, r)
		})
	}
}
