package mocks_metadata

import (
	"net/http"
	"path"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/outscale/goutils/sdk/metadata"
	"github.com/samber/lo"
)

var mockClient *http.Client

func Setup() {
	mockClient = &http.Client{Transport: http.DefaultTransport.(*http.Transport).Clone()}
	httpmock.ActivateNonDefault(mockClient)
	metadata.DefaultService = metadata.NewService(mockClient)
}

func Teardown() {
	if mockClient != nil {
		httpmock.DeactivateNonDefault(mockClient)
	}
}

func mock(k, v string) {
	httpmock.RegisterResponder(http.MethodGet, metadata.MetadataServer+k,
		httpmock.NewStringResponder(200, v))
}

func MockSubRegion(az string) {
	mock(metadata.Subregion, az)
}

func MockInstanceID(id string) {
	mock(metadata.InstanceID, id)
}

func MockDevideMappings(mappings map[string]string) {
	if mappings == nil {
		mappings = map[string]string{}
	}
	mappings["root"] = "/dev/sda1"
	mappings["ami"] = "/dev/sda1"
	for k, v := range mappings {
		mock(path.Join(metadata.DeviceMapping, k), v)
	}
	mock(metadata.DeviceMapping, strings.Join(lo.Keys(mappings), "\n"))
}
