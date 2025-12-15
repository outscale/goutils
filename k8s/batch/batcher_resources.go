/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package batch

import (
	"time"

	k8slog "github.com/outscale/goutils/k8s/log"
	"github.com/outscale/goutils/sdk/batch"
	"github.com/outscale/goutils/sdk/log"
	osc "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

func init() {
	log.Default = k8slog.Logger{}
}

func NewSnapshotBatcherByID(interval time.Duration, client osc.ClientInterface) *batch.BatcherByID[osc.Snapshot] {
	return batch.NewSnapshotBatcherByID(interval, client)
}

func NewSnapshotBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *batch.BatcherSameQuery[osc.ReadSnapshotsRequest, osc.ReadSnapshotsResponse] {
	return batch.NewSnapshotBatcherSameQuery(interval, client)
}

func NewVolumeBatcherByID(interval time.Duration, client osc.ClientInterface) *batch.BatcherByID[osc.Volume] {
	return batch.NewVolumeBatcherByID(interval, client)
}

func NewVolumeBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *batch.BatcherSameQuery[osc.ReadVolumesRequest, osc.ReadVolumesResponse] {
	return batch.NewVolumeBatcherSameQuery(interval, client)
}

func NewVmBatcherByID(interval time.Duration, client osc.ClientInterface) *batch.BatcherByID[osc.Vm] {
	return batch.NewVmBatcherByID(interval, client)
}

func NewVmBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *batch.BatcherSameQuery[osc.ReadVmsRequest, osc.ReadVmsResponse] {
	return batch.NewVmBatcherSameQuery(interval, client)
}
