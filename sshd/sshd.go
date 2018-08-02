package sshd

import "github.com/yankeguo/bastion/types"

type SSHD struct {
	opts types.SSHDOptions
}

func New(opts types.SSHDOptions) *SSHD {
	return &SSHD{opts: opts}
}

func (s *SSHD) Run() (err error) {
	return
}

func (s *SSHD) Shutdown() {
}
