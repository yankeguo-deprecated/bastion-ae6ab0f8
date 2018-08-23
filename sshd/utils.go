package sshd

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

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

func isClosedError(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err == os.ErrClosed
	}
	return err == os.ErrClosed
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

type DirectTCPIPExtraData struct {
	Host           string
	Port           uint32
	OriginatorIP   string
	OriginatorPort uint32
}

type PtyRequestPayload struct {
	Term   string
	Cols   uint32
	Rows   uint32
	Width  uint32
	Height uint32
	Modes  string
}

type EnvRequestPayload struct {
	Name  string
	Value string
}

type ExecRequestPayload struct {
	Command string
}

type WindowChangeRequestPayload struct {
	Cols   uint32
	Rows   uint32
	Width  uint32
	Height uint32
}

type ExitStatusRequestPayload struct {
	Code uint32
}
