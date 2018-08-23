package sshd

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"sync"
)

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
