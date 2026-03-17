package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"go-api/internal/auth"
)

func AuthMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Проверяем заголовок
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Error("Authorization header required")
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// 2. Проверяем Bearer
			if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				logger.Error("Invalid auth format", "header", authHeader[:10]+"...")
				http.Error(w, "Invalid header format", http.StatusUnauthorized)
				return
			}

			// 3. Валидируем токен
			tokenString := authHeader[7:]
			claims, err := auth.ValidateToken(tokenString, logger)
			if err != nil {
				logger.Error("Invalid token", "error", err.Error())
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// 4. Добавляем в контекст и передаем дальше
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
