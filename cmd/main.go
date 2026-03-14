package main

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	vcs "testFagprove/internal"
	"testFagprove/internal/config"
	"testFagprove/internal/data"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	version = vcs.Version()
)

type application struct {
	logger *slog.Logger
	models data.Models
	config *config.Config
}

// @version 1.0
func main() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	instanceLogger := logger.With(
		slog.Group(
			"application_instance",
			slog.String("version", version),
			slog.String("instance_id", uuid.New().String()),
		),
	)

	slog.Info("loadig configuration")
	cfg, err := config.New()
	if err != nil {
		slog.Error("unable to load configuration", "error", err)
		os.Exit(1)
	}
	slog.Info("configuration loaded", "configuration", cfg)

	if cfg.DB == nil {
		logger.Error("unable to load configuration")
		os.Exit(1)
	}

	logger.Info("opening database connection pool...")
	logger.Info("DSN", "dsn", cfg.DB.DSN)
	db, err := openDB(cfg.DB.DSN)
	if err != nil {
		logger.Error("unable to open database connection pool", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	queryTimeout := time.Duration(cfg.DB.Timeout) * time.Second
	logger.Info("database connection pool established")

	app := &application{
		logger: instanceLogger,
		models: data.NewModels(db, &queryTimeout),
	}

	app.logger.Info("Starting API")

	err = app.serve(":8080")
	if err != nil {
		logger.Error("unable to start server", "error", err)
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
