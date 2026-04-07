package main

import (
	"go-api/internal/auth"
	"go-api/internal/config"
	handlerAdmin "go-api/internal/handlers/admin"
	handlerAds "go-api/internal/handlers/ads"
	handlerAuth "go-api/internal/handlers/auth"
	handlerInfo "go-api/internal/handlers/info"
	handlerReviews "go-api/internal/handlers/reviews"
	handlerSys "go-api/internal/handlers/sys"
	handlerWork "go-api/internal/handlers/worker"
	appmw "go-api/internal/middleware"
	"go-api/internal/storage"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad() //init config cleanenv

	logger := setupLogger(cfg.Env) //init logger slog
	logger.Info("starting", slog.String("env", cfg.Env))

	store, err := storage.NewDB(cfg.DB, logger) // init storage postgresql
	auth.Init(cfg.JWT.SecretKey)                //init secret key

	if err != nil {
		logger.Error("DB init failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer store.Close()

	r := chi.NewRouter() // init router chi

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(appmw.CORSMiddleware)

	handlerAuth.SetupRoutes(store.DB(), logger, r)
	handlerSys.SetupRoutes(store.DB(), logger, r)
	handlerWork.SetupRoutes(store.DB(), logger, r)
	handlerInfo.SetupRoutes(store.DB(), logger, r)
	handlerAds.SetupRoutes(store.DB(), logger, r)
	handlerReviews.SetupRoutes(store.DB(), logger, r)
	handlerAdmin.SetupRoutes(store.DB(), logger, r) // Админ-панель

	logger.Info("server started", slog.String("port", ":8080"))
	http.ListenAndServe(":8080", r)
}

// константы логгера
const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// создание логгера
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
