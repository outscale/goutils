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

const metadataServer = "http://169.254.169.254/latest/meta-data/"

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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, metadataServer+path, nil)
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
	return s.fetch(ctx, "placement/availability-zone")
}

func (s *Service) fetchInstanceID(ctx context.Context) (string, error) {
	return s.fetch(ctx, "instance-id")
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
	svc := NewService(http.DefaultClient)
	return svc.Fetch(ctx)
}

// GetRegion fetches the region from the metadata server.
func GetRegion(ctx context.Context) (string, error) {
	svc := NewService(http.DefaultClient)
	subregion, err := svc.fetchSubRegion(ctx)
	if err != nil {
		return "", err
	}
	return subregion[0 : len(subregion)-1], nil
}
