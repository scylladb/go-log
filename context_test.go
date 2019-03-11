package log

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithFields(t *testing.T) {
	table := []struct {
		name    string
		keyvals []interface{}
		flds    []zapcore.Field
	}{
		{
			"empty keyvals",
			nil,
			nil,
		},
		{
			"invalid number of keyvals",
			[]interface{}{"key1", 1, "key2"},
			nil,
		},
		{
			"valid number of keyvals",
			[]interface{}{"key1", 1, "key2", "val2"},
			[]zapcore.Field{
				{Key: "key1", Integer: 1},
				{Key: "key2", String: "val2"},
			},
		},
		{
			"duplicate keys",
			[]interface{}{"key1", 1, "key1", "val1"},
			[]zapcore.Field{
				{Key: "key1", String: "val1"},
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			ctx := WithFields(context.Background(), test.keyvals...)
			flds := Fields(ctx)
			opts := cmpopts.IgnoreFields(zapcore.Field{}, "Type")
			if diff := cmp.Diff(test.flds, flds, opts); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestWithFieldsLogger(t *testing.T) {
	ctx := WithFields(context.Background(), "key1", 1)
	core, o := observer.New(zap.InfoLevel)
	l := NewLogger(zap.New(core))

	t.Run("info", func(t *testing.T) {
		l.Info(ctx, "testing", "key2", "val2")
		flds := o.TakeAll()[0].Context

		if len(flds) != 2 {
			t.Fatalf("expected 2 fields got %d", len(flds))
		}
		if flds[0].Key != "key1" || flds[0].Integer != 1 {
			t.Errorf(`expected "key1"=1, got %q=%+v`, flds[0].Key, flds[0].Integer)
		}
	})

	t.Run("duplicate key collision with log call", func(t *testing.T) {
		l.Info(ctx, "testing", "key1", "val2")
		flds := o.TakeAll()[0].Context

		if len(flds) != 1 {
			t.Log(flds)
			t.Fatalf("expected one field got %d", len(flds))
		}
		if flds[0].Key != "key1" || flds[0].String != "val2" {
			t.Errorf(`expected "key1"=1, got %q=%+v`, flds[0].Key, flds[0].String)
		}
	})

	t.Run("duplicate keys collision with With", func(t *testing.T) {
		withLogger := l.With("key1", "val2", "key1", "val3")
		withLogger.Info(ctx, "testing")
		flds := o.TakeAll()[0].Context

		if len(flds) != 1 {
			t.Log(flds)
			t.Fatalf("expected one field got %d", len(flds))
		}
		if flds[0].Key != "key1" || flds[0].String != "val3" {
			t.Errorf(`expected "key1"=1, got %q=%+v`, flds[0].Key, flds[0].String)
		}
	})
}
