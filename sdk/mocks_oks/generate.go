package mocks_oks

import (
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
)

//go:generate mockgen -destination mock_oks.go -package mocks_oks -source generate.go
type Client interface {
	oks.ClientInterface
}
