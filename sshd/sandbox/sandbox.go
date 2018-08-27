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
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

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
	var idRes dockerTypes.IDResponse
	if idRes, err = s.client.ContainerExecCreate(context.Background(), s.name, execCfg); err != nil {
		log.Error().Str("containerName", s.name).Err(err).Msg("failed to docker exec/create")
		return
	}
	// exec attach
	var hjRes dockerTypes.HijackedResponse
	if hjRes, err = s.client.ContainerExecAttach(context.Background(), idRes.ID, execCfg); err != nil {
		log.Error().Str("containerName", s.name).Err(err).Msg("failed to docker exec/attach")
		return
	}
	// pipe window size
	if opts.IsPty && opts.WindowChan != nil {
		go sandboxPipeWindowSize(s.client, idRes.ID, opts.WindowChan)
	}
	// pipe stdin
	go sandboxPipeStdin(hjRes, opts.Stdin, &err)
	// pipe stdout/stderr
	sandboxPipeStdoutStderr(hjRes, opts.Stdout, opts.Stderr, opts.IsPty, &err)
	// close hr
	hjRes.Close()
	// wait 1 second
	time.Sleep(time.Second)
	// SIGTERM
	s.signalExecIfNotExited(idRes.ID, syscall.SIGTERM)
	// wait to check process exit
	go func() {
		// wait 10 seconds
		time.Sleep(time.Second * 10)
		// SIGKILL
		s.signalExecIfNotExited(idRes.ID, syscall.SIGKILL)
	}()
	return
}

func (s *sandbox) signalExecIfNotExited(execId string, sig os.Signal) (err error) {
	var eiRes dockerTypes.ContainerExecInspect
	if eiRes, err = s.client.ContainerExecInspect(context.Background(), execId); err != nil {
		log.Error().Str("containerName", s.name).Str("execId", execId).Err(err).Msg("failed to inspect docker exec")
		return
	}
	// return if already stopped
	if !eiRes.Running || eiRes.Pid == 0 {
		return
	}
	// find the process
	var proc *os.Process
	if proc, err = os.FindProcess(eiRes.Pid); err != nil {
		log.Error().Str("containerName", s.name).Int("pid", eiRes.Pid).Str("execId", execId).Err(err).Msg("failed to find docker exec process")
		return
	}
	// send the signal
	if err = proc.Signal(sig); err != nil {
		log.Error().Str("containerName", s.name).Int("pid", eiRes.Pid).Str("execId", execId).Str("signal", sig.String()).Err(err).Msg("failed to send signal to docker exec")
		return
	}
	log.Info().Str("containerName", s.name).Int("pid", eiRes.Pid).Str("execId", execId).Str("signal", sig.String()).Msg("send signal to docker exec")
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
