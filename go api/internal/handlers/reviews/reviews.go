package reviews

import (
	"encoding/json"
	"go-api/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func GetWorkerReviewsHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		workerID := uint(id)

		limit := 10
		offset := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil && v >= 0 {
				offset = v
			}
		}

		approved, err := storage.IsApprovedWorker(db, workerID)
		if err != nil {
			logger.Error("failed to check worker status", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if !approved {
			http.Error(w, "Worker not found", http.StatusNotFound)
			return
		}

		reviews, total, err := storage.ListWorkerReviews(db, workerID, limit, offset)
		if err != nil {
			logger.Error("failed to list reviews", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"reviews": reviews,
			"pagination": map[string]interface{}{
				"total":  total,
				"limit":  limit,
				"offset": offset,
			},
		})
	}
}

func CreateReviewHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		type input struct {
			WorkerID uint   `json:"worker_id"`
			Rating   int    `json:"rating"`
			Text     string `json:"text"`
		}

		var in input
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		in.Text = strings.TrimSpace(in.Text)
		if in.WorkerID == 0 {
			http.Error(w, "worker_id required", http.StatusBadRequest)
			return
		}
		if in.WorkerID == userID {
			http.Error(w, "Cannot review yourself", http.StatusBadRequest)
			return
		}
		if in.Rating < 1 || in.Rating > 5 {
			http.Error(w, "rating must be 1..5", http.StatusBadRequest)
			return
		}
		if in.Text == "" {
			http.Error(w, "text required", http.StatusBadRequest)
			return
		}

		approved, err := storage.IsApprovedWorker(db, in.WorkerID)
		if err != nil {
			logger.Error("failed to check worker status", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if !approved {
			http.Error(w, "Worker not found", http.StatusNotFound)
			return
		}

		created, err := storage.CreateReview(db, userID, in.WorkerID, in.Rating, in.Text)
		if err != nil {
			logger.Error("failed to create review", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	}
}
