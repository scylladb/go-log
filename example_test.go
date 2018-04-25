package log_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/scylladb/go-log"
	"go.uber.org/zap/zapcore"
)

func TestExample(t *testing.T) {
	ctx := log.WithTraceID(context.Background())

	logger, err := log.NewProduction(log.Config{
		Mode:  log.SyslogMode,
		Level: zapcore.InfoLevel,
	})
	if err != nil {
		t.Fatal(err)
	}
	logger.Info(ctx, "Could not connect to database",
		"sleep", 5*time.Second,
		"error", errors.New("I/O error"),
	)

	logger.Named("sub").Error(ctx, "Unexpected error", "error", errors.New("unexpected"))
}
