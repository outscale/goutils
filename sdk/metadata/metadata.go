/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package metadata

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

var DefaultService = NewService(http.DefaultClient)

const (
	MetadataServer = "http://169.254.169.254/latest/meta-data/"

	Hostname         = "hostname"
	Subregion        = "placement/availability-zone"
	PlacementServer  = "placement/server"
	PlacementCluster = "placement/cluster"
	InstanceID       = "instance-id"
	OMIID            = "ami-id"
	InstanceType     = "instance-type"
	DeviceMapping    = "block-device-mapping"
	MAC              = "mac"
	Tags             = "tags"
)

func getRegion(az string) string {
	if len(az) <= 1 {
		return ""
	}
	return az[0 : len(az)-1]
}

type Placement struct {
	Subregion string `json:"availability_zone"`
	Cluster   string `json:"cluster"`
	Server    string `json:"server"`
}

func (p *Placement) GetRegion() string {
	return getRegion(p.Subregion)
}

// Metadata is the metadata returned by the metadata server.
type Metadata struct {
	Hostname      string            `json:"hostname"`
	InstanceID    string            `json:"instance_id"`
	InstanceType  string            `json:"instance_type"`
	OMIID         string            `json:"ami_id"`
	MAC           string            `json:"mac"`
	Placement     Placement         `json:"placement"`
	DeviceMapping map[string]string `json:"block_device_mapping"`
	Tags          map[string]string `json:"tags"`
}

// Service is a metadata service.
type Service struct {
	client *http.Client
}

// NewService builds a metadata service.
func NewService(client *http.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) fetch(ctx context.Context, path string) (res string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, MetadataServer+path, nil)
	if err != nil {
		return "", fmt.Errorf("get metadata: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("get metadata: %w", err)
	}
	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			err = cerr
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get metadata: %v returned %s", req.URL, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("get metadata: %w", err)
	}
	return string(body), nil
}

func (s *Service) GetHostname(ctx context.Context) (string, error) {
	return s.fetch(ctx, Hostname)
}

func (s *Service) GetSubregion(ctx context.Context) (string, error) {
	return s.fetch(ctx, Subregion)
}

func (s *Service) GetRegion(ctx context.Context) (string, error) {
	subregion, err := s.GetSubregion(ctx)
	if err != nil {
		return "", err
	}
	return getRegion(subregion), nil
}

func (s *Service) GetInstanceID(ctx context.Context) (string, error) {
	return s.fetch(ctx, InstanceID)
}

func (s *Service) GetOMIID(ctx context.Context) (string, error) {
	return s.fetch(ctx, OMIID)
}

func (s *Service) GetInstanceType(ctx context.Context) (string, error) {
	return s.fetch(ctx, InstanceType)
}

func (s *Service) GetMAC(ctx context.Context) (string, error) {
	return s.fetch(ctx, MAC)
}

func (s *Service) GetPlacementCluster(ctx context.Context) (string, error) {
	return s.fetch(ctx, PlacementCluster)
}

func (s *Service) GetPlacementServer(ctx context.Context) (string, error) {
	return s.fetch(ctx, PlacementServer)
}

func (s *Service) GetDeviceMappings(ctx context.Context) (map[string]string, error) {
	return s.fetchKeyValue(ctx, DeviceMapping)
}

func (s *Service) GetTags(ctx context.Context) (map[string]string, error) {
	return s.fetchKeyValue(ctx, Tags)
}

func (s *Service) fetchKeyValue(ctx context.Context, p string) (map[string]string, error) {
	res, err := s.fetch(ctx, p)
	if err != nil {
		return nil, err
	}
	kv := make(map[string]string)
	scan := bufio.NewScanner(strings.NewReader(res))
	for scan.Scan() {
		res, err := s.fetch(ctx, path.Join(p, scan.Text()))
		if err != nil {
			return nil, err
		}
		kv[scan.Text()] = res
	}
	if scan.Err() != nil {
		return nil, scan.Err()
	}
	return kv, nil
}

// Fetch fetches metadata from the metadata server.
func (s *Service) Fetch(ctx context.Context) (Metadata, error) {
	hostname, err := s.GetHostname(ctx)
	if err != nil {
		return Metadata{}, err
	}
	instanceID, err := s.GetInstanceID(ctx)
	if err != nil {
		return Metadata{}, err
	}
	az, err := s.GetSubregion(ctx)
	if err != nil {
		return Metadata{}, err
	}
	omi, err := s.GetOMIID(ctx)
	if err != nil {
		return Metadata{}, err
	}
	instanceType, err := s.GetInstanceType(ctx)
	if err != nil {
		return Metadata{}, err
	}
	mac, err := s.GetMAC(ctx)
	if err != nil {
		return Metadata{}, err
	}
	cluster, err := s.GetPlacementCluster(ctx)
	if err != nil {
		return Metadata{}, err
	}
	server, err := s.GetPlacementServer(ctx)
	if err != nil {
		return Metadata{}, err
	}
	mapping, err := s.GetDeviceMappings(ctx)
	if err != nil {
		return Metadata{}, err
	}
	tags, err := s.GetTags(ctx)
	if err != nil {
		return Metadata{}, err
	}
	return Metadata{
		Hostname:     hostname,
		InstanceID:   instanceID,
		OMIID:        omi,
		InstanceType: instanceType,
		MAC:          mac,
		Placement: Placement{
			Subregion: az,
			Cluster:   cluster,
			Server:    server,
		},
		DeviceMapping: mapping,
		Tags:          tags,
	}, nil
}

// Fetch fetches metadata from the metadata server.
func Fetch(ctx context.Context) (Metadata, error) {
	return DefaultService.Fetch(ctx)
}

// GetHostname fetches the hostname from the metadata server.
func GetHostname(ctx context.Context) (string, error) {
	return DefaultService.GetHostname(ctx)
}

// GetSubregion fetches the sub region from the metadata server.
func GetSubregion(ctx context.Context) (string, error) {
	return DefaultService.GetSubregion(ctx)
}

// GetRegion fetches the region from the metadata server.
func GetRegion(ctx context.Context) (string, error) {
	return DefaultService.GetRegion(ctx)
}

// GetInstanceID fetches the instance ID from the metadata server.
func GetInstanceID(ctx context.Context) (string, error) {
	return DefaultService.GetInstanceID(ctx)
}

// GetInstanceType fetches the instance type from the metadata server.
func GetInstanceType(ctx context.Context) (string, error) {
	return DefaultService.GetInstanceType(ctx)
}

// GetOMIID fetches the OMI ID from the metadata server.
func GetOMIID(ctx context.Context) (string, error) {
	return DefaultService.GetOMIID(ctx)
}

// GetMAC fetches the MAC from the metadata server.
func GetMAC(ctx context.Context) (string, error) {
	return DefaultService.GetMAC(ctx)
}

// GetPlacementCluster fetches the cluster where the VM is located.
func GetPlacementCluster(ctx context.Context) (string, error) {
	return DefaultService.GetPlacementCluster(ctx)
}

// GetPlacementServer fetches the physical server where the VM is located.
func GetPlacementServer(ctx context.Context) (string, error) {
	return DefaultService.GetPlacementServer(ctx)
}

// GetDeviceMappings fetches the device mapping of the VM.
func GetDeviceMappings(ctx context.Context) (map[string]string, error) {
	return DefaultService.GetDeviceMappings(ctx)
}

// GetTags fetches the tags of the VM.
func GetTags(ctx context.Context) (map[string]string, error) {
	return DefaultService.GetTags(ctx)
}
