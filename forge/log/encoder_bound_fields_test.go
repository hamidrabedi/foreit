package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Finding 4: One-line mode must include fields bound via logger.With in order, before per-call fields.
func TestConsoleEncoder_OneLine_IncludesBoundFields(t *testing.T) {
	enc := NewConsoleEncoder(DevelopmentConfig{
		Colored: false,
		OneLine: true,
	})

	buf := &bytes.Buffer{}
	core := zapcore.NewCore(enc, zapcore.AddSync(buf), zapcore.DebugLevel)
	logger := zap.New(core)

	logger.With(zap.String("request_id", "abc")).Info("msg")

	output := buf.String()
	assert.Contains(t, output, "request_id=abc")
}

func TestConsoleEncoder_OneLine_FieldOrder(t *testing.T) {
	enc := NewConsoleEncoder(DevelopmentConfig{
		Colored: false,
		OneLine: true,
	})

	buf := &bytes.Buffer{}
	core := zapcore.NewCore(enc, zapcore.AddSync(buf), zapcore.DebugLevel)
	logger := zap.New(core)

	// Bound fields must appear in order added, before per-call fields
	logger.With(zap.String("first", "1"), zap.String("second", "2")).Info("msg", zap.String("call", "3"))

	output := buf.String()
	assert.Contains(t, output, "first=1 second=2 call=3")
}
