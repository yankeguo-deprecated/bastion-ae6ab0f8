package sshd

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func loadSSHPrivateKeyFile(filename string) (s ssh.Signer, err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	s, err = ssh.ParsePrivateKey(buf)
	return
}
