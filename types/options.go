package types

// Options options for bunker
type Options struct {
	// Web
	// web options
	Web WebOptions `yaml:"web"`

	// Daemon
	// daemon options
	Daemon DaemonOptions `yaml:"daemon"`
}

// WebOptions web options for bunker
type WebOptions struct {
	// Host
	// host to bind for bunker web, default to "127.0.0.1"
	Host string `yaml:"host"`

	// Port
	// port to listen for bunker web, default to 9778
	Port int `yaml:"port"`

	// RPCEndpoint
	// address of bunker rpc service, default to "127.0.0.1:9777"
	RPCEndpoint string `yaml:"rpc_endpoint"`
}

// DaemonOptions daemon options for bunker
type DaemonOptions struct {
	// DB
	// database file path, using bolt, default to "/var/lib/bunker/database"
	DB string `yaml:"db"`

	// RPCHost
	// host to bind for bunker rpc, default to "127.0.0.1"
	RPCHost string `yaml:"rpc_host"`

	// RPCPort
	// port to bind for bunker rpc, default to 9777
	RPCPort int `yaml:"rpc_port"`

	// SSHDHost
	// host to bind for bunker sshd, default to "0.0.0.0"
	SSHDHost string `yaml:"sshd_host"`

	// SSHDPort
	// port to bind for bunker sshd, default to 22
	SSHDPort int `yaml:"sshd_port"`

	// ClientKey
	// client key file path for bunker ssh proxy, should be present on all target hosts' /root/.ssh/authorized_keys
	// default to "/etc/bunker/client_rsa"
	ClientKey string `yaml:"client_key"`

	// HostKey
	// host key file path for bunker sshd
	// default to "/etc/bunker/host_rsa"
	HostKey string `yaml:"host_key"`

	// SandboxImage
	// sandbox image, default to "bunker-sandbox"
	SandboxImage string `yaml:"sandbox_image"`

	// SandboxData
	// data dir for sandboxes, default to "/var/lib/bunker/sandboxes"
	SandboxData string `yaml:"sandbox_data"`

	// SandboxEndpoint
	// bunker IP address inside sandbox, basically the IP of docker0 virtual network adapter
	// default to "172.17.0.1"
	SandboxEndpoint string `yaml:"sandbox_endpoint"`
}
