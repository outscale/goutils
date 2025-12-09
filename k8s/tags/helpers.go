package tags

import (
	"github.com/outscale/goutils/sdk/tags"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

var (
	Has       = tags.Has
	HasPrefix = tags.HasPrefix
	GetValue  = tags.GetValue
	GetName   = tags.GetName
	Must      = tags.Must
)

// GetClusterID fetches the cluster ID from tags.
func GetClusterID(t []osc.ResourceTag) string {
	suffix, _, _ := HasPrefix(t, ClusterIDPrefix)
	return suffix
}

// ClusterIDKey returns the key for a cluster ID tag.
func ClusterIDKey(id string) string {
	return ClusterIDPrefix + id
}

func GetServiceName(t []osc.ResourceTag) string {
	return tags.Must(GetValue(t, ServiceName))
}
