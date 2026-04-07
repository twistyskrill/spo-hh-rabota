package worker

import (
	"encoding/json"
	"errors"
	"go-api/internal/models"
	"go-api/internal/storage"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type availabilitySlotResponse struct {
	ID       uint      `json:"id"`
	StartAt  time.Time `json:"start_at"`
	IsBooked bool      `json:"is_booked"`
}

func parseRange(r *http.Request) (time.Time, time.Time, error) {
	now := time.Now()
	from := now
	to := now.AddDate(0, 0, 60)

	if fromStr := strings.TrimSpace(r.URL.Query().Get("from")); fromStr != "" {
		parsed, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		from = parsed
	}
	if toStr := strings.TrimSpace(r.URL.Query().Get("to")); toStr != "" {
		parsed, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		to = parsed
	}
	if !to.After(from) {
		return time.Time{}, time.Time{}, errors.New("invalid range")
	}
	return from, to, nil
}

func normalizeSlots(raw []string) ([]time.Time, error) {
	unique := make(map[time.Time]struct{}, len(raw))
	now := time.Now().Add(-1 * time.Minute)

	for _, s := range raw {
		ts := strings.TrimSpace(s)
		if ts == "" {
			continue
		}
		parsed, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return nil, err
		}
		if parsed.Before(now) {
			continue
		}
		// Normalize to minute precision to avoid duplicate seconds/nanos.
		normalized := parsed.UTC().Truncate(time.Minute)
		unique[normalized] = struct{}{}
	}

	out := make([]time.Time, 0, len(unique))
	for t := range unique {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Before(out[j]) })
	return out, nil
}

func GetWorkerAvailabilityHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		workerID := uint(id)

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

		from, to, err := parseRange(r)
		if err != nil {
			http.Error(w, "Invalid date range", http.StatusBadRequest)
			return
		}

		var slots []models.WorkerAvailabilitySlot
		err = db.
			Where("worker_id = ? AND start_at >= ? AND start_at <= ? AND is_booked = ?", workerID, from, to, false).
			Order("start_at ASC").
			Find(&slots).Error
		if err != nil {
			logger.Error("failed to list worker availability", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		resp := make([]availabilitySlotResponse, 0, len(slots))
		for _, s := range slots {
			resp = append(resp, availabilitySlotResponse{ID: s.ID, StartAt: s.StartAt, IsBooked: s.IsBooked})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"slots": resp})
	}
}

func MyAvailabilityHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var profile models.WorkerProfile
		if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
			http.Error(w, "Worker profile not found", http.StatusNotFound)
			return
		}
		if !profile.HaveWorkerProfile {
			http.Error(w, "Worker profile not active", http.StatusForbidden)
			return
		}

		from, to, err := parseRange(r)
		if err != nil {
			http.Error(w, "Invalid date range", http.StatusBadRequest)
			return
		}

		var slots []models.WorkerAvailabilitySlot
		err = db.
			Where("worker_id = ? AND start_at >= ? AND start_at <= ?", userID, from, to).
			Order("start_at ASC").
			Find(&slots).Error
		if err != nil {
			logger.Error("failed to list my availability", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		resp := make([]availabilitySlotResponse, 0, len(slots))
		for _, s := range slots {
			resp = append(resp, availabilitySlotResponse{ID: s.ID, StartAt: s.StartAt, IsBooked: s.IsBooked})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"slots": resp})
	}
}

func UpdateMyAvailabilityHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var profile models.WorkerProfile
		if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
			http.Error(w, "Worker profile not found", http.StatusNotFound)
			return
		}
		if !profile.HaveWorkerProfile {
			http.Error(w, "Worker profile not active", http.StatusForbidden)
			return
		}

		var in struct {
			Slots []string `json:"slots"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		normalized, err := normalizeSlots(in.Slots)
		if err != nil {
			http.Error(w, "Invalid slot format", http.StatusBadRequest)
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		var existing []models.WorkerAvailabilitySlot
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("worker_id = ? AND start_at >= ?", userID, time.Now().Add(-1*time.Minute)).
			Find(&existing).Error; err != nil {
			tx.Rollback()
			logger.Error("failed to load existing slots", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		targetSet := make(map[time.Time]struct{}, len(normalized))
		for _, t := range normalized {
			targetSet[t] = struct{}{}
		}

		existingFree := make(map[time.Time]models.WorkerAvailabilitySlot)
		for _, s := range existing {
			if !s.IsBooked {
				existingFree[s.StartAt.UTC().Truncate(time.Minute)] = s
			}
		}

		for _, t := range normalized {
			if _, ok := existingFree[t]; ok {
				continue
			}
			newSlot := models.WorkerAvailabilitySlot{WorkerID: userID, StartAt: t, IsBooked: false}
			if err := tx.Create(&newSlot).Error; err != nil {
				tx.Rollback()
				logger.Error("failed to create availability slot", "error", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

		for t, s := range existingFree {
			if _, ok := targetSet[t]; ok {
				continue
			}
			if err := tx.Delete(&models.WorkerAvailabilitySlot{}, s.ID).Error; err != nil {
				tx.Rollback()
				logger.Error("failed to delete availability slot", "error", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "availability updated", "slots_count": len(normalized)})
	}
}

func BookWorkerSlotHandler(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		workerID := uint(id)
		if workerID == userID {
			http.Error(w, "Cannot book yourself", http.StatusBadRequest)
			return
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

		var in struct {
			StartAt string `json:"start_at"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(in.StartAt))
		if err != nil {
			http.Error(w, "Invalid start_at", http.StatusBadRequest)
			return
		}
		startAt = startAt.UTC().Truncate(time.Minute)
		if startAt.Before(time.Now().Add(-1 * time.Minute)) {
			http.Error(w, "Slot is in the past", http.StatusBadRequest)
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		var slot models.WorkerAvailabilitySlot
		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("worker_id = ? AND start_at = ?", workerID, startAt).
			First(&slot).Error
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Slot not found", http.StatusNotFound)
				return
			}
			logger.Error("failed to load slot", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if slot.IsBooked {
			tx.Rollback()
			http.Error(w, "Slot already booked", http.StatusConflict)
			return
		}

		now := time.Now().UTC()
		slot.IsBooked = true
		slot.BookedByUserID = &userID
		slot.BookedAt = &now
		if err := tx.Save(&slot).Error; err != nil {
			tx.Rollback()
			logger.Error("failed to book slot", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"slot_id":   slot.ID,
			"worker_id": workerID,
			"start_at":  slot.StartAt,
			"booked":    true,
		})
	}
}
