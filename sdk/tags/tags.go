package tags

import (
	"slices"
	"strings"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

// Has checks if a tag slice contains a tag.
// If called with 1 string param, it checks if a key exists.
// If called with 2 string params, it checks if a key exists with the exact value.
func Has(tags []osc.ResourceTag, kv ...string) bool {
	if len(kv) == 0 {
		return true
	}
	return slices.ContainsFunc(tags, func(t osc.ResourceTag) bool {
		return kv[0] == t.Key && (len(kv) == 1 || kv[1] == t.Value)
	})
}

// HasPrefix checks if a tag slice contains a tag key with the specified prefix.
// It returns the suffix of the key found and the tag value.
func HasPrefix(tags []osc.ResourceTag, prefix string) (suffix, value string, found bool) {
	for _, t := range tags {
		if strings.HasPrefix(t.Key, prefix) {
			return strings.TrimPrefix(t.Key, prefix), t.Value, true
		}
	}
	return "", "", false
}
