package metadata

import (
	"context"
	"fmt"

	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
)

// SetProfileDefaults is a profile option that fetches the region from the metadata server if not already set.
func SetProfileDefaults(ctx context.Context) profile.Option {
	return func(prof *profile.Profile) error {
		if prof.Region == "" {
			var err error
			prof.Region, err = GetRegion(ctx)
			if err != nil {
				return fmt.Errorf("unable to fetch metadata: %w", err)
			}
		}
		return nil
	}
}
