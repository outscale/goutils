package v1beta

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	schemeBuilder = &runtime.SchemeBuilder{}
	// AddToScheme adds to the scheme.
	AddToScheme = schemeBuilder.AddToScheme

	// SchemeGroupVersion is group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{Group: "oks.dev", Version: "v1beta"}
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	schemeBuilder.Register(addKnownTypes)
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&IPPool{},
		&IPPoolList{},
		&NetPeering{},
		&NetPeeringList{},
		&NetPeeringRequest{},
		&NetPeeringRequestList{},
		&NetPeeringAcceptance{},
		&NetPeeringAcceptanceList{},
		&OOSAccess{},
		&OOSAccessList{},
		&VpnConnection{},
		&VpnConnectionList{},
	)

	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&metav1.Status{},
	)

	metav1.AddToGroupVersion(
		scheme,
		SchemeGroupVersion,
	)

	return nil
}
