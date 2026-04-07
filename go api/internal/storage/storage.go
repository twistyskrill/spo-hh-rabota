package storage

import (
	"context"
	"fmt"
	"go-api/internal/config"
	"go-api/internal/models"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	db *gorm.DB
}

func NewDB(cfg config.DB, basicLogger *slog.Logger) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Europe/Moscow client_encoding=UTF8",
		cfg.Host, cfg.User, cfg.Password, cfg.DBname, cfg.Port, cfg.SSLmode)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // стандартный log.Logger
		logger.Config{
			SlowThreshold:             time.Second, // Медленные запросы
			LogLevel:                  logger.Warn, // Только ошибки + медленные
			IgnoreRecordNotFoundError: true,        // Игнор "record not found"
			Colorful:                  true,        // Цветной вывод
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{ //открытие бд
		Logger: gormLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("postgres connect failed: %w", err)
	}

	sqlDB, err := db.DB() //получение пула соединений

	if err != nil {
		return nil, fmt.Errorf("get sql.DB failed: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)       //Устанавливает максимальное количество простаивающих соединений в пуле
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)       //Максимальное количество открытых соединений одновременно
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) //Максимальное время жизни соединения

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	basicLogger.Info("PostgreSQL подключен",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("db", cfg.DBname))

	// TODO: миграции
	db.AutoMigrate(
		&models.Role{},
		&models.Category{},
		&models.PriceUnit{},

		&models.User{},
		&models.WorkerProfile{},

		&models.Ad{},
		&models.Review{},
		&models.Response{},
		&models.WorkerAvailabilitySlot{},

		&models.WorkerCategory{},
		&models.BlackList{},
	)

	if err := migrateWorkerProfileDescription(db); err != nil {
		return nil, fmt.Errorf("migrate worker_profiles.description failed: %w", err)
	}

	return &Postgres{db: db}, nil
}

type columnMeta struct {
	DataType               string `gorm:"column:data_type"`
	CharacterMaximumLength *int   `gorm:"column:character_maximum_length"`
}

func migrateWorkerProfileDescription(db *gorm.DB) error {
	// Ensure we can store longer descriptions than 255 chars.
	// AutoMigrate does not always widen existing VARCHAR columns reliably.
	var meta columnMeta
	err := db.Raw(`
		SELECT data_type, character_maximum_length
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'worker_profiles'
		  AND column_name = 'description'
		LIMIT 1
	`).Scan(&meta).Error
	if err != nil {
		return err
	}

	// If column doesn't exist yet, AutoMigrate will create it as TEXT from the model tag.
	if meta.DataType == "" {
		return nil
	}

	if meta.DataType == "character varying" && meta.CharacterMaximumLength != nil && *meta.CharacterMaximumLength > 0 {
		return db.Exec(`ALTER TABLE worker_profiles ALTER COLUMN description TYPE text`).Error
	}

	return nil
}

func (p *Postgres) DB() *gorm.DB {
	return p.db
}

func (p *Postgres) Close() error {
	sqlDB, err := p.db.DB()

	if err != nil {
		return fmt.Errorf("postgres close failed: %w", err)
	}

	return sqlDB.Close()
}
