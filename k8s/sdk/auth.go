/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package sdk

import (
	"context"
	"errors"
	"fmt"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"k8s.io/klog/v2"
)

// ErrInvalidCredentials is returned by CheckCredentials if credentials are invalid.
var ErrInvalidCredentials = errors.New("authentication error (invalid credentials ?)")

// CheckCredentials checks if credentials are valid, and returns ErrInvalidCredentials if not.
// It uses ReadVms, as it is already used by CCM/CSI/CAPOSC and do not require a EIM policy change.
func CheckCredentials(ctx context.Context, client osc.ClientInterface) error {
	klog.FromContext(ctx).V(5).Info("Checking credentials")
	_, err := client.ReadVms(ctx, osc.ReadVmsRequest{DryRun: ptr.To(true)})
	if err == nil {
		return nil
	}
	if osc.IsAuthError(err) {
		err = ErrInvalidCredentials
	}
	return fmt.Errorf("check credentials: %w", err)
}
