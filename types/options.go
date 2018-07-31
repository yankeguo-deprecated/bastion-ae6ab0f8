package types

import (
		"encoding/json"
)

// Options options for bastion
type Options struct {
	// Daemon daemon options
	Daemon DaemonOptions `yaml:"daemon"`

	// Web web options
	Web WebOptions `yaml:"web"`

	// SSHD sshd options
	SSHD SSHDOptions `yaml:"sshd"`
}

func (o Options) String() string {
	buf, _ := json.MarshalIndent(&o, "", "  ")
	return "Options" + string(buf)
}

// DaemonOptions daemon options for bastion
type DaemonOptions struct {
	// DB database file path, using bolt, default to "/var/lib/bastion/database"
	DB string `yaml:"db"`

	// Host host to bind for bastion rpc, default to "127.0.0.1"
	Host string `yaml:"host"`

	// Port port to bind for bastion rpc, default to 9777
	Port int `yaml:"port"`

	// Consul whether using consul nodes catalog
	Consul bool `yaml:"consul"`
}

func (o DaemonOptions) String() string {
	buf, _ := json.MarshalIndent(&o, "", "  ")
	return "DaemonOptions" + string(buf)
}

// WebOptions web options for bastion
type WebOptions struct {
	// Host host to bind for bastion web, default to "127.0.0.1"
	Host string `yaml:"host"`

	// Port port to listen for bastion web, default to 9778
	Port int `yaml:"port"`

	// DaemonEndpoint address of bastion daemon rpc service, default to "127.0.0.1:9777"
	DaemonEndpoint string `yaml:"daemon_endpoint"`
}

func (o WebOptions) String() string {
	buf, _ := json.MarshalIndent(&o, "", "  ")
	return "WebOptions" + string(buf)
}

// SSHDOptions sshd options
type SSHDOptions struct {
	// Host host to bind for bastion sshd, default to "0.0.0.0"
	Host string `yaml:"sshd_host"`

	// Port port to bind for bastion sshd, default to 22
	Port int `yaml:"sshd_port"`

	// DaemonEndpoint address of bastion daemon rpc service, default to "127.0.0.1:9777"
	DaemonEndpoint string `yaml:"daemon_endpoint"`

	// ClientKeys client key file path for bastion ssh proxy
	// should presents on all target hosts' /root/.ssh/authorized_keys
	// default to "/etc/bastion/client_rsa"
	ClientKeys []string `yaml:"client_key"`

	// HostKey host key file path for bastion sshd
	// default to "/etc/bastion/host_rsa"
	HostKey string `yaml:"host_key"`

	// SandboxImage sandbox image, default to "bastion-sandbox"
	SandboxImage string `yaml:"sandbox_image"`

	// SandboxDir storage dir for sandboxes, default to "/var/lib/bastion/sandboxes"
	SandboxDir string `yaml:"sandbox_dir"`

	// ReplayDir directory of replay files
	ReplayDir string `yaml:"replay_dir"`

	// SandboxEndpoint accessible bastion IP from sandbox, basically the IP of docker0 virtual network adapter
	// default to "172.17.0.1"
	SandboxEndpoint string `yaml:"sandbox_endpoint"`
}

func (o SSHDOptions) String() string {
	buf, _ := json.MarshalIndent(&o, "", "  ")
	return "SSHDOptions" + string(buf)
}
