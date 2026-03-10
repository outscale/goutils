/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	osc "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

func NewSecurityGroupBatcherByID(interval time.Duration, client osc.ClientInterface) *BatcherByID[osc.SecurityGroup] {
	return NewBatcherByID(interval, func(ctx context.Context, ids []string) (resultFn[string, osc.SecurityGroup], error) {
		req := osc.ReadSecurityGroupsRequest{
			Filters: &osc.FiltersSecurityGroup{
				SecurityGroupIds: &ids,
			},
			ResultsPerPage: ptr.To(len(ids)),
		}
		resp, err := client.ReadSecurityGroups(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("read security groups: %w", err)
		}
		res := *resp.SecurityGroups
		return func(query string) (*osc.SecurityGroup, bool) {
			for i := range res {
				if res[i].SecurityGroupId == query {
					return &res[i], true
				}
			}
			return nil, false
		}, nil
	})
}

func NewSecurityGroupBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *BatcherSameQuery[osc.ReadSecurityGroupsRequest, osc.ReadSecurityGroupsResponse] {
	return NewBatcherSameQuery(interval, func(ctx context.Context, queries []osc.ReadSecurityGroupsRequest) (resultFn[osc.ReadSecurityGroupsRequest, osc.ReadSecurityGroupsResponse], error) {
		resp, err := client.ReadSecurityGroups(ctx, queries[0])
		if err != nil {
			return nil, fmt.Errorf("read security groups: %w", err)
		}
		return func(_ osc.ReadSecurityGroupsRequest) (*osc.ReadSecurityGroupsResponse, bool) {
			return resp, true
		}, nil
	})
}

func NewSubnetBatcherByID(interval time.Duration, client osc.ClientInterface) *BatcherByID[osc.Subnet] {
	return NewBatcherByID(interval, func(ctx context.Context, ids []string) (resultFn[string, osc.Subnet], error) {
		req := osc.ReadSubnetsRequest{
			Filters: &osc.FiltersSubnet{
				SubnetIds: &ids,
			},
			ResultsPerPage: ptr.To(len(ids)),
		}
		resp, err := client.ReadSubnets(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("read subnets: %w", err)
		}
		res := *resp.Subnets
		return func(query string) (*osc.Subnet, bool) {
			for i := range res {
				if res[i].SubnetId == query {
					return &res[i], true
				}
			}
			return nil, false
		}, nil
	})
}

func NewSubnetBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *BatcherSameQuery[osc.ReadSubnetsRequest, osc.ReadSubnetsResponse] {
	return NewBatcherSameQuery(interval, func(ctx context.Context, queries []osc.ReadSubnetsRequest) (resultFn[osc.ReadSubnetsRequest, osc.ReadSubnetsResponse], error) {
		resp, err := client.ReadSubnets(ctx, queries[0])
		if err != nil {
			return nil, fmt.Errorf("read subnets: %w", err)
		}
		return func(_ osc.ReadSubnetsRequest) (*osc.ReadSubnetsResponse, bool) {
			return resp, true
		}, nil
	})
}

func NewNetBatcherByID(interval time.Duration, client osc.ClientInterface) *BatcherByID[osc.Net] {
	return NewBatcherByID(interval, func(ctx context.Context, ids []string) (resultFn[string, osc.Net], error) {
		req := osc.ReadNetsRequest{
			Filters: &osc.FiltersNet{
				NetIds: &ids,
			},
			ResultsPerPage: ptr.To(len(ids)),
		}
		resp, err := client.ReadNets(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("read nets: %w", err)
		}
		res := *resp.Nets
		return func(query string) (*osc.Net, bool) {
			for i := range res {
				if res[i].NetId == query {
					return &res[i], true
				}
			}
			return nil, false
		}, nil
	})
}

func NewNetBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *BatcherSameQuery[osc.ReadNetsRequest, osc.ReadNetsResponse] {
	return NewBatcherSameQuery(interval, func(ctx context.Context, queries []osc.ReadNetsRequest) (resultFn[osc.ReadNetsRequest, osc.ReadNetsResponse], error) {
		resp, err := client.ReadNets(ctx, queries[0])
		if err != nil {
			return nil, fmt.Errorf("read nets: %w", err)
		}
		return func(_ osc.ReadNetsRequest) (*osc.ReadNetsResponse, bool) {
			return resp, true
		}, nil
	})
}

func NewSnapshotBatcherByID(interval time.Duration, client osc.ClientInterface) *BatcherByID[osc.Snapshot] {
	return NewBatcherByID(interval, func(ctx context.Context, ids []string) (resultFn[string, osc.Snapshot], error) {
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
		res := *resp.Snapshots
		return func(query string) (*osc.Snapshot, bool) {
			for i := range res {
				if res[i].SnapshotId == query {
					return &res[i], true
				}
			}
			return nil, false
		}, nil
	})
}

func NewSnapshotBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *BatcherSameQuery[osc.ReadSnapshotsRequest, osc.ReadSnapshotsResponse] {
	return NewBatcherSameQuery(interval, func(ctx context.Context, queries []osc.ReadSnapshotsRequest) (resultFn[osc.ReadSnapshotsRequest, osc.ReadSnapshotsResponse], error) {
		resp, err := client.ReadSnapshots(ctx, queries[0])
		if err != nil {
			return nil, fmt.Errorf("read snapshots: %w", err)
		}
		return func(_ osc.ReadSnapshotsRequest) (*osc.ReadSnapshotsResponse, bool) {
			return resp, true
		}, nil
	})
}

func NewVolumeBatcherByID(interval time.Duration, client osc.ClientInterface) *BatcherByID[osc.Volume] {
	return NewBatcherByID(interval, func(ctx context.Context, ids []string) (resultFn[string, osc.Volume], error) {
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
		res := *resp.Volumes
		return func(query string) (*osc.Volume, bool) {
			for i := range res {
				if res[i].VolumeId == query {
					return &res[i], true
				}
			}
			return nil, false
		}, nil
	})
}

func NewVolumeBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *BatcherSameQuery[osc.ReadVolumesRequest, osc.ReadVolumesResponse] {
	return NewBatcherSameQuery(interval, func(ctx context.Context, queries []osc.ReadVolumesRequest) (resultFn[osc.ReadVolumesRequest, osc.ReadVolumesResponse], error) {
		resp, err := client.ReadVolumes(ctx, queries[0])
		if err != nil {
			return nil, fmt.Errorf("read volumes: %w", err)
		}
		return func(_ osc.ReadVolumesRequest) (*osc.ReadVolumesResponse, bool) {
			return resp, true
		}, nil
	})
}

func NewVmBatcherByID(interval time.Duration, client osc.ClientInterface) *BatcherByID[osc.Vm] {
	return NewBatcherByID(interval, func(ctx context.Context, ids []string) (resultFn[string, osc.Vm], error) {
		req := osc.ReadVmsRequest{
			Filters: &osc.FiltersVm{
				VmIds: &ids,
			},
			ResultsPerPage: ptr.To(len(ids)),
		}
		resp, err := client.ReadVms(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("read vms: %w", err)
		}
		res := *resp.Vms
		return func(query string) (*osc.Vm, bool) {
			for i := range res {
				if res[i].VmId == query {
					return &res[i], true
				}
			}
			return nil, false
		}, nil
	})
}

func NewVmBatcherSameQuery(interval time.Duration, client osc.ClientInterface) *BatcherSameQuery[osc.ReadVmsRequest, osc.ReadVmsResponse] {
	return NewBatcherSameQuery(interval, func(ctx context.Context, queries []osc.ReadVmsRequest) (resultFn[osc.ReadVmsRequest, osc.ReadVmsResponse], error) {
		resp, err := client.ReadVms(ctx, queries[0])
		if err != nil {
			return nil, fmt.Errorf("read vms: %w", err)
		}
		return func(_ osc.ReadVmsRequest) (*osc.ReadVmsResponse, bool) {
			return resp, true
		}, nil
	})
}
