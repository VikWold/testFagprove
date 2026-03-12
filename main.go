package main

import (
	"log/slog"
	"os"

	"github.com/google/uuid"
)

var (
	version = vcs.Version()
)

type application struct {
	logger *slog.Logger
}

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

	app := &application{
		logger: instanceLogger,
	}

	app.logger.Info("Starting API")
}
