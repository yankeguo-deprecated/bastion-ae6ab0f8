package utils

import (
	"github.com/yankeguo/bastion/types"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

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

// LoadOptions load options from file
func LoadOptions(f string) (opt types.Options, err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(f); err != nil {
		return
	}
	if err = yaml.Unmarshal(buf, &opt); err != nil {
		return
	}
	// default values
	defaultStr(&opt.Daemon.DB, "/var/lib/bastion/database.sqlite3")
	defaultStr(&opt.Daemon.Host, "127.0.0.1")
	defaultInt(&opt.Daemon.Port, 9777)
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
	defaultStr(&opt.SSHD.SandboxEndpoint, "172.17.0.1")
	return
}
