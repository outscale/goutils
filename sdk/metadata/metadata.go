/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package metadata

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

var DefaultService = NewService(http.DefaultClient)

const (
	MetadataServer = "http://169.254.169.254/latest/meta-data/"

	SubRegion  = "placement/availability-zone"
	InstanceID = "instance-id"
)

// Metadata is the metadata returned by the metadata server.
type Metadata struct {
	InstanceID string
	Region     string
	SubRegion  string
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

func (s *Service) fetchSubRegion(ctx context.Context) (string, error) {
	return s.fetch(ctx, SubRegion)
}

func (s *Service) fetchInstanceID(ctx context.Context) (string, error) {
	return s.fetch(ctx, InstanceID)
}

// Fetch fetches metadata from the metadata server.
func (s *Service) Fetch(ctx context.Context) (Metadata, error) {
	instanceID, err := s.fetchInstanceID(ctx)
	if err != nil {
		return Metadata{}, err
	}
	availabilityZone, err := s.fetchSubRegion(ctx)
	if err != nil {
		return Metadata{}, err
	}
	region := availabilityZone[0 : len(availabilityZone)-1]
	return Metadata{
		InstanceID: instanceID,
		Region:     region,
		SubRegion:  availabilityZone,
	}, nil
}

// Fetch fetches metadata from the metadata server.
func Fetch(ctx context.Context) (Metadata, error) {
	return DefaultService.Fetch(ctx)
}

// GetSubRegion fetches the sub region from the metadata server.
func GetSubRegion(ctx context.Context) (string, error) {
	return DefaultService.fetchSubRegion(ctx)
}

// GetRegion fetches the region from the metadata server.
func GetRegion(ctx context.Context) (string, error) {
	subregion, err := DefaultService.fetchSubRegion(ctx)
	if err != nil {
		return "", err
	}
	return subregion[0 : len(subregion)-1], nil
}

// GetInstanceID fetches the instance ID from the metadata server.
func GetInstanceID(ctx context.Context) (string, error) {
	return DefaultService.fetchInstanceID(ctx)
}
