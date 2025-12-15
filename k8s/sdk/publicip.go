/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package sdk

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/outscale/goutils/k8s/tags"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"k8s.io/klog/v2"
)

var (
	// ErrEmptyPool is returned when one requests an API from an empty pool.
	ErrEmptyPool = errors.New("no available IP in pool")
)

func AllocateIPFromPool(ctx context.Context, pool string, c osc.ClientInterface) (*osc.PublicIp, error) {
	log := klog.FromContext(ctx)
	log.V(4).Info("Fetching publicIps from pool", "pool", pool)
	req := osc.ReadPublicIpsRequest{
		Filters: &osc.FiltersPublicIp{
			TagKeys:   &[]string{tags.PublicIPPool},
			TagValues: &[]string{pool},
		},
	}
	resp, err := c.ReadPublicIps(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("allocate from pool: %w", err)
	}
	if len(*resp.PublicIps) == 0 {
		return nil, ErrEmptyPool
	}
	pips := *resp.PublicIps
	// randomly fetch from the list, to limit the chance of allocating
	// the same IP to two concurrent requests
	off := rand.IntN(len(pips)) //nolint:gosec
	for i := range pips {
		pip := pips[(off+i)%len(pips)]
		if pip.LinkPublicIpId == nil {
			log.V(3).Info("Found publicIp in pool", "publicIpId", pip.PublicIpId, "publicIp", pip.PublicIp)
			return &pip, nil
		}
	}
	return nil, ErrEmptyPool
}
