package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/goutils/sdk/generics/maps"
	"github.com/outscale/goutils/sdk/ptr"
	osc "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

func NewSnapshotBatcher(interval time.Duration, client osc.ClientInterface) *Batcher[osc.Snapshot] {
	return NewBatcher(interval, func(ctx context.Context, ids []string) (map[string]*osc.Snapshot, error) {
		req := osc.ReadSnapshotsRequest{
			Filters: &osc.FiltersSnapshot{
				SnapshotIds: &ids,
			},
			ResultsPerPage: ptr.To(len(ids)),
		}
		resp, err := client.ReadSnapshots(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("read snapshots: %w", err)
		}
		return maps.FromSlice(*resp.Snapshots, func(snap osc.Snapshot) (string, *osc.Snapshot) {
			return snap.SnapshotId, &snap
		}), nil
	})
}

func NewVolumeBatcher(interval time.Duration, client osc.ClientInterface) *Batcher[osc.Volume] {
	return NewBatcher(interval, func(ctx context.Context, ids []string) (map[string]*osc.Volume, error) {
		req := osc.ReadVolumesRequest{
			Filters: &osc.FiltersVolume{
				VolumeIds: &ids,
			},
			ResultsPerPage: ptr.To(len(ids)),
		}
		resp, err := client.ReadVolumes(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("read volumes: %w", err)
		}
		return maps.FromSlice(*resp.Volumes, func(vol osc.Volume) (string, *osc.Volume) {
			return vol.VolumeId, &vol
		}), nil
	})
}
