package models

import (
	"time"

	"gorm.io/gorm"
)

// ======================================================================
// СПРАВОЧНИКИ
// ======================================================================

type Role struct {
	gorm.Model
	RoleName string `gorm:"size:255;not null;uniqueIndex" json:"role_name"`

	Users []User `gorm:"foreignKey:RoleID" json:"-"`
}

type Category struct {
	gorm.Model
	Name string `gorm:"size:255;not null;uniqueIndex" json:"name"`

	Ads              []Ad             `gorm:"foreignKey:CategoryID" json:"ads,omitempty"`
	WorkerCategories []WorkerCategory `gorm:"foreignKey:CategoryID" json:"worker_categories,omitempty"`
}

type PriceUnit struct {
	gorm.Model
	Name string `gorm:"size:255;not null;uniqueIndex" json:"name"`

	Ads []Ad `gorm:"foreignKey:PriceUnitID" json:"ads,omitempty"`
}

// ======================================================================
// ОСНОВНЫЕ СУЩНОСТИ
// ======================================================================

type User struct {
	gorm.Model
	Name         string    `gorm:"size:255;not null" json:"name"`
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	RoleID       uint      `gorm:"not null;index;default:1" json:"role_id"`
	CreatedAt    time.Time `gorm:"not null;index" json:"created_at"`
	Phone        string    `gorm:"size:255" json:"phone"`

	// Связи
	Role          Role           `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Ads           []Ad           `gorm:"foreignKey:UserID" json:"ads,omitempty"`
	ReviewsGiven  []Review       `gorm:"foreignKey:UserID" json:"reviews_given,omitempty"`
	WorkerProfile *WorkerProfile `gorm:"foreignKey:UserID" json:"worker_profile,omitempty"` // 1:1
}

type WorkerProfile struct {
	gorm.Model
	UserID      uint    `gorm:"not null;uniqueIndex" json:"-"`
	ExpYears    *int    `json:"exp_years"`   // NULLABLE
	HourlyRate  *int    `json:"hourly_rate"` // NULLABLE
	Description *string `gorm:"type:text" json:"description"`
	IsBusy      bool    `gorm:"default:false;not null" json:"is_busy"`
	Location    string  `gorm:"size:255" json:"location"` // место жительства / работы
	Schedule    string  `gorm:"size:255" json:"schedule"` // расписание (часы / дни)
	// Пометка, что профиль реально заполнен и пользователь считается "рабочим"
	HaveWorkerProfile bool   `gorm:"default:false;not null" json:"have_worker_profile"`
	Status            string `gorm:"size:20;not null;default:'pending';index" json:"status"` // pending, approved, rejected

	// Связи (UserID вместо ID!)
	User            User             `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	ReviewsReceived []Review         `gorm:"foreignKey:WorkerID;references:UserID" json:"reviews_received,omitempty"`
	Categories      []WorkerCategory `gorm:"foreignKey:WorkerID;references:UserID" json:"categories,omitempty"`
}

type Ad struct {
	gorm.Model
	Title       string    `gorm:"size:255;not null" json:"title"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	PriceUnitID uint      `gorm:"not null;index" json:"price_unit_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Location    string    `gorm:"size:255" json:"location"` // локация объявления
	Schedule    string    `gorm:"size:255" json:"schedule"` // когда актуально объявление
	CreatedAt   time.Time `gorm:"not null;index" json:"created_at"`
	Status      string    `gorm:"size:20;not null;default:'pending';index" json:"status"` // pending, approved, rejected, in_progress, completed
	ExecutorID  *uint     `gorm:"index" json:"executor_id,omitempty"`                     // ID мастера, который взял в работу

	// Связи
	Category  Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	PriceUnit PriceUnit `gorm:"foreignKey:PriceUnitID" json:"price_unit,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Executor  *User     `gorm:"foreignKey:ExecutorID" json:"executor,omitempty"`
}

// ======================================================================
// M2M и вспомогательные таблицы
// ======================================================================

type WorkerCategory struct {
	WorkerID   uint `gorm:"primaryKey" json:"worker_id"`
	CategoryID uint `gorm:"primaryKey" json:"category_id"`

	Worker   WorkerProfile `gorm:"foreignKey:WorkerID;references:UserID" json:"-"`
	Category Category      `gorm:"foreignKey:CategoryID" json:"-"`
}

type Review struct {
	gorm.Model
	Rating   int       `gorm:"check:rating >= 1 AND rating <= 5;not null" json:"rating"`
	Text     string    `gorm:"size:255;not null" json:"text"`
	Date     time.Time `gorm:"not null;index" json:"date"`
	UserID   uint      `gorm:"not null;index" json:"user_id"`
	WorkerID uint      `gorm:"not null;index" json:"worker_id"`

	// Связи
	User   User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Worker WorkerProfile `gorm:"foreignKey:WorkerID;references:UserID" json:"worker,omitempty"`
}

type Response struct {
	gorm.Model
	AdID          uint      `gorm:"not null;index" json:"ad_id"`
	WorkerID      uint      `gorm:"not null;index" json:"worker_id"`                        // UserID мастера
	Message       string    `gorm:"size:500" json:"message"`                                // сообщение от мастера
	ProposedPrice *float64  `gorm:"type:decimal(10,2)" json:"proposed_price"`               // предлагаемая цена (опционально)
	Status        string    `gorm:"size:50;not null;default:'pending';index" json:"status"` // pending, accepted, rejected, cancelled
	CreatedAt     time.Time `gorm:"not null;index" json:"created_at"`

	// Связи
	Ad     Ad            `gorm:"foreignKey:AdID" json:"ad,omitempty"`
	Worker WorkerProfile `gorm:"foreignKey:WorkerID;references:UserID" json:"worker,omitempty"`
}

type BlackList struct {
	Email string `gorm:"primaryKey;size:255;not null" json:"email"`
}
