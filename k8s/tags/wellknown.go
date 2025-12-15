/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package tags

const (
	// ClusterIDPrefix is the tag key prefix we use to differentiate multiple
	// logically independent clusters running in the same AZ.
	// The tag key = ClusterIDPrefix + clusterID
	// The tag value is an ownership value
	ClusterIDPrefix = "OscK8sClusterID/"

	// MainSGPrefix The main SG Tag (deprecated)
	// The tag key = OscK8sMainSG/clusterId
	MainSGPrefix = "OscK8sMainSG/"

	// RolePrefix is the prefix of tag key storing the role.
	RolePrefix = "OscK8sRole/"

	// ServiceName is the tag key storing the service name.
	ServiceName = "OscK8sService"

	// VmNodeName stores the node name.
	VmNodeName = "OscK8sNodeName"

	// PublicIPPool stores the name of a Public IP.
	PublicIPPool = "OscK8sIPPool"
)

// ResourceLifecycle is the cluster lifecycle state used in tagging
type ResourceLifecycle string

const (
	// ResourceLifecycleOwned is the value we use when tagging resources to indicate
	// that the resource is considered owned and managed by the cluster,
	// and in particular that the lifecycle is tied to the lifecycle of the cluster.
	ResourceLifecycleOwned = "owned"
	// ResourceLifecycleShared is the value we use when tagging resources to indicate
	// that the resource is shared between multiple clusters, and should not be destroyed
	// if the cluster is destroyed.
	ResourceLifecycleShared = "shared"
)
