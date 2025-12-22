package mocks_osc

import (
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

//go:generate mockgen -destination mock_osc.go -package mocks_osc -source generate.go
type Client interface {
	osc.ClientInterface
}
