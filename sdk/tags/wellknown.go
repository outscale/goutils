package tags

const (
	// EIPAutoAttach: Automatically associates a public IP with a VM
	EIPAutoAttach = "osc.fcu.eip.auto-attach"
	// RepulseServer: Places VMs with the same tag value on different servers
	RepulseServer = "osc.fcu.repulse_server"
	// AttractServer: Places VMs with the same tag value on the same server
	AttractServer = "osc.fcu.attract_server"
	// RepulseCluster: Places VMs with the same tag value on different Cisco UCS clusters
	RepulseCluster = "osc.fcu.repulse_cluster"
	// AttractCluster: Places VMs with the same tag value on the same Cisco UCS cluster
	AttractCluster = "osc.fcu.attract_cluster"
	// PrivateOnly: Blocks attribution of a public IP to a VM in the public Cloud (true or false)
	PrivateOnly = "private_only"
)
