//nolint:modernize
package v1beta

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPPool defines the model for a IPPool.
// +genclient
// +genclient:nonNamespaced
// +groupName=oks.dev
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IPPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   IPPoolSpec   `json:"spec"`
	Status IPPoolStatus `json:"status,omitempty"`
}

// IPPoolList defines model for list of IPPools.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IPPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []IPPool `json:"items"`
}

// NetPeering defines the model for a NetPeering.
// +genclient
// +genclient:nonNamespaced
// +groupName=oks.dev
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NetPeering struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   NetPeeringSpec   `json:"spec"`
	Status NetPeeringStatus `json:"status,omitempty"`
}

// NetPeeringList defines model for list of NetPeerings.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NetPeeringList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []NetPeering `json:"items"`
}

// NetPeeringRequest defines the model for a NetPeeringRequest.
// +genclient
// +genclient:nonNamespaced
// +groupName=oks.dev
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NetPeeringRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   NetPeeringRequestSpec `json:"spec"`
	Status NetPeeringStatus      `json:"status,omitempty"`
}

// NetPeeringRequestList defines model for list of NetPeeringRequests.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NetPeeringRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []NetPeeringRequest `json:"items"`
}

// NetPeeringAcceptance defines the model for a NetPeeringAcceptance.
// +genclient
// +genclient:nonNamespaced
// +groupName=oks.dev
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NetPeeringAcceptance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   NetPeeringSpec   `json:"spec"`
	Status NetPeeringStatus `json:"status,omitempty"`
}

// NetPeeringAcceptanceList defines model for list of NetPeeringAcceptances.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NetPeeringAcceptanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []NetPeeringAcceptance `json:"items"`
}

// OOSAccess defines the model for a OOSAccess.
// +genclient
// +genclient:nonNamespaced
// +groupName=oks.dev
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type OOSAccess struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   OOSAccessSpec   `json:"spec"`
	Status OOSAccessStatus `json:"status,omitempty"`
}

// OOSAccessList defines model for list of OOSAccesss.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type OOSAccessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []OOSAccess `json:"items"`
}

// VpnConnection defines the model for a VpnConnection.
// +genclient
// +genclient:nonNamespaced
// +groupName=oks.dev
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VpnConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   VpnConnectionSpec   `json:"spec"`
	Status VpnConnectionStatus `json:"status,omitempty"`
}

// VpnConnectionList defines model for list of VpnConnections.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VpnConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []VpnConnection `json:"items"`
}
