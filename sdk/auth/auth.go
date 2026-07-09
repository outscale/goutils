package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

// ErrInvalidCredentials is returned by CheckCredentials if credentials are invalid.
var ErrInvalidCredentials = errors.New("authentication error (invalid credentials ?)")

// CheckCredentials checks if credentials are valid, and returns ErrInvalidCredentials if not.
func CheckCredentials(ctx context.Context, client osc.ClientInterface) error {
	_, err := client.ReadVms(ctx, osc.ReadVmsRequest{DryRun: ptr.To(true)})
	if err == nil {
		return nil
	}
	if osc.IsAuthError(err) {
		err = ErrInvalidCredentials
	}
	return fmt.Errorf("check credentials: %w", err)
}
