package mocks_metadata

import (
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/outscale/goutils/sdk/metadata"
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
	mock(metadata.SubRegion, az)
}

func MockInstanceID(id string) {
	mock(metadata.InstanceID, id)
}
