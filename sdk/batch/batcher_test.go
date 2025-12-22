/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package batch_test

import (
	"context"
	"errors"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/outscale/goutils/sdk/batch"
	"github.com/outscale/goutils/sdk/mocks_osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBatcherById_Volumes(t *testing.T) {
	t.Run("When concurrent calls are made, the right status is returned to the right volume", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadVolumes(gomock.Any(), gomock.Cond(func(req osc.ReadVolumesRequest) bool {
			return len(*req.Filters.VolumeIds) == 4
		})).Return(&osc.ReadVolumesResponse{Volumes: &[]osc.Volume{
			{VolumeId: "id-creating", State: osc.VolumeStateCreating},
			{VolumeId: "id-available", State: osc.VolumeStateAvailable},
			{VolumeId: "id-in-use", State: osc.VolumeStateInUse},
			{VolumeId: "id-error", State: osc.VolumeStateError},
		}}, nil).MinTimes(1)

		rw := batch.NewVolumeBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 99*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		for range 2 {
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
		}
		wg.Wait()
	})
	t.Run("ErrNotFound is returned if a volume does not exist anymore", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadVolumes(gomock.Any(), gomock.Any()).Return(&osc.ReadVolumesResponse{Volumes: &[]osc.Volume{}}, nil).MinTimes(1)

		rw := batch.NewVolumeBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 99*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		wg.Go(func() {
			_, err := rw.Read(ctx, "id-foo")
			require.ErrorIs(t, err, batch.ErrNotFound)
		})
		wg.Wait()
	})
	t.Run("Errors are returned", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadVolumes(gomock.Any(), gomock.Any()).Return(&osc.ReadVolumesResponse{Volumes: &[]osc.Volume{
			{VolumeId: "id-error"},
		}}, nil).MinTimes(1)

		rw := batch.NewVolumeBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		_, err := rw.WaitUntil(ctx, "id-error", func(v *osc.Volume) (bool, error) {
			return false, errors.New("error")
		})
		require.Error(t, err)
	})
}

func TestBatcherSameQuery_Volumes(t *testing.T) {
	t.Run("When concurrent calls are made, a single query is made", func(t *testing.T) {
		req := osc.ReadVolumesRequest{Filters: &osc.FiltersVolume{
			VolumeStates: &[]osc.VolumeState{osc.VolumeStateCreating, osc.VolumeStateAvailable},
		}}
		res := &[]osc.Volume{
			{VolumeId: "id-foo", State: osc.VolumeStateCreating},
			{VolumeId: "id-bar", State: osc.VolumeStateAvailable},
			{VolumeId: "id-baz", State: osc.VolumeStateAvailable},
		}
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadVolumes(gomock.Any(), gomock.Eq(req)).Return(&osc.ReadVolumesResponse{Volumes: res}, nil)
		rw := batch.NewVolumeBatcherSameQuery(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		for range 4 {
			wg.Go(func() {
				vs, err := rw.Read(ctx, req)
				require.NoError(t, err)
				require.NotNil(t, vs.Volumes)
				require.Len(t, *vs.Volumes, 3)
			})
		}
		wg.Wait()
	})
}

func TestBatcherById_Snapshots(t *testing.T) {
	t.Run("When concurrent calls are made, the right status is returned to the right snapshot", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadSnapshots(gomock.Any(), gomock.Cond(func(req osc.ReadSnapshotsRequest) bool {
			return len(*req.Filters.SnapshotIds) == 4
		})).Return(&osc.ReadSnapshotsResponse{Snapshots: &[]osc.Snapshot{
			{SnapshotId: "id-completed", State: osc.SnapshotStateCompleted},
			{SnapshotId: "id-pending", State: osc.SnapshotStatePending},
			{SnapshotId: "id-deleting", State: osc.SnapshotStateDeleting},
			{SnapshotId: "id-error", State: osc.SnapshotStateError},
		}}, nil).MinTimes(1)

		rw := batch.NewSnapshotBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		for range 2 {
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
		}
		wg.Wait()
	})
	t.Run("ErrNotFound is returned if a snapshot does not exist anymore", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadSnapshots(gomock.Any(), gomock.Any()).Return(&osc.ReadSnapshotsResponse{Snapshots: &[]osc.Snapshot{}}, nil).MinTimes(1)

		rw := batch.NewSnapshotBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 99*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		wg.Go(func() {
			_, err := rw.Read(ctx, "id-foo")
			require.ErrorIs(t, err, batch.ErrNotFound)
		})
		wg.Wait()
	})
	t.Run("Errors are returned", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadSnapshots(gomock.Any(), gomock.Any()).Return(&osc.ReadSnapshotsResponse{Snapshots: &[]osc.Snapshot{
			{SnapshotId: "id-error"},
		}}, nil).MinTimes(1)

		rw := batch.NewSnapshotBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		_, err := rw.WaitUntil(ctx, "id-error", func(v *osc.Snapshot) (bool, error) {
			return false, errors.New("error")
		})
		require.Error(t, err)
	})
}

func TestBatcherSameQuery_Snapshots(t *testing.T) {
	t.Run("When concurrent calls are made, a single query is made", func(t *testing.T) {
		req := osc.ReadSnapshotsRequest{Filters: &osc.FiltersSnapshot{
			States: &[]osc.SnapshotState{osc.SnapshotStateCompleted, osc.SnapshotStatePending},
		}}
		res := &[]osc.Snapshot{
			{SnapshotId: "id-foo", State: osc.SnapshotStateCompleted},
			{SnapshotId: "id-bar", State: osc.SnapshotStatePending},
			{SnapshotId: "id-baz", State: osc.SnapshotStatePending},
		}
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadSnapshots(gomock.Any(), gomock.Eq(req)).Return(&osc.ReadSnapshotsResponse{Snapshots: res}, nil)
		rw := batch.NewSnapshotBatcherSameQuery(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		for range 4 {
			wg.Go(func() {
				vs, err := rw.Read(ctx, req)
				require.NoError(t, err)
				require.NotNil(t, vs.Snapshots)
				require.Len(t, *vs.Snapshots, 3)
			})
		}
		wg.Wait()
	})
}

func TestBatcherById_Vms(t *testing.T) {
	t.Run("When concurrent calls are made, the right status is returned to the right Vm", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadVms(gomock.Any(), gomock.Cond(func(req osc.ReadVmsRequest) bool {
			return len(*req.Filters.VmIds) == 4
		})).Return(&osc.ReadVmsResponse{Vms: &[]osc.Vm{
			{VmId: "id-pending", State: osc.VmStatePending},
			{VmId: "id-running", State: osc.VmStateRunning},
			{VmId: "id-stopped", State: osc.VmStateStopped},
			{VmId: "id-terminated", State: osc.VmStateTerminated},
		}}, nil).MinTimes(1)

		rw := batch.NewVmBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 99*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		for range 2 {
			for _, state := range []osc.VmState{osc.VmStatePending, osc.VmStateRunning, osc.VmStateStopped, osc.VmStateTerminated} {
				wg.Go(func() {
					time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
					v, err := rw.WaitUntil(ctx, "id-"+string(state), func(v *osc.Vm) (bool, error) {
						if v.State != state {
							return false, errors.New("invalid state")
						}
						return true, nil
					})
					require.NoError(t, err)
					assert.Equal(t, "id-"+string(state), v.VmId)
					assert.Equal(t, state, v.State)
				})
			}
		}
		wg.Wait()
	})
	t.Run("ErrNotFound is returned if a Vm does not exist anymore", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadVms(gomock.Any(), gomock.Any()).Return(&osc.ReadVmsResponse{Vms: &[]osc.Vm{}}, nil).MinTimes(1)

		rw := batch.NewVmBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 99*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		wg.Go(func() {
			_, err := rw.Read(ctx, "id-foo")
			require.ErrorIs(t, err, batch.ErrNotFound)
		})
		wg.Wait()
	})
	t.Run("Errors are returned", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadVms(gomock.Any(), gomock.Any()).Return(&osc.ReadVmsResponse{Vms: &[]osc.Vm{
			{VmId: "id-error"},
		}}, nil).MinTimes(1)

		rw := batch.NewVmBatcherByID(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		_, err := rw.WaitUntil(ctx, "id-error", func(v *osc.Vm) (bool, error) {
			return false, errors.New("error")
		})
		require.Error(t, err)
	})
}

func TestBatcherSameQuery_Vms(t *testing.T) {
	t.Run("When concurrent calls are made, a single query is made", func(t *testing.T) {
		req := osc.ReadVmsRequest{Filters: &osc.FiltersVm{
			VmStateNames: &[]osc.VmState{osc.VmStatePending, osc.VmStateRunning},
		}}
		res := &[]osc.Vm{
			{VmId: "id-foo", State: osc.VmStatePending},
			{VmId: "id-bar", State: osc.VmStateRunning},
			{VmId: "id-baz", State: osc.VmStateRunning},
		}
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadVms(gomock.Any(), gomock.Eq(req)).Return(&osc.ReadVmsResponse{Vms: res}, nil)
		rw := batch.NewVmBatcherSameQuery(time.Second, mockSDK)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go rw.Run(ctx)

		wg := sync.WaitGroup{}
		for range 4 {
			wg.Go(func() {
				vs, err := rw.Read(ctx, req)
				require.NoError(t, err)
				require.NotNil(t, vs.Vms)
				require.Len(t, *vs.Vms, 3)
			})
		}
		wg.Wait()
	})
}
