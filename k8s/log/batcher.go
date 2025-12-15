/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package log

import (
	"context"

	"k8s.io/klog/v2"
)

type Logger struct{}

func (Logger) Info(ctx context.Context, msg string, kv ...any) {
	klog.FromContext(ctx).V(5).Info(msg, kv...)
}

func (Logger) Error(ctx context.Context, err error, msg string, kv ...any) {
	klog.FromContext(ctx).V(3).Error(err, msg, kv...)
}
