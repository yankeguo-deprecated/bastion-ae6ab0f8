package sshd

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"strings"
	"net"
	"os"
)

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

func isClosedError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err == os.ErrClosed
	}
	return err == os.ErrClosed
}

type DirectTCPIPExtraData struct {
	Host           string
	Port           uint32
	OriginatorIP   string
	OriginatorPort uint32
}

func decodeDirectTCPIPExtraData(p []byte) (pl DirectTCPIPExtraData, err error) {
	err = ssh.Unmarshal(p, &pl)
	return
}
