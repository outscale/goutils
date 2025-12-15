package sdk_test

import (
	"testing"

	"github.com/outscale/goutils/k8s/sdk"
	"github.com/outscale/goutils/k8s/tags"
	"github.com/outscale/goutils/sdk/mocks_osc"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAllocateIPFromPool(t *testing.T) {
	t.Run("An IP is allocated", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadPublicIps(gomock.Any(), gomock.Eq(osc.ReadPublicIpsRequest{
			Filters: &osc.FiltersPublicIp{
				TagKeys:   &[]string{tags.PublicIPPool},
				TagValues: &[]string{"foo"},
			},
		})).Return(&osc.ReadPublicIpsResponse{PublicIps: &[]osc.PublicIp{
			{PublicIpId: "bar", LinkPublicIpId: ptr.To("bar")},
			{PublicIpId: "baz"},
		}}, nil)

		ip, err := sdk.AllocateIPFromPool(t.Context(), "foo", mockSDK)
		require.NoError(t, err)
		assert.Equal(t, "baz", ip.PublicIpId)
	})
	t.Run("ErrEmptyPool is returned if no IP is found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadPublicIps(gomock.Any(), gomock.Eq(osc.ReadPublicIpsRequest{
			Filters: &osc.FiltersPublicIp{
				TagKeys:   &[]string{tags.PublicIPPool},
				TagValues: &[]string{"foo"},
			},
		})).Return(&osc.ReadPublicIpsResponse{PublicIps: &[]osc.PublicIp{}}, nil)

		_, err := sdk.AllocateIPFromPool(t.Context(), "foo", mockSDK)
		require.ErrorIs(t, err, sdk.ErrEmptyPool)
	})
	t.Run("ErrEmptyPool is returned if all IP are already allocated", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockSDK := mocks_osc.NewMockClient(mockCtrl)
		mockSDK.EXPECT().ReadPublicIps(gomock.Any(), gomock.Eq(osc.ReadPublicIpsRequest{
			Filters: &osc.FiltersPublicIp{
				TagKeys:   &[]string{tags.PublicIPPool},
				TagValues: &[]string{"foo"},
			},
		})).Return(&osc.ReadPublicIpsResponse{PublicIps: &[]osc.PublicIp{
			{PublicIpId: "bar", LinkPublicIpId: ptr.To("bar")},
			{PublicIpId: "baz", LinkPublicIpId: ptr.To("baz")},
		}}, nil)

		_, err := sdk.AllocateIPFromPool(t.Context(), "foo", mockSDK)
		require.ErrorIs(t, err, sdk.ErrEmptyPool)
	})
}
