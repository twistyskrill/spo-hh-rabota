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

		&models.WorkerCategory{},
		&models.BlackList{},
	)

	return &Postgres{db: db}, nil
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
