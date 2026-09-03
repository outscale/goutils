package log

import (
	"context"
	"errors"

	"github.com/outscale/goutils/sdk/sanitize"
	"github.com/samber/lo"
)

type sanitizingLogger struct {
	l Logger
}

// Sanitize will sanitize all log lines: messages, key/values, errors
func Sanitize(l Logger) Logger {
	return &sanitizingLogger{
		l: l,
	}
}

func (l *sanitizingLogger) Info(ctx context.Context, msg string, kv ...any) {
	l.l.Info(
		ctx,
		sanitize.Sanitize(msg),
		lo.Map(kv, func(v any, _ int) any {
			return sanitize.Sanitize(v)
		}),
	)
}

func (l *sanitizingLogger) Error(ctx context.Context, err error, msg string, kv ...any) {
	l.l.Error(
		ctx,
		// Wrapped errors cannot be sanitized directly, as the underlying error is within a private field.
		errors.New(sanitize.Sanitize(err.Error())),
		sanitize.Sanitize(msg),
		lo.Map(kv, func(v any, _ int) any {
			return sanitize.Sanitize(v)
		}),
	)
}

var _ Logger = (*sanitizingLogger)(nil)
