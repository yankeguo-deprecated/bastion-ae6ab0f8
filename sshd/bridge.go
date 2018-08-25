package sshd

import (
	"context"
	"fmt"
	"github.com/kballard/go-shellquote"
	"github.com/yankeguo/bastion/sshd/recorder"
	"github.com/yankeguo/bastion/sshd/sandbox"
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/utils"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"sync"
)

func bridgeLv1DirectTCPIPChannel(sc ssh.Channel, tp *TunnelPool, address string, port int) (err error) {
	log.Println("lv1 direct-tcpip channel opened")
	defer log.Println("lv1 direct-tcpip channel finished, error =", err)
	defer sc.Close()
	// dial remote address
	var c net.Conn
	if c, err = tp.Dial(address, port); err != nil {
		log.Println("failed to dial ssh tunnel connection:", err)
		return
	}
	defer c.Close()
	// bi-copy streams
	if err = utils.DualCopy(c, sc); err != nil {
		log.Println("failed to pipe:", err)
		return
	}
	return
}

func bridgeLv2SessionChannel(sc ssh.Channel, srchan <-chan *ssh.Request, tc ssh.Channel, trchan <-chan *ssh.Request, user string) (err error) {
	log.Println("lv2 session channel opened")
	defer log.Println("lv2 session channel finished, error =", err)
	defer sc.Close()
	defer tc.Close()
	// stream stdin, stdout, stderr, srchan <-> trchan
	wr := &sync.WaitGroup{}
	wr.Add(5)
	// srchan -> trchan
	go func() {
		for req := range srchan {
			// modify request with sudo
			switch req.Type {
			case "exec":
				var pl ExecRequestPayload
				if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
					log.Println("failed to decode exec request payload", err)
					continue
				}
				// switch user
				pl.Command = commandSwitchUser(user, pl.Command)
				// change payload
				req.Payload = ssh.Marshal(&pl)
			case "shell":
				pl := ExecRequestPayload{
					Command: commandSwitchUser(user, ""),
				}
				// change request type to "exec" and update payload
				req.Type = "exec"
				req.Payload = ssh.Marshal(&pl)
			}
			// ban "x11-req", "subsystem" requests
			switch req.Type {
			case "x11-req", "subsystem":
				{
					if req.WantReply {
						req.Reply(false, nil)
					}
				}
			default:
				ok, _ := tc.SendRequest(req.Type, req.WantReply, req.Payload)
				if req.WantReply {
					req.Reply(ok, nil)
				}
			}
		}
		wr.Done()
	}()
	// trchan -> srchan
	go func() {
		// transparent bridge requests
		for req := range trchan {
			ok, _ := sc.SendRequest(req.Type, req.WantReply, req.Payload)
			if req.WantReply {
				req.Reply(ok, nil)
			}
		}
		wr.Done()
	}()
	go utils.CopyWG(sc, tc, wr, &err)
	go utils.CopyWG(tc, sc, wr, &err)
	go utils.CopyWG(sc.Stderr(), tc.Stderr(), wr, &err)
	wr.Wait()
	return
}

func bridgeLv1SessionChannel(chn ssh.Channel, crchan <-chan *ssh.Request, sb sandbox.Sandbox, account string, ss types.SessionServiceClient, rs types.ReplayServiceClient) (err error) {
	log.Println("lv1 session channel opened")
	defer log.Println("lv1 session channel finished, error =", err)
	defer chn.Close()
	// variables
	cmd, cmdReady, cmdMissing, cmdCond := "", false, false, sync.NewCond(&sync.Mutex{})
	env := make([]string, 0)
	pty, ptyCols, ptyRows, term, wch := false, uint32(0), uint32(0), "", make(chan sandbox.Window, 4)
	// remember to close wch/rwch
	defer close(wch)
	// range all requests
	go func() {
		for req := range crchan {
			switch req.Type {
			case "pty-req":
				var pl PtyRequestPayload
				if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
					log.Println("failed to decode pty-req payload:", err)
					continue
				}
				pty, ptyCols, ptyRows, term = true, pl.Cols, pl.Rows, pl.Term
				wch <- sandbox.Window{
					Width:  uint(pl.Cols),
					Height: uint(pl.Rows),
				}
				if req.WantReply {
					req.Reply(true, nil)
				}
			case "env":
				var pl EnvRequestPayload
				if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
					log.Println("failed to decode env request payload:", err)
					continue
				}
				env = append(env, fmt.Sprintf("%s=%s", pl.Name, pl.Value))
				if req.WantReply {
					req.Reply(true, nil)
				}
			case "window-change":
				var pl WindowChangeRequestPayload
				if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
					log.Println("failed to decode window-change request payload", err)
					continue
				}
				wch <- sandbox.Window{
					Width:  uint(pl.Cols),
					Height: uint(pl.Rows),
				}
				if req.WantReply {
					req.Reply(true, nil)
				}
			case "shell", "exec":
				// decode command
				if req.Type == "exec" {
					var pl ExecRequestPayload
					if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
						log.Println("failed to decode exec request payload", err)
						continue
					}
					cmd = pl.Command
				}
				if req.WantReply {
					req.Reply(true, nil)
				}
				// signal cmdCond
				cmdReady = true
				cmdCond.Signal()
			default:
				if req.WantReply {
					req.Reply(false, nil)
				}
			}
		}

		// if cmdReady not set, then cmd is missing, ensure cmdCond is always signaled
		if !cmdReady {
			cmdMissing = true
			cmdCond.Signal()
		}
	}()
	// wait for cmdCond
	cmdCond.L.Lock()
	for !cmdReady && !cmdMissing {
		cmdCond.Wait()
	}
	cmdCond.L.Unlock()
	// check if cmd is missing
	if cmdMissing {
		log.Println("command is missing")
		return
	}
	// split the command
	var cmds []string
	if cmds, err = shellquote.Split(cmd); err != nil {
		log.Println("failed to split command")
		return
	}
	// determine should be recorded
	isRecorded := shouldCommandBeRecorded(cmds)
	// start session
	var sRes *types.CreateSessionResponse
	if sRes, err = ss.CreateSession(context.Background(), &types.CreateSessionRequest{
		Account:    account,
		Command:    cmd,
		IsRecorded: isRecorded,
	}); err != nil {
		log.Println("failed to create session", err)
		return
	}
	// build the exec options
	opts := sandbox.ExecAttachOptions{
		Env:     env,
		Command: cmds,
		Stdin:   chn,
		Stdout:  chn,
		Stderr:  chn.Stderr(),
	}
	if pty {
		opts.IsPty = true
		opts.Term = term
		opts.WindowChan = wch
	}
	// wrap options if isRecorded
	if isRecorded {
		r := recorder.StartRecording(&opts, sRes.Session.Id, rs)
		defer r.Close()
	}
	// execute and returns exit status
	es := ExitStatusRequestPayload{}
	if err = sb.ExecAttach(opts); err != nil {
		log.Println("exec attach returns error:", err)
		es.Code = 1
	}
	// finish session
	ss.FinishSession(context.Background(), &types.FinishSessionRequest{Id: sRes.Session.Id})
	// send exit-status
	chn.SendRequest("exit-status", false, ssh.Marshal(&es))
	return
}
