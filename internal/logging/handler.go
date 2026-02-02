package logging

import (
	"log/slog"
	"os"
)

func init() {
	var handler slog.Handler

	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	switch env {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, nil)
	case "development":
		handler = slog.NewTextHandler(os.Stdout, nil)
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
