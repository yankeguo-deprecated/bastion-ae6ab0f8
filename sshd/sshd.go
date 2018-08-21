package sshd

import (
	"github.com/yankeguo/bastion/types"
	"golang.org/x/crypto/ssh"
)

type SSHD struct {
	opts            types.SSHDOptions
	clientSigners   []ssh.Signer
	hostSigner      ssh.Signer
	sshServerConfig *ssh.ServerConfig
}

func New(opts types.SSHDOptions) *SSHD {
	return &SSHD{opts: opts}
}

func (s *SSHD) initHostSigner() (err error) {
	s.hostSigner, err = loadSSHPrivateKeyFile(s.opts.HostKey)
	return
}

func (s *SSHD) initClientSigners() (err error) {
	s.clientSigners = []ssh.Signer{}
	for _, key := range s.opts.ClientKeys {
		var cs ssh.Signer
		if cs, err = loadSSHPrivateKeyFile(key); err != nil {
			return
		}
		s.clientSigners = append(s.clientSigners, cs)
	}
	return
}

func (s *SSHD) Run() (err error) {
	// load host signer
	if err = s.initHostSigner(); err != nil {
		return
	}
	// load client signers
	if err = s.initClientSigners(); err != nil {
		return
	}
	return
}

func (s *SSHD) Shutdown() {
}
