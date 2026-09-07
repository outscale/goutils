/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package log

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/outscale/goutils/sdk/sanitize"
	"k8s.io/klog/v2"
)

const maxResponseLength = 500

func cleanBody(r io.ReadCloser) (string, error) {
	defer r.Close() //nolint:errcheck
	buf, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(buf), `"`, ``), nil
}

func truncatedBody(body string) string {
	str := []rune(body)
	if len(str) > maxResponseLength {
		return string(str[:maxResponseLength/2]) + " [truncated] " + string(str[len(str)-maxResponseLength/2:])
	}
	return string(str)
}

type OAPILogger struct{}

func callName(r *http.Request) string {
	return path.Base(r.URL.Path)
}

func (OAPILogger) Request(ctx context.Context, req any)   {}
func (OAPILogger) Response(ctx context.Context, resp any) {}

func (l OAPILogger) RequestHttp(ctx context.Context, req *http.Request) {
	logger := klog.FromContext(ctx).WithCallDepth(1)
	if !logger.V(5).Enabled() {
		return
	}
	req = sanitize.HTTPRequest(req)
	body, err := cleanBody(req.Body)
	if err != nil {
		l.Error(ctx, fmt.Errorf("log request: %w", err))
		return
	}
	logger.Info("OAPI request: "+body, "OAPI", callName(req))
}

func (l OAPILogger) ResponseHttp(ctx context.Context, resp *http.Response, d time.Duration) {
	logger := klog.FromContext(ctx).WithCallDepth(1)
	call := callName(resp.Request)
	if resp.StatusCode < 300 && !logger.V(5).Enabled() {
		return
	}
	resp = sanitize.HTTPResponse(resp) //nolint:bodyclose
	body, err := cleanBody(resp.Body)
	if err != nil {
		l.Error(ctx, fmt.Errorf("log response: %w", err))
		return
	}
	switch {
	case resp.StatusCode > 299:
		logger.V(3).Info("OAPI error response: "+body, "OAPI", call, "http_status", resp.Status, "duration", d)
	case logger.V(5).Enabled(): // no error
		// do not truncate response in level 6
		if !logger.V(6).Enabled() {
			body = truncatedBody(body)
		}
		logger.Info("OAPI response: "+body, "OAPI", call, "duration", d)
	}
}

func (OAPILogger) Error(ctx context.Context, err error) {
	logger := klog.FromContext(ctx).WithCallDepth(1)
	logger.V(3).Error(err, "OAPI error")
}
