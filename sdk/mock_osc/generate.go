package mock_osc

import (
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

//go:generate mockgen -destination osc_mock.go -package mock_osc -source generate.go
type Client interface {
	osc.ClientInterface
}
