package log

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithFields(t *testing.T) {
	ctx := WithFields(context.Background())
	flds := Fields(ctx)

	t.Run("empty log context", func(t *testing.T) {
		if len(flds) != 0 {
			t.Fatal("initial fields should be empty", flds)
		}
	})

	t.Run("invalid number of keyvals", func(t *testing.T) {
		ctx = WithFields(ctx, "key1", 1, "key2")
		flds = Fields(ctx)
		if len(flds) != 0 {
			t.Fatalf("expected 0 fields on error got %d", len(flds))
		}
	})

	t.Run("valid log context", func(t *testing.T) {
		ctx = WithFields(ctx, "key1", 1, "key2", "val2")
		flds = Fields(ctx)
		if len(flds) != 2 {
			t.Fatalf("expected 2 fields got %d", len(flds))
		}
		if flds[0].Key != "key1" || flds[0].Integer != 1 {
			t.Errorf(`expected "key1"=1, got %q=%+v`, flds[0].Key, flds[0].Integer)
		}
	})
}

func TestWithFieldsInfo(t *testing.T) {
	ctx := WithFields(context.Background(), "key1", 1)
	core, o := observer.New(zap.InfoLevel)
	l := NewLogger(zap.New(core))
	l.Info(ctx, "testing", "key2", "val2")

	flds := o.FilterMessage("testing").All()[0].Context

	if len(flds) != 2 {
		t.Fatalf("expected 2 fields got %d", len(flds))
	}
	if flds[1].Key != "key1" || flds[1].Integer != 1 {
		t.Errorf(`expected "key1"=1, got %q=%+v`, flds[0].Key, flds[0].Integer)
	}
}
