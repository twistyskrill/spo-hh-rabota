package auth

import (
	"encoding/json"
	"errors"
	"go-api/internal/auth"

	// "go-api/internal/models"
	"go-api/internal/storage"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func LoginHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Error("Неправильный метод",
				"method", r.Method,
				"path", r.URL.Path)
			http.Error(w, "Method not allowes", http.StatusMethodNotAllowed)
			return
		}

		type LoginInput struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var input LoginInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Ошибка парсинга JSON",
				"req body", r.Body)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		user, err := storage.UserByEmail(db, input.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
			logger.Debug("Неверный пароль", "user", user)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Email, logger)

		if err != nil {
			http.Error(w, "Token generation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":    user.ID,
				"email": user.Email,
				"role":  user.Role.RoleName,
			},
		})
	}
}
