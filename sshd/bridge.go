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
	"net"
	"sync"
)

func handleLv1DirectTCPIPChannel(conn *ssh.ServerConn, sc ssh.Channel, tp *TunnelPool, address string, port int) (err error) {
	ILog(conn).Str("channel", ChannelTypeDirectTCPIP).Msg("channel opened")
	defer ILog(conn).Str("channel", ChannelTypeDirectTCPIP).Err(err).Msg("channel finished")
	// remember to close channel
	defer sc.Close()
	// dial remote address
	var c net.Conn
	if c, err = tp.Dial(address, port); err != nil {
		ELog(conn).Str("channel", ChannelTypeDirectTCPIP).Err(err).Msg("failed to dial ssh tunnel connection")
		return
	}
	defer c.Close()
	// bi-copy streams
	if err = utils.DualCopy(c, sc); err != nil {
		ELog(conn).Str("channel", ChannelTypeDirectTCPIP).Err(err).Msg("failed to pipe bridge connection")
		return
	}
	return
}

func handleLv1SessionChannel(conn *ssh.ServerConn, sc ssh.Channel, srchan <-chan *ssh.Request, sb sandbox.Sandbox, account string, ss types.SessionServiceClient, rs types.ReplayServiceClient) (err error) {
	ILog(conn).Str("channel", ChannelTypeSession).Msg("channel opened")
	defer ILog(conn).Str("channel", ChannelTypeSession).Err(err).Msg("channel finished")
	// remember to close channel
	defer sc.Close()
	// variables
	cmd, cmdReady, cmdMissing, cmdCond := "", false, false, sync.NewCond(&sync.Mutex{})
	env := make([]string, 0)
	pty, ptyCols, ptyRows, term, wch := false, uint32(0), uint32(0), "", make(chan sandbox.Window, 4)
	// remember to close wch/rwch
	defer close(wch)
	// range all requests
	go func() {
		for req := range srchan {
			DLog(conn).Str("channel", ChannelTypeSession).Str("request", req.Type).Msg("request received from user")
			switch req.Type {
			case RequestTypePtyReq:
				var pl PtyRequestPayload
				if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
					ELog(conn).Str("channel", ChannelTypeSession).Str("request", req.Type).Err(err).Msg("failed to decode payload")
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
			case RequestTypeEnv:
				var pl EnvRequestPayload
				if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
					ELog(conn).Str("channel", ChannelTypeSession).Str("request", req.Type).Err(err).Msg("failed to decode payload")
					continue
				}
				env = append(env, fmt.Sprintf("%s=%s", pl.Name, pl.Value))
				if req.WantReply {
					req.Reply(true, nil)
				}
			case RequestTypeWindowChange:
				var pl WindowChangeRequestPayload
				if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
					ELog(conn).Str("channel", ChannelTypeSession).Str("request", req.Type).Err(err).Msg("failed to decode payload")
					continue
				}
				wch <- sandbox.Window{
					Width:  uint(pl.Cols),
					Height: uint(pl.Rows),
				}
				if req.WantReply {
					req.Reply(true, nil)
				}
			case RequestTypeExec, RequestTypeShell:
				// decode exec command
				if req.Type == RequestTypeExec {
					var pl ExecRequestPayload
					if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
						ELog(conn).Str("channel", ChannelTypeSession).Str("request", req.Type).Err(err).Msg("failed to decode payload")
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
		ELog(conn).Msg("command is missing")
		return
	}
	// split the command
	var cmds []string
	if cmds, err = shellquote.Split(cmd); err != nil {
		ELog(conn).Err(err).Str("command", cmd).Msg("failed to split command")
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
		ELog(conn).Err(err).Msg("failed to create session")
		return
	}
	// build the exec options
	opts := sandbox.ExecAttachOptions{
		Env:     env,
		Command: cmds,
		Stdin:   sc,
		Stdout:  sc,
		Stderr:  sc.Stderr(),
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
		ELog(conn).Err(err).Msg("exec attach returns error")
		es.Code = 1
	}
	// finish session
	ss.FinishSession(context.Background(), &types.FinishSessionRequest{Id: sRes.Session.Id})
	// send exit-status
	sc.SendRequest(RequestTypeExitStatus, false, ssh.Marshal(&es))
	return
}

func handleLv2SessionChannel(conn *ssh.ServerConn, sc ssh.Channel, srchan <-chan *ssh.Request, tc ssh.Channel, trchan <-chan *ssh.Request, user string) (err error) {
	ILog(conn).Str("channel", ChannelTypeSession).Msg("channel opened")
	defer ILog(conn).Str("channel", ChannelTypeSession).Err(err).Msg("channel finished")
	// remember to close channels
	defer sc.Close()
	defer tc.Close()
	// stream stdin, stdout, stderr, srchan <-> trchan
	wr := &sync.WaitGroup{}
	wr.Add(3)
	// srchan -> trchan
	go func() {
		for req := range srchan {
			DLog(conn).Str("channel", ChannelTypeSession).Str("request", req.Type).Msg("request received from remote server")
			// modify request with sudo
			switch req.Type {
			case RequestTypeExec:
				var pl ExecRequestPayload
				if err = ssh.Unmarshal(req.Payload, &pl); err != nil {
					ELog(conn).Str("channel", ChannelTypeSession).Str("request", req.Type).Err(err).Msg("failed to decode payload")
					continue
				}
				// switch user
				pl.Command = commandSwitchUser(user, pl.Command)
				// change payload
				req.Payload = ssh.Marshal(&pl)
			case RequestTypeShell:
				pl := ExecRequestPayload{
					Command: commandSwitchUser(user, ""),
				}
				// change request type to "exec" and update payload
				req.Type = RequestTypeExec
				req.Payload = ssh.Marshal(&pl)
			}
			// ban "x11-req", "subsystem" requests
			switch req.Type {
			case RequestTypeX11Req, RequestTypeSubsystem:
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
	}()
	// trchan -> srchan
	go func() {
		// transparent bridge requests
		for req := range trchan {
			DLog(conn).Str("channel", ChannelTypeSession).Str("request", req.Type).Msg("request received from user")
			ok, _ := sc.SendRequest(req.Type, req.WantReply, req.Payload)
			if req.WantReply {
				req.Reply(ok, nil)
			}
		}
	}()
	go utils.CopyWG(sc, tc, wr, &err)
	go utils.CopyWG(tc, sc, wr, &err)
	go utils.CopyWG(sc.Stderr(), tc.Stderr(), wr, &err)
	wr.Wait()
	return
}
