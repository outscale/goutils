/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package sdk

import (
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/middleware"
	options "github.com/outscale/osc-sdk-go/v3/pkg/utils"
	"github.com/spf13/pflag"
)

const (
	// DefaultRateLimit defines the default rate limit.
	DefaultRateLimit = 5
	// DefaultRetryWaitMin defines the default wait between retries.
	DefaultRetryWaitMin = time.Second
	// DefaultRetryWaitMax defines the max wait between retries
	DefaultRetryWaitMax = 30 * time.Second
	// DefaultRetryCount defines the max number of retries.
	DefaultRetryCount = 5
)

// Options defines the SDK options (rate limit & backoff)
type Options struct {
	RateLimit                  int
	RetryWaitMin, RetryWaitMax time.Duration
	RetryCount                 int
}

// AddFlags adds flags for SDK options to a flag set.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.RateLimit, "oapi-rate-limit", DefaultRateLimit, "Maximum rate of Outscale API calls (per second)")
	fs.DurationVar(&o.RetryWaitMin, "oapi-retry-wait-min", DefaultRetryWaitMin, "Minimum wait between retries")
	fs.DurationVar(&o.RetryWaitMax, "oapi-retry-wait-max", DefaultRetryWaitMax, "Maximum wait between retries")
	fs.IntVar(&o.RetryCount, "oapi-retry-count", DefaultRetryCount, "Maximum number of retries")
}

func (o *Options) middleware() []middleware.MiddlewareChainOption {
	opts := []middleware.MiddlewareChainOption{options.WithRatelimit(o.RateLimit)}
	if o.RetryCount > 0 {
		opts = append(opts, options.WithRetry(&o.RetryWaitMin, &o.RetryWaitMax, &o.RetryCount))
	} else {
		opts = append(opts, options.WithoutRetry())
	}
	return opts
}
