package batch_test

import (
	"context"
	"errors"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/outscale/goutils/sdk/batch"
	"github.com/outscale/goutils/sdk/mock_osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestResourceWatcher_Volumes(t *testing.T) {
	t.Run("When concurrent calls are made, the right status is returned to the right volume", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		moskSDK := mock_osc.NewMockClient(mockCtrl)
		moskSDK.EXPECT().ReadVolumes(gomock.Any(), gomock.Cond(func(req osc.ReadVolumesRequest) bool {
			return len(*req.Filters.VolumeIds) == 4
		})).Return(&osc.ReadVolumesResponse{Volumes: &[]osc.Volume{
			{VolumeId: "id-creating", State: osc.VolumeStateCreating},
			{VolumeId: "id-available", State: osc.VolumeStateAvailable},
			{VolumeId: "id-in-use", State: osc.VolumeStateInUse},
			{VolumeId: "id-error", State: osc.VolumeStateError},
		}}, nil).MinTimes(1)

		rw := batch.NewVolumeBatcher(time.Second, moskSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		for _, state := range []osc.VolumeState{osc.VolumeStateCreating, osc.VolumeStateAvailable, osc.VolumeStateInUse, osc.VolumeStateError} {
			wg.Go(func() {
				time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
				v, err := rw.WaitUntil(ctx, "id-"+string(state), func(v *osc.Volume) (bool, error) {
					if v.State != state {
						return false, errors.New("invalid state")
					}
					return true, nil
				})
				require.NoError(t, err)
				assert.Equal(t, "id-"+string(state), v.VolumeId)
				assert.Equal(t, state, v.State)
			})
		}
		wg.Wait()
	})
	t.Run("Errors are returned", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		moskSDK := mock_osc.NewMockClient(mockCtrl)
		moskSDK.EXPECT().ReadVolumes(gomock.Any(), gomock.Any()).Return(&osc.ReadVolumesResponse{Volumes: &[]osc.Volume{
			{VolumeId: "id-error"},
		}}, nil).MinTimes(1)

		rw := batch.NewVolumeBatcher(time.Second, moskSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		_, err := rw.WaitUntil(ctx, "id-error", func(v *osc.Volume) (bool, error) {
			return false, errors.New("error")
		})
		require.Error(t, err)
	})
}

func TestResourceWatcher_Snapshots(t *testing.T) {
	t.Run("When concurrent calls are made, the right status is returned to the right snapshot", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		moskSDK := mock_osc.NewMockClient(mockCtrl)
		moskSDK.EXPECT().ReadSnapshots(gomock.Any(), gomock.Cond(func(req osc.ReadSnapshotsRequest) bool {
			return len(*req.Filters.SnapshotIds) == 4
		})).Return(&osc.ReadSnapshotsResponse{Snapshots: &[]osc.Snapshot{
			{SnapshotId: "id-completed", State: osc.SnapshotStateCompleted},
			{SnapshotId: "id-pending", State: osc.SnapshotStatePending},
			{SnapshotId: "id-deleting", State: osc.SnapshotStateDeleting},
			{SnapshotId: "id-error", State: osc.SnapshotStateError},
		}}, nil).MinTimes(1)

		rw := batch.NewSnapshotBatcher(time.Second, moskSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		for _, state := range []osc.SnapshotState{osc.SnapshotStateCompleted, osc.SnapshotStatePending, osc.SnapshotStateDeleting, osc.SnapshotStateError} {
			wg.Go(func() {
				time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
				s, err := rw.WaitUntil(ctx, "id-"+string(state), func(v *osc.Snapshot) (bool, error) {
					if v.State != state {
						return false, errors.New("invalid state")
					}
					return true, nil
				})
				require.NoError(t, err)
				assert.Equal(t, "id-"+string(state), s.SnapshotId)
				assert.Equal(t, state, s.State)
			})
		}
		wg.Wait()
	})
	t.Run("Errors are returned", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		moskSDK := mock_osc.NewMockClient(mockCtrl)
		moskSDK.EXPECT().ReadSnapshots(gomock.Any(), gomock.Any()).Return(&osc.ReadSnapshotsResponse{Snapshots: &[]osc.Snapshot{
			{SnapshotId: "id-error"},
		}}, nil).MinTimes(1)

		rw := batch.NewSnapshotBatcher(time.Second, moskSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		_, err := rw.WaitUntil(ctx, "id-error", func(v *osc.Snapshot) (bool, error) {
			return false, errors.New("error")
		})
		require.Error(t, err)
	})
}
