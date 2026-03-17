package auth

import (
	"encoding/json"
	"go-api/internal/auth"
	"go-api/internal/models"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Error("Неправильный метод",
				"method", r.Method,
				"path", r.URL.Path)
			http.Error(w, "Method not allowes", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Email       string  `json:"email"`
			Name        string  `json:"name"`
			Password    string  `json:"password"`
			RoleID      uint    `json:"role"`
			Description *string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Ошибка парсинга JSON",
				"req body", r.Body)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if len(input.Email) == 0 || len(input.Name) == 0 {
			logger.Error("Пустые поля ввода",
				"email", input.Email,
				"name", input.Name,
			)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
		}

		if userExists(db, input.Email) {
			logger.Info("Пользователь с таким email уже существует", "email", input.Email)
			http.Error(w, "User already exist", http.StatusConflict)
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		passwordHash, err := auth.HashPassword(input.Password)

		if err != nil {
			tx.Rollback()
			logger.Error("Ошибка генерации хэша", "err:", err)
			http.Error(w, "Ошибка сервера при создании пароля", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Name:         input.Name,
			Email:        input.Email,
			PasswordHash: passwordHash,
			RoleID:       input.RoleID,
		}

		// создаём пользователя в рамках транзакции
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			logger.Error("Ошибка вставки пользователя", "err:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// и сразу создаём связанный WorkerProfile для любого пользователя (user_id == worker_id)
		haveProfile := input.RoleID == 2
		workerProfile := models.WorkerProfile{
			UserID:            user.ID,
			IsBusy:            false,
			HaveWorkerProfile: haveProfile,
			Description:       input.Description,
		}
		if haveProfile {
			workerProfile.Status = "approved" // автоматом одобряем для примера
		}
		if err := tx.Create(&workerProfile).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to create worker profile", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			logger.Error("Ошибка транзакции", "err", err)
			http.Error(w, "Transaction failed", http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Email, logger)
		if err != nil {
			http.Error(w, "Ошибка токена", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"role":    user.RoleID,
			"message": "Пользователь зарегистрирован",
			"token":   token,
		})
	}
}

func userExists(db *gorm.DB, email string) bool {
	var user models.User
	result := db.Where("email = ?", email).First(&user)

	return result.Error == nil
}
