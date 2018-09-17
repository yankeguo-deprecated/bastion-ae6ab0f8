package sshd

import (
	"bytes"
	"fmt"
	"github.com/kballard/go-shellquote"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strings"
	"sync"
)

func discardRequests(in <-chan *ssh.Request) {
	for req := range in {
		log.Debug().Str("requestType", req.Type).Msg("request discarded")
		if req.WantReply {
			req.Reply(false, nil)
		}
	}
}

func sshClientOverrideKeys(client *ssh.Client, keys []ssh.Signer) (err error) {
	if len(keys) == 0 {
		return
	}
	var aks []byte
	for _, key := range keys {
		buf := ssh.MarshalAuthorizedKey(key.PublicKey())
		buf = bytes.TrimSpace(buf)
		aks = append(aks, buf...)
		aks = append(aks, '\n')
	}
	var session *ssh.Session
	if session, err = client.NewSession(); err != nil {
		return
	}
	defer session.Close()
	session.Stdin = bytes.NewReader([]byte(aks))
	if err = session.Run("cat > /root/.ssh/authorized_keys"); err != nil {
		log.Error().Err(err).Msg("failed to execute command")
		return
	}
	return
}

func fixSSHAddress(address string) string {
	if len(strings.Split(address, ":")) < 2 {
		return address + ":22"
	}
	return address
}

func loadSSHPrivateKeyFile(filename string) (s ssh.Signer, err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	s, err = ssh.ParsePrivateKey(buf)
	return
}

func decodeTargetServer(input string) (user string, host string) {
	ds := strings.Split(input, "@")
	if len(ds) == 2 {
		user = ds[0]
		host = ds[1]
	}
	return
}

func isSandboxConnection(conn ssh.ConnMetadata, endpoint string) bool {
	hostIP := net.ParseIP(endpoint)
	if addr, ok := conn.LocalAddr().(*net.TCPAddr); ok {
		return addr.IP.Equal(hostIP)
	}
	return false
}

func shouldCommandBeRecorded(cmd []string) bool {
	if len(cmd) == 0 {
		return true
	}
	if strings.ToLower(strings.TrimSpace(cmd[0])) == "scp" {
		return false
	}
	return true
}

func commandSwitchUser(user string, input string) string {
	if len(input) > 0 {
		return shellquote.Join("sudo", "-S", "-n", "-u", user, "-i", "--", "bash", "-c", input)
	}
	return shellquote.Join("sudo", "-S", "-n", "-u", user, "-i")
}

type TunnelPool struct {
	clientSigners []ssh.Signer
	clients       map[string]*ssh.Client
	clientsMutex  *sync.Mutex
}

func NewTunnelPool(clientSigners []ssh.Signer) *TunnelPool {
	return &TunnelPool{
		clientSigners: clientSigners,
		clients:       map[string]*ssh.Client{},
		clientsMutex:  &sync.Mutex{},
	}
}

func (t *TunnelPool) GetClient(address string) (c *ssh.Client, err error) {
	t.clientsMutex.Lock()
	defer t.clientsMutex.Unlock()
	address = fixSSHAddress(address)
	c = t.clients[address]
	if c == nil {
		if c, err = ssh.Dial("tcp", address, &ssh.ClientConfig{
			User:            "root",
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(t.clientSigners...)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}); err != nil {
			return
		}
		t.clients[address] = c
	}
	return
}

func (t *TunnelPool) Dial(address string, port int) (c net.Conn, err error) {
	var cl *ssh.Client
	if cl, err = t.GetClient(address); err != nil {
		return
	}
	return cl.Dial("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", port))
}

func (t *TunnelPool) Close() {
	t.clientsMutex.Lock()
	defer t.clientsMutex.Unlock()
	for _, c := range t.clients {
		c.Close()
	}
	t.clients = map[string]*ssh.Client{}
}
