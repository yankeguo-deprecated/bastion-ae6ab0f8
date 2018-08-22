/**
 * sandbox/sandbox.go
 * Copyright (c) 2018 Yanke Guo <guoyk.cn@gmail.com>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"syscall"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Window pty window size
type Window struct {
	Width  uint
	Height uint
}

// Pty pty information
type Pty struct {
	Term   string
	Window Window
}

// ExecAttachOptions opts for exec attach
type ExecAttachOptions struct {
	Env        []string
	Command    []string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	IsPty      bool
	Term       string
	WindowChan chan Window
}

// Sandbox interface
type Sandbox interface {
	GetContainerName() string
	Start() error
	GenerateSSHKey() error
	GetSSHPublicKey() (string, error)
	ExecScript(sc string) (string, string, error)
	ExecAttach(opts ExecAttachOptions) error
}

type sandbox struct {
	client *client.Client
	name   string
}

func (s *sandbox) GetContainerName() string {
	return s.name
}

func (s *sandbox) Start() error {
	return s.client.ContainerStart(context.Background(), s.name, dockerTypes.ContainerStartOptions{})
}

func (s *sandbox) GenerateSSHKey() (err error) {
	_, _, err = s.ExecScript(scriptGenerateSSHKey)
	return
}

func (s *sandbox) GetSSHPublicKey() (pkey string, err error) {
	pkey, _, err = s.ExecScript(`cat /root/.ssh/id_rsa.pub`)
	pkey = strings.TrimSpace(pkey)
	return
}

func (s *sandbox) ExecScript(sc string) (stdout string, stderr string, err error) {
	// create exec
	var id dockerTypes.IDResponse
	if id, err = s.client.ContainerExecCreate(
		context.Background(),
		s.name,
		dockerTypes.ExecConfig{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Cmd: []string{
				"/bin/bash",
			},
		},
	); err != nil {
		return
	}
	// exec attach
	var hr dockerTypes.HijackedResponse
	if hr, err = s.client.ContainerExecAttach(context.Background(), id.ID, dockerTypes.ExecConfig{}); err != nil {
		return
	}
	defer hr.Close()
	wg := &sync.WaitGroup{}
	// pipe stdin
	wg.Add(1)
	go func() {
		sandboxPipeStdin(hr, bytes.NewReader([]byte(sc)), &err)
		wg.Done()
	}()
	// pipe stdout/stderr
	bout := &bytes.Buffer{}
	berr := &bytes.Buffer{}
	wg.Add(1)
	go func() {
		sandboxPipeStdoutStderr(hr, bout, berr, false, &err)
		wg.Done()
	}()
	// wait
	wg.Wait()

	// output as string
	stdout = string(bout.Bytes())
	stderr = string(berr.Bytes())

	return
}

func (s *sandbox) ExecAttach(opts ExecAttachOptions) (err error) {
	execCfg := dockerTypes.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          opts.IsPty,
		Cmd:          opts.Command,
		Env:          opts.Env,
	}
	// append env if TERM is set
	if len(opts.Term) > 0 {
		if execCfg.Env == nil {
			execCfg.Env = make([]string, 0)
		}
		execCfg.Env = append(execCfg.Env, fmt.Sprintf("TERM=%s", opts.Term))
	}
	// use /bin/bash if no command specified
	if len(execCfg.Cmd) == 0 {
		execCfg.Cmd = []string{"/bin/bash"}
	}
	// create exec
	var id dockerTypes.IDResponse
	if id, err = s.client.ContainerExecCreate(context.Background(), s.name, execCfg); err != nil {
		return
	}
	// exec attach
	var hr dockerTypes.HijackedResponse
	if hr, err = s.client.ContainerExecAttach(context.Background(), id.ID, execCfg); err != nil {
		return
	}
	// pipe window size
	if opts.IsPty && opts.WindowChan != nil {
		go sandboxPipeWindowSize(s.client, id.ID, opts.WindowChan)
	}
	// pipe stdin
	go sandboxPipeStdin(hr, opts.Stdin, &err)
	// pipe stdout/stderr
	sandboxPipeStdoutStderr(hr, opts.Stdout, opts.Stderr, opts.IsPty, &err)
	// close hr
	hr.Close()
	// inspect exec
	var is dockerTypes.ContainerExecInspect
	if is, err = s.client.ContainerExecInspect(context.Background(), id.ID); err != nil {
		return
	}
	// send SIGTERM to zombie process
	if is.Running == true && is.Pid > 0 {
		log.Println("Exec:", id.ID, "not terminated properly, sending SIGTERM")
		var p *os.Process
		if p, err = os.FindProcess(is.Pid); err != nil {
			return
		}
		p.Signal(syscall.SIGTERM)
	}
	return
}

func sandboxPipeWindowSize(c *client.Client, id string, wchan chan Window) {
	for {
		w, ok := <-wchan
		if !ok {
			break
		}
		c.ContainerExecResize(context.Background(), id, dockerTypes.ResizeOptions{Height: w.Height, Width: w.Width})
	}
}

func sandboxPipeStdin(hr dockerTypes.HijackedResponse, stdin io.Reader, errout *error) {
	var err error
	_, err = io.Copy(hr.Conn, stdin)
	hr.CloseWrite()
	// clear EOF
	if err != nil && err != io.EOF {
		*errout = err
	}
	return
}

func sandboxPipeStdoutStderr(hr dockerTypes.HijackedResponse, stdout, stderr io.Writer, isPty bool, errout *error) {
	var err error
	if isPty {
		_, err = io.Copy(stdout, hr.Reader)
	} else {
		_, err = stdcopy.StdCopy(stdout, stderr, hr.Reader)
	}
	// clear EOF
	if err != nil && err != io.EOF {
		*errout = err
	}
	return
}
