package role

// Role defines a role, used to define roles of subnets/SGs.
type Role string

const (
	// ControlPlane is used for control-plane nodes.
	ControlPlane Role = "controlplane"
	// Worker is used for worker nodes.
	Worker Role = "worker"
	// LoadBalancer is used for the kube api LB.
	LoadBalancer Role = "loadbalancer"
	// Bastion is used for the bastion VM.
	Bastion Role = "bastion"
	// Nat is used for NAT services.
	Nat Role = "nat"
	// Service is used for public service LBs.
	Service Role = "service"
	// InternalService is used for internal service LBs.
	InternalService Role = "service.internal"
)
