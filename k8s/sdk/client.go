/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package sdk

import (
	"context"
	"errors"
	"fmt"

	"dario.cat/mergo"
	"github.com/outscale/goutils/k8s/log"
	"github.com/outscale/goutils/sdk/metadata"
	"github.com/outscale/osc-sdk-go/v3/pkg/middleware"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	options "github.com/outscale/osc-sdk-go/v3/pkg/utils"
)

// NewSDKClient builds a new SDK client, querying the metadata server is OSC_REGION is not set,
// configuring backoff & ratelimiter based on opts, and checking credentials.
func NewSDKClient(ctx context.Context, ua string, opts ...Options) (*profile.Profile, osc.ClientInterface, error) {
	prof, err := profile.NewProfileFromStandardConfiguration("", "")
	if err != nil {
		return nil, nil, fmt.Errorf("load config: %w", err)
	}
	if prof.AccessKey == "" || prof.SecretKey == "" {
		return nil, nil, errors.New("OSC_ACCESS_KEY/OSC_SECRET_KEY are required")
	}
	if prof.Region == "" {
		prof.Region, err = metadata.GetRegion(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to fetch metadata: %w", err)
		}
	}
	lg := log.OAPILogger{}
	copts := []middleware.MiddlewareChainOption{options.WithUseragent(ua), options.WithLogging(lg)}
	if len(opts) > 0 {
		opt := opts[0]
		// no default is set on RetryCount, it might be valid to run without backoff.
		// ratelimiter is always configured. 0 values will be replaced by defaults.
		err = mergo.Merge(&opt, Options{
			RateLimit:    DefaultRateLimit,
			RetryWaitMin: DefaultRetryWaitMin,
			RetryWaitMax: DefaultRetryWaitMax,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("unable to set OAPI SDK options: %w", err)
		}
		copts = append(copts, opt.middleware()...)
	}
	client, err := osc.NewClient(prof, copts...)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize OAPI client: %w", err)
	}
	if err := CheckCredentials(ctx, client); err != nil {
		return nil, nil, err
	}
	return prof, client, nil
}
