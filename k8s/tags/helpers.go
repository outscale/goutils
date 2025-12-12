package tags

import (
	"github.com/outscale/goutils/k8s/role"
	"github.com/outscale/goutils/sdk/tags"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

// GetClusterID fetches the cluster ID from tags.
func GetClusterID(t []osc.ResourceTag) string {
	suffix, _, _ := tags.HasPrefix(t, ClusterIDPrefix)
	return suffix
}

// ClusterIDKey returns the key for a cluster ID tag.
func ClusterIDKey(id string) string {
	return ClusterIDPrefix + id
}

// GetServiceName fetches the name of a service from tags.
func GetServiceName(t []osc.ResourceTag) string {
	return tags.Must(tags.GetValue(t, ServiceName))
}

// RoleKey returns the key for a role tag.
func RoleKey(role role.Role) string {
	return RolePrefix + string(role)
}

// HasRole checks if a resource has a role.
func HasRole(t []osc.ResourceTag, role role.Role) bool {
	return tags.Has(t, RoleKey(role))
}

// Wrappers for sdk/tags.

func Has(t []osc.ResourceTag, kv ...string) bool {
	return tags.Has(t, kv...)
}

func HasPrefix(t []osc.ResourceTag, prefix string) (suffix, value string, found bool) {
	return tags.HasPrefix(t, prefix)
}

func GetValue(t []osc.ResourceTag, k string) (string, bool) {
	return tags.GetValue(t, k)
}

func GetName(t []osc.ResourceTag) (string, bool) {
	return tags.GetName(t)
}

func Must(v string, _ bool) string {
	return tags.Must(v, true)
}
