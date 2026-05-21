package logging

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"

	"yv35.com/dotfiles-cli/internal/cli"
)

func init() {
	var handler slog.Handler

	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	logPath := path.Join(cli.LOGS_PATH, "log.txt")

	// Ensure logs directory exists
	logDir := path.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Errorf("failed to create logs directory: %v", err))
	}

	fd, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("failed to setup logging: %v", err))
	}

	switch env {
	case "production":
		handler = slog.NewJSONHandler(fd, nil)
	case "development":

		handler = slog.NewTextHandler(fd, nil)
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Logging initialized", "env", env)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		slog.Info("Received interrupt signal, flushing logs and exiting...")

		if err := fd.Sync(); err != nil {
			slog.Error("Failed to sync log file", "error", err)
		}

		if err := fd.Close(); err != nil {
			slog.Error("Failed to close log file", "error", err)
		}
	}()
}
