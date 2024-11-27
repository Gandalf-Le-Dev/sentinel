package main

import (
	"context"
	"errors"

	"github.com/Gandalf-Le-Dev/sentinel"
)

func main() {
	cfg := sentinel.DefaultConfig()
	cfg.ServiceName = "my-service"
	cfg.Level = sentinel.LevelDebug

	log, err := sentinel.New(cfg)
	if err != nil {
		panic(err)
	}

	// Basic logging
	log.Info("Server starting", "port", 8080)

	// Logging with additional fields
	log = log.With(
		"environment", "production",
		"version", "1.0.0",
	)

	// Logging with error
	err = errors.New("database connection failed")
	log.WithError(err).Error("Failed to connect to database")

	// Logging with context
	ctx := context.WithValue(context.Background(), "trace_id", "abc-123")
	log.WithContext(ctx).Info("Processing request")
}
