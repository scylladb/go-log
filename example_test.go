package log_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/scylladb/go-log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestExample(t *testing.T) {
	ctx := log.WithTraceID(context.Background())

	atom := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, err := log.NewProduction(log.Config{
		Mode:     log.StderrMode,
		Level:    atom,
		Encoding: log.JSONEncoding,
	})
	if err != nil {
		t.Fatal(err)
	}
	logger.Info(ctx, "Could not connect to database",
		"sleep", 5*time.Second,
		"error", errors.New("I/O error"),
	)

	logger.Named("sub").Error(ctx, "Unexpected error", "error", errors.New("unexpected"))

	logger.Debug(ctx, "Logging on debug level is not printing anything now",
		"sleep", 5*time.Second,
		"error", errors.New("I/O error"),
	)

	atom.SetLevel(zapcore.DebugLevel)

	logger.Debug(ctx, "Logging on debug level requires the logger to be at that level",
		"sleep", 5*time.Second,
		"error", errors.New("I/O error"),
	)
}
