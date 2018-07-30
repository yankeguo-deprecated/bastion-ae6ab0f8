package utils

import (
	"github.com/yankeguo/bunker/types"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"bufio"
	"bytes"
	"log"
)

func defaultStr(s *string, d string) {
	if len(*s) == 0 {
		*s = d
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
	defaultStr(&opt.Web.Host, "127.0.0.1")
	defaultInt(&opt.Web.Port, 9778)
	defaultStr(&opt.Web.RPCEndpoint, "127.0.0.1:9777")
	defaultStr(&opt.Daemon.DB, "/var/lib/bunker/database.sqlite3")
	defaultStr(&opt.Daemon.RPCHost, "127.0.0.1")
	defaultInt(&opt.Daemon.RPCPort, 9777)
	defaultStr(&opt.Daemon.SSHDHost, "0.0.0.0")
	defaultInt(&opt.Daemon.SSHDPort, 22)
	defaultStr(&opt.Daemon.ClientKey, "/etc/bunker/client_rsa")
	defaultStr(&opt.Daemon.HostKey, "/etc/bunker/host_rsa")
	defaultStr(&opt.Daemon.SandboxImage, "bunker-sandbox")
	defaultStr(&opt.Daemon.SandboxData, "/var/lib/bunker/sandboxes")
	defaultStr(&opt.Daemon.SandboxEndpoint, "172.17.0.1")
	return
}

func PrintOptions(opt types.Options) {
	log.Println("CONFIG:")
	buf, _ := yaml.Marshal(&opt)
	r := bufio.NewReader(bytes.NewReader(buf))
	for {
		if l, err := r.ReadString('\n'); err != nil {
			return
		} else {
			log.Print("  " + l)
		}
	}
	log.Print("\n")
}
