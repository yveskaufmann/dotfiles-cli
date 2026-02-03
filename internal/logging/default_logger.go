package logging

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"
)

func init() {
	var handler slog.Handler

	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to setup logging: %v", err))
	}

	fd, err := os.OpenFile(path.Join(cwd, "logs/log.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
