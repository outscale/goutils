/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package log

import "context"

var Default Logger = noLogger{}

type Logger interface {
	Info(ctx context.Context, msg string, kv ...any)
	Error(ctx context.Context, err error, msg string, kv ...any)
}

type noLogger struct{}

func (noLogger) Info(ctx context.Context, msg string, kv ...any)             {}
func (noLogger) Error(ctx context.Context, err error, msg string, kv ...any) {}
