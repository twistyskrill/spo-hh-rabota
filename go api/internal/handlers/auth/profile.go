package auth

import (
	"encoding/json"
	"errors"
	"go-api/internal/models"
	"go-api/internal/storage"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func ProfileHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProfile(db, logger)(w, r)
		case http.MethodPatch:
			editProfile(db, logger)(w, r)
		default:
			http.Error(w, "Method not allows", http.StatusMethodNotAllowed)
			logger.Error("Ошибка метода в хендлера роутера", "Метод", r.Method)
			return
		}
	}

}

func getProfile(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		user, err := storage.UserById(db, userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Загружаем Worker-представление через единый storage-слой (с категориями)
		workerResp, _ := storage.WorkerByUserID(db, userID)

		response := map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role.RoleName,
			"name":  user.Name,
			"phone": user.Phone,
		}

		if workerResp != nil && workerResp.HaveWorkerProfile {
			response["have_worker_profile"] = true
			response["worker"] = map[string]interface{}{
				"specialization": workerResp.Categories,
				"experience":     workerResp.ExpYears,
				"hourly_rate":    workerResp.HourlyRate,
				"description":    workerResp.Description,
				"is_busy":        workerResp.IsBusy,
				"location":       workerResp.Location,
				"schedule":       workerResp.Schedule,
			}
		} else {
			response["have_worker_profile"] = false
		}

		json.NewEncoder(w).Encode(response)
	}
}

// самый важный момент - решить проблему, если сначала регаешься как юзер, а потом как рабочий
func editProfile(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		type ProfileInput struct {
			// User поля
			Name  *string `json:"name,omitempty"`
			Phone *string `json:"phone,omitempty"`

			// WorkerProfile поля
			ExpYears    *int    `json:"exp_years,omitempty"`
			HourlyRate  *int    `json:"hourly_rate,omitempty"`
			Description *string `json:"description,omitempty"`
			IsBusy      *bool   `json:"is_busy,omitempty"`
			Location    *string `json:"location,omitempty"`
			Schedule    *string `json:"schedule,omitempty"`

			// Категории по НАЗВАНИЯМ (полная замена списка)
			CategoryNames []string `json:"category_names,omitempty"`
		}

		var input ProfileInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		userUpdates := map[string]interface{}{}
		if input.Name != nil {
			userUpdates["name"] = *input.Name
		}
		if input.Phone != nil {
			userUpdates["phone"] = *input.Phone
		}

		if len(userUpdates) > 0 {
			result := tx.Model(&models.User{}).Where("id = ?", userID).Updates(userUpdates)
			if result.Error != nil {
				tx.Rollback()
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			if result.RowsAffected == 0 {
				tx.Rollback()
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
		}

		var user models.User
		if err := tx.Preload("Role").First(&user, userID).Error; err != nil {
			tx.Rollback()
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		workerUpdates := map[string]interface{}{}
		if input.ExpYears != nil {
			workerUpdates["exp_years"] = *input.ExpYears
		}
		if input.HourlyRate != nil {
			workerUpdates["hourly_rate"] = *input.HourlyRate
		}
		if input.Description != nil {
			workerUpdates["description"] = *input.Description
		}
		if input.IsBusy != nil {
			workerUpdates["is_busy"] = *input.IsBusy
		}
		if input.Location != nil {
			workerUpdates["location"] = *input.Location
		}
		if input.Schedule != nil {
			workerUpdates["schedule"] = *input.Schedule
		}

		// Обновление категорий работника (полная замена списка по НАЗВАНИЯМ)
		if len(input.CategoryNames) > 0 {
			// Гарантируем наличие WorkerProfile (для старых пользователей)
			var wp models.WorkerProfile
			if err := tx.FirstOrCreate(&wp, models.WorkerProfile{UserID: userID}).Error; err != nil {
				tx.Rollback()
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			// Сначала чистим все старые категории
			if err := tx.Where("worker_id = ?", userID).Delete(&models.WorkerCategory{}).Error; err != nil {
				tx.Rollback()
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			// Вставляем новые связи по именам категорий
			for _, name := range input.CategoryNames {
				var category models.Category
				if err := tx.Where("name ILIKE ?", name).First(&category).Error; err != nil {
					tx.Rollback()
					http.Error(w, "Category not found", http.StatusBadRequest)
					return
				}

				wc := models.WorkerCategory{WorkerID: userID, CategoryID: category.ID}
				if err := tx.Create(&wc).Error; err != nil {
					tx.Rollback()
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}
			}
		}

		// Если пришли данные для воркера или категорий, помечаем профиль как активный
		if len(workerUpdates) > 0 || len(input.CategoryNames) > 0 {
			workerUpdates["have_worker_profile"] = true
		}

		if len(userUpdates) == 0 && len(workerUpdates) == 0 && len(input.CategoryNames) == 0 {
			tx.Rollback()
			http.Error(w, "Database error, no fields", http.StatusInternalServerError)
			return
		}

		if len(workerUpdates) > 0 {
			result := tx.Model(&models.WorkerProfile{}).Where("user_id = ?", userID).Updates(workerUpdates)
			if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				tx.Rollback()
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		var updatedUser models.User
		if err := db.Preload("Role").Preload("WorkerProfile").First(&updatedUser, userID).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":             updatedUser.ID,
			"email":          updatedUser.Email,
			"name":           updatedUser.Name,
			"role":           updatedUser.Role.RoleName,
			"worker":         updatedUser.WorkerProfile != nil,
			"user_updates":   userUpdates,
			"worker_updates": workerUpdates,
		})
	}
}
