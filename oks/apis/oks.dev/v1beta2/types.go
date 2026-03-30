package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodePool defines the model for a NodePool.
// +genclient
// +genclient:nonNamespaced
// +groupName=oks.dev
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NodePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   NodePoolSpec    `json:"spec"`
	Status *NodePoolStatus `json:"status,omitempty"`
}

// NodePoolList defines model for list of Nodepools.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NodePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []NodePool `json:"items"`
}
