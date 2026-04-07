package storage

import (
	"go-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReviewResponse struct {
	ID     uint      `json:"id"`
	Author string    `json:"author"`
	Rating int       `json:"rating"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}

func IsApprovedWorker(db *gorm.DB, workerID uint) (bool, error) {
	var count int64
	err := db.Model(&models.WorkerProfile{}).
		Where("user_id = ? AND have_worker_profile = ? AND status = ?", workerID, true, "approved").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ListWorkerReviews(db *gorm.DB, workerID uint, limit, offset int) ([]ReviewResponse, int64, error) {
	var total int64
	if err := db.Model(&models.Review{}).
		Where("worker_id = ?", workerID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reviews []ReviewResponse
	err := db.Table("reviews r").
		Select("r.id, u.name as author, r.rating, r.text, r.date").
		Joins("JOIN users u ON u.id = r.user_id").
		Where("r.worker_id = ? AND r.deleted_at IS NULL", workerID).
		Order("r.date DESC").
		Limit(limit).
		Offset(offset).
		Scan(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func CreateReview(db *gorm.DB, userID uint, workerID uint, rating int, text string) (*ReviewResponse, error) {
	review := models.Review{
		Rating:   rating,
		Text:     text,
		Date:     time.Now(),
		UserID:   userID,
		WorkerID: workerID,
	}
	if err := db.Create(&review).Error; err != nil {
		return nil, err
	}

	var authorName string
	_ = db.Model(&models.User{}).Select("name").Where("id = ?", userID).Scan(&authorName).Error

	return &ReviewResponse{
		ID:     review.ID,
		Author: authorName,
		Rating: review.Rating,
		Text:   review.Text,
		Date:   review.Date,
	}, nil
}
