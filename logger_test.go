// Copyright (C) 2017 ScyllaDB

package log

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	pkgErrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

var diffOpts = cmp.Options{
	cmpopts.IgnoreTypes(zapcore.EntryCaller{}),
	cmpopts.IgnoreFields(zapcore.Entry{}, "Stack"),
}

func TestLogger(t *testing.T) {
	table := []struct {
		msg       string
		expectMsg string
	}{
		{"foo", "foo"},
		{"", ""},
	}

	// Common to all test cases.
	ctx := WithTraceID(context.Background())
	base := []interface{}{"foo", "bar"}
	extra := []interface{}{"baz", false}
	expectedFields := []zapcore.Field{
		zap.String("foo", "bar"),
		zap.Bool("baz", false),
		zap.String("_trace_id", TraceID(ctx)),
	}

	for _, test := range table {
		withLogger(t, Config{Level: zap.NewAtomicLevelAt(zap.DebugLevel)}, func(logger Logger, logs *observer.ObservedLogs) {
			logger.With(base...).Debug(ctx, test.msg, extra...)
			logger.With(base...).Info(ctx, test.msg, extra...)
			logger.With(base...).Error(ctx, test.msg, extra...)

			expected := make([]observer.LoggedEntry, 3)
			for i, lvl := range []zapcore.Level{zap.DebugLevel, zap.InfoLevel, zap.ErrorLevel} {
				expected[i] = observer.LoggedEntry{
					Entry:   zapcore.Entry{Message: test.expectMsg, Level: lvl},
					Context: expectedFields,
				}
			}
			if diff := cmp.Diff(expected, logs.AllUntimed(), diffOpts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestStringifyErrors(t *testing.T) {
	ctx := context.Background()
	err := pkgErrors.Wrapf(pkgErrors.New("inner"), "outer")

	withLogger(t, Config{Level: zap.NewAtomicLevelAt(zap.DebugLevel)}, func(logger Logger, logs *observer.ObservedLogs) {
		logger.Debug(ctx, "msg", "error", err)
		logger.Info(ctx, "msg", "error", err)
		logger.Error(ctx, "msg", "error", err)

		expected := []observer.LoggedEntry{
			{
				Entry:   zapcore.Entry{Message: "msg", Level: zap.DebugLevel},
				Context: []zapcore.Field{zap.String("error", err.Error())},
			},
			{
				Entry:   zapcore.Entry{Message: "msg", Level: zap.InfoLevel},
				Context: []zapcore.Field{zap.String("error", err.Error())},
			},
			{
				Entry:   zapcore.Entry{Message: "msg", Level: zap.ErrorLevel},
				Context: []zapcore.Field{zap.String("error", err.Error()), zap.String("errorStack", "github.com/scylladb/go-log.TestStringifyErrors")},
			},
		}

		opt := cmp.Comparer(func(x, y string) bool {
			if len(x) > 20 {
				x = x[:20]
			}
			if len(y) > 20 {
				y = y[:20]
			}
			return x == y
		})

		if diff := cmp.Diff(expected, logs.AllUntimed(), append(diffOpts, opt)); diff != "" {
			t.Error(diff)
		}
	})
}

func withLogger(t *testing.T, c Config, f func(Logger, *observer.ObservedLogs)) {
	t.Helper()

	fac, logs := observer.New(c.Level)
	l, err := NewProduction(c, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, fac)
	}))
	if err != nil {
		t.Fatal("NewProduction() error", err)
	}
	f(l, logs)
}

func newZapLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	enc := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	return zap.New(zapcore.NewCore(
		enc,
		&zaptest.Discarder{},
		zapcore.DebugLevel,
	))
}

func BenchmarkZap(b *testing.B) {
	t := newTraceID()
	l := newZapLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("log message", zap.Int("key0", 0), zap.String("key1", "key1"), zap.String("key2", "key2"), zap.String("_trace_id", t))
	}
}

func BenchmarkZapSugared(b *testing.B) {
	t := newTraceID()
	l := newZapLogger().Sugar()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debugw("log message", "key0", 0, "key1", "key1", "key2", "key2", "_trace_id", t)
	}
}

func BenchmarkLogger(b *testing.B) {
	ctx := WithTraceID(context.Background())
	l := Logger{base: newZapLogger()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug(ctx, "log message", "key0", 0, "key1", "key1", "key2", "key2")
	}
}
