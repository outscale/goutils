package log_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/outscale/goutils/sdk/log"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	logged string
}

func (l *testLogger) Info(ctx context.Context, msg string, kv ...any) {
	l.logged = msg + " " + fmt.Sprint(kv...)
}

func (l *testLogger) Error(ctx context.Context, err error, msg string, kv ...any) {
	l.logged = err.Error() + " " + msg + " " + fmt.Sprint(kv...)
}

func TestSanitizeLogs(t *testing.T) {
	type tc struct {
		err      error
		msg      string
		kv       []any
		sanitize bool
	}
	sk := "0A1B2C3D4D5E6F7G8H9I0A1B2C3D4D5E6F7G8H9I"
	tcs := []tc{
		{msg: sk, sanitize: true},
		{kv: []any{"foo", sk}, sanitize: true},
		{kv: []any{"foo", osc.AccessKeySecretKey{SecretKey: &sk}}, sanitize: true},
		{err: fmt.Errorf("sk error %q", sk), msg: "foo", sanitize: true},
		{err: errors.New("foo"), msg: sk, sanitize: true},
		{err: errors.New("foo"), msg: sk, sanitize: true},
		{err: errors.New("foo"), kv: []any{"foo", sk}, sanitize: true},
		{err: errors.New(sk), sanitize: true},
		{err: fmt.Errorf("%w", errors.New(sk)), sanitize: true},
		{err: errors.Join(errors.New("foo"), errors.New(sk)), sanitize: true},
	}
	for _, tc := range tcs {
		l := &testLogger{}
		sl := log.Sanitize(l)
		if tc.err != nil {
			sl.Error(t.Context(), tc.err, tc.msg, tc.kv...)
		} else {
			sl.Info(t.Context(), tc.msg, tc.kv...)
		}
		if tc.sanitize {
			assert.NotContainsf(t, l.logged, sk, "sanitized string %q should not contain sk %q", l.logged, sk)
		}
	}
}
