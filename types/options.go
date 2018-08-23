package types

import (
	"encoding/json"
	"os"
	"path"
	"io/ioutil"
	"gopkg.in/yaml.v2"
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

	// ReplayDir directory of replay files
	ReplayDir string `yaml:"replay_dir"`
}

func (o DaemonOptions) String() string {
	buf, _ := json.MarshalIndent(&o, "", "  ")
	return "DaemonOptions" + string(buf)
}

// WebOptions web options for bastion
type WebOptions struct {
	// SSHDomain ssh target for display
	SSHDomain string `yaml:"ssh_domain"`

	// Host host to bind for bastion web, default to "127.0.0.1"
	Host string `yaml:"host"`

	// Port port to listen for bastion web, default to 9778
	Port int `yaml:"port"`

	// Dev development mode, will not use BinFS
	Dev bool `yaml:"dev"`

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
	Host string `yaml:"host"`

	// Port port to bind for bastion sshd, default to 22
	Port int `yaml:"port"`

	// DaemonEndpoint address of bastion daemon rpc service, default to "127.0.0.1:9777"
	DaemonEndpoint string `yaml:"daemon_endpoint"`

	// ClientKeys client key file path for bastion ssh proxy
	// should presents on all target hosts' /root/.ssh/authorized_keys
	// default to "/etc/bastion/client_rsa"
	ClientKeys []string `yaml:"client_keys"`

	// HostKey host key file path for bastion sshd
	// default to "/etc/bastion/host_rsa"
	HostKey string `yaml:"host_key"`

	// SandboxImage sandbox image, default to "bastion-sandbox"
	SandboxImage string `yaml:"sandbox_image"`

	// SandboxDir storage dir for sandboxes, default to "/var/lib/bastion/sandboxes"
	SandboxDir string `yaml:"sandbox_dir"`

	// SandboxEndpoint accessible bastion IP from sandbox, basically the IP of docker0 virtual network adapter
	// default to "172.17.0.1"
	SandboxEndpoint string `yaml:"sandbox_endpoint"`
}

func (o SSHDOptions) String() string {
	buf, _ := json.MarshalIndent(&o, "", "  ")
	return "SSHDOptions" + string(buf)
}

func defaultStr(s *string, d string) {
	if len(*s) == 0 {
		*s = d
	}
}

func defaultSts(s *[]string, d string) {
	if len(*s) == 0 {
		*s = []string{d}
	}
}

func defaultInt(i *int, d int) {
	if *i == 0 {
		*i = d
	}
}

func resolveDir(s *string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if !path.IsAbs(*s) {
		*s = path.Join(wd, *s)
	}
}

// LoadOptions load options from file
func LoadOptions(f string) (opt Options, err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(f); err != nil {
		return
	}
	if err = yaml.Unmarshal(buf, &opt); err != nil {
		return
	}
	// default values
	defaultStr(&opt.Daemon.DB, "/var/lib/bastion/database.bolt")
	defaultStr(&opt.Daemon.Host, "127.0.0.1")
	defaultInt(&opt.Daemon.Port, 9777)
	defaultStr(&opt.Daemon.ReplayDir, "/var/lib/bastion/replays")
	defaultStr(&opt.Web.Host, "127.0.0.1")
	defaultInt(&opt.Web.Port, 9778)
	defaultStr(&opt.Web.DaemonEndpoint, "127.0.0.1:9777")
	defaultStr(&opt.SSHD.Host, "0.0.0.0")
	defaultInt(&opt.SSHD.Port, 22)
	defaultStr(&opt.SSHD.DaemonEndpoint, "127.0.0.1:9777")
	defaultSts(&opt.SSHD.ClientKeys, "/etc/bastion/client_rsa")
	defaultStr(&opt.SSHD.HostKey, "/etc/bastion/host_rsa")
	defaultStr(&opt.SSHD.SandboxImage, "bastion-sandbox")
	defaultStr(&opt.SSHD.SandboxDir, "/var/lib/bastion/sandboxes")
	resolveDir(&opt.SSHD.SandboxDir)
	defaultStr(&opt.SSHD.SandboxEndpoint, "172.17.0.1")
	return
}
