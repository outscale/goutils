package mock_oks

import (
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
)

//go:generate mockgen -destination oks_mock.go -package mock_oks -source generate.go
type Client interface {
	oks.ClientInterface
}
