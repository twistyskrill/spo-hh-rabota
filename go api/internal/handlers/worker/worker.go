package worker

import (
	"encoding/json"
	"errors"
	"go-api/internal/storage"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func AllWorkersHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 10
		offset := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			limit, _ = strconv.Atoi(l)
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			offset, _ = strconv.Atoi(o)
		}

		workers, total, err := storage.ListApprovedWorkers(db, limit, offset)
		if err != nil {
			logger.Error("Failed to get workers", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"workers": workers,
			"pagination": map[string]interface{}{
				"total":  total,
				"limit":  limit,
				"offset": offset,
				"page":   offset/limit + 1,
				"pages":  (total + int64(limit) - 1) / int64(limit),
			},
		})
	}
}

func WorkerHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
		}

		worker, err := storage.WorkerByID(db, uint(id))
		if err != nil || worker == nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Worker not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(worker)
	}
}
