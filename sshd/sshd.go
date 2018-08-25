package sshd

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/sshd/sandbox"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"net"
)

type SSHD struct {
	opts            types.SSHDOptions
	listener        net.Listener
	clientSigners   []ssh.Signer
	hostSigner      ssh.Signer
	sshServerConfig *ssh.ServerConfig

	rpcConn        *grpc.ClientConn
	sessionService types.SessionServiceClient
	replayService  types.ReplayServiceClient
	userService    types.UserServiceClient
	keyService     types.KeyServiceClient
	nodeService    types.NodeServiceClient
	grantService   types.GrantServiceClient

	sandboxManager sandbox.Manager
}

func New(opts types.SSHDOptions) *SSHD {
	return &SSHD{
		opts: opts,
	}
}

func (s *SSHD) initSandboxManager() (err error) {
	s.sandboxManager, err = sandbox.NewManager(s.opts)
	return
}

func (s *SSHD) initRPCConn() (err error) {
	if s.rpcConn, err = grpc.Dial(s.opts.DaemonEndpoint, grpc.WithInsecure()); err != nil {
		return
	}
	s.sessionService = types.NewSessionServiceClient(s.rpcConn)
	s.replayService = types.NewReplayServiceClient(s.rpcConn)
	s.userService = types.NewUserServiceClient(s.rpcConn)
	s.keyService = types.NewKeyServiceClient(s.rpcConn)
	s.nodeService = types.NewNodeServiceClient(s.rpcConn)
	s.grantService = types.NewGrantServiceClient(s.rpcConn)
	return
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

func (s *SSHD) initSSHServerConfig() (err error) {
	s.sshServerConfig = &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (ms *ssh.Permissions, err error) {
			ILog(conn).Msg("connection accepted")
			// decode target user and target node
			tu, th := decodeTargetServer(conn.User())
			// find the key
			var kRes *types.GetKeyResponse
			fp := ssh.FingerprintSHA256(key)
			if kRes, err = s.keyService.GetKey(context.Background(), &types.GetKeyRequest{Fingerprint: fp}); err != nil {
				ELog(conn).Str("fingerprint", fp).Err(err).Msg("failed to lookup key")
				err = errors.New("internal error: failed to lookup key")
				return
			}
			// touch the key
			if _, ierr := s.keyService.TouchKey(context.Background(), &types.TouchKeyRequest{Fingerprint: kRes.Key.Fingerprint}); ierr != nil {
				ELog(conn).Str("fingerprint", fp).Str("account", kRes.Key.Account).Err(ierr).Msg("failed to touch key")
			}
			// find the user
			var uRes *types.GetUserResponse
			if uRes, err = s.userService.GetUser(context.Background(), &types.GetUserRequest{Account: kRes.Key.Account}); err != nil {
				ELog(conn).Str("fingerprint", fp).Str("account", kRes.Key.Account).Err(err).Msg("failed to lookup user")
				err = errors.New("internal error: failed to lookup user")
				return
			}
			// check blocked user
			if uRes.User.IsBlocked {
				ILog(conn).Str("account", uRes.User.Account).Msg("trying to login a blocked user")
				err = errors.New("error: user is blocked")
				return
			}
			// touch the user
			if _, ierr := s.userService.TouchUser(context.Background(), &types.TouchUserRequest{Account: uRes.User.Account}); ierr != nil {
				ELog(conn).Str("account", uRes.User.Account).Err(ierr).Msg("failed to touch user")
			}
			// check internal connection
			if isSandboxConnection(conn, s.opts.SandboxEndpoint) {
				// check key source
				if kRes.Key.Source != types.KeySourceSandbox {
					ILog(conn).Str("fingerprint", fp).Str("account", uRes.User.Account).Msg("trying to enter lv2 stage with a non-sandbox key")
					err = errors.New("error: invalid key source")
					return
				}
				// check format
				if len(tu) == 0 || len(th) == 0 {
					ILog(conn).Str("account", uRes.User.Account).Str("sshUser", conn.User()).Msg("invalid lv2 stage ssh user format")
					err = errors.New("error: invalid format")
					return
				}
				// check node
				var nRes *types.GetNodeResponse
				if nRes, err = s.nodeService.GetNode(context.Background(), &types.GetNodeRequest{Hostname: th}); err != nil {
					ELog(conn).Str("account", uRes.User.Account).Err(err).Msg("failed to lookup node")
					err = errors.New("internal error: failed to lookup node")
					return
				}
				// check grant
				var cRes *types.CheckGrantResponse
				if cRes, err = s.grantService.CheckGrant(context.Background(), &types.CheckGrantRequest{
					User:     tu,
					Account:  uRes.User.Account,
					Hostname: nRes.Node.Hostname,
				}); err != nil {
					ELog(conn).Str("account", uRes.User.Account).Str("hostname", nRes.Node.Hostname).Str("user", tu).Err(err).Msg("failed to check grant")
					err = errors.New("internal error: failed to check permission")
					return
				}
				if !cRes.Ok {
					ILog(conn).Str("account", uRes.User.Account).Str("hostname", nRes.Node.Hostname).Str("user", tu).Msg("trying to access a not granted server")
					err = errors.New("error: no permission")
					return
				}
				ms = &ssh.Permissions{
					Extensions: map[string]string{
						extKeyAccount:  uRes.User.Account,
						extKeyUser:     tu,
						extKeyAddress:  nRes.Node.Address,
						extKeyHostname: nRes.Node.Hostname,
						extKeyStage:    stageLv2,
					},
				}
			} else {
				// connection from external
				// check recursive sandbox connection
				if kRes.Key.Source == types.KeySourceSandbox {
					ILog(conn).Str("fingerprint", fp).Str("account", uRes.User.Account).Msg("trying to enter lv1 stage with a sandbox key")
					err = errors.New("error: invalid key source")
					return
				}
				ms = &ssh.Permissions{
					Extensions: map[string]string{
						extKeyAccount: uRes.User.Account,
						extKeyStage:   stageLv1,
					},
				}
			}
			return
		},
	}
	s.sshServerConfig.AddHostKey(s.hostSigner)
	return
}

func (s *SSHD) initListener() (err error) {
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.opts.Host, s.opts.Port))
	return
}

func (s *SSHD) Run() (err error) {
	// init host signer
	if err = s.initHostSigner(); err != nil {
		return
	}
	// init client signers
	if err = s.initClientSigners(); err != nil {
		return
	}
	// init sandbox manager
	if err = s.initSandboxManager(); err != nil {
		return
	}
	// init rpcConn
	if err = s.initRPCConn(); err != nil {
		return
	}
	// init sshServerConfig, must after host signer and rpcConn
	if err = s.initSSHServerConfig(); err != nil {
		return
	}
	// init listener
	if err = s.initListener(); err != nil {
		return
	}
	for {
		var c net.Conn
		if c, err = s.listener.Accept(); err != nil {
			return
		}
		go s.handleConnection(c)
	}
	return
}

func (s *SSHD) handleConnection(c net.Conn) {
	var err error
	// variables
	var conn *ssh.ServerConn
	var nchan <-chan ssh.NewChannel
	var rchan <-chan *ssh.Request
	// upgrade connection to ssh connection
	if conn, nchan, rchan, err = ssh.NewServerConn(c, s.sshServerConfig); err != nil {
		log.Error().Err(err).Msg("failed to handshake")
		return
	}
	// handle the connection
	ILog(conn).Msg("connection established")
	if conn.Permissions.Extensions[extKeyStage] == stageLv1 {
		err = s.handleLv1Connection(conn, nchan, rchan)
	} else {
		err = s.handleLv2Connection(conn, nchan, rchan)
	}
	ILog(conn).Msg("connection finished")
	return
}

func (s *SSHD) updateSandboxPublicKey(sb sandbox.Sandbox, account string) (err error) {
	var ak string
	if ak, err = sb.GetSSHPublicKey(); err != nil {
		return
	}
	var pk ssh.PublicKey
	if pk, _, _, _, err = ssh.ParseAuthorizedKey([]byte(ak)); err != nil {
		return
	}
	fp := ssh.FingerprintSHA256(pk)
	_, err = s.keyService.CreateKey(context.Background(), &types.CreateKeyRequest{
		Name:        "sandbox",
		Account:     account,
		Fingerprint: fp,
		Source:      types.KeySourceSandbox,
	})
	return
}

func (s *SSHD) updateSandboxSSHConfig(sb sandbox.Sandbox, account string) (err error) {
	var riRes *types.ListGrantItemsResponse
	if riRes, err = s.grantService.ListGrantItems(context.Background(), &types.ListGrantItemsRequest{Account: account}); err != nil {
		return
	}
	se := make([]sandbox.SSHEntry, 0)
	for _, ri := range riRes.GrantItems {
		// skip the special tunnel user
		if ri.User == types.GrantUserTunnel {
			continue
		}
		se = append(se, sandbox.SSHEntry{
			Name: fmt.Sprintf("%s-%s", ri.Hostname, ri.User),
			Host: s.opts.SandboxEndpoint,
			Port: uint(s.opts.Port),
			User: fmt.Sprintf("%s@%s", ri.User, ri.Hostname),
		})
	}
	_, _, err = sb.ExecScript(sandbox.ScriptSeedSSHConfig(se))
	return
}

func (s *SSHD) handleLv1Connection(conn *ssh.ServerConn, ncchan <-chan ssh.NewChannel, grchan <-chan *ssh.Request) (err error) {
	// remember to close the connection
	defer conn.Close()
	account := conn.Permissions.Extensions[extKeyAccount]
	// discard global requests
	go discardRequests(grchan)
	// pre-create a connection-local tunnel pool for failure isolation
	tp := NewTunnelPool(s.clientSigners)
	defer tp.Close()
	// handle new channels
	for nc := range ncchan {
		if nc.ChannelType() == ChannelTypeDirectTCPIP {
			// if channel type is 'direct-tcpip'
			// extract host and port from extra data
			var pl DirectTCPIPExtraData
			if err = ssh.Unmarshal(nc.ExtraData(), &pl); err != nil {
				nc.Reject(ssh.UnknownChannelType, "internal error: invalid extra data for 'direct-tcpip' new channel request")
				ELog(conn).Str("channel", nc.ChannelType()).Hex("extraData", nc.ExtraData()).Err(err).Msg("invalid extra data for 'direct-tcpip'")
				continue
			}
			// find the node
			var nRes *types.GetNodeResponse
			if nRes, err = s.nodeService.GetNode(context.Background(), &types.GetNodeRequest{Hostname: pl.Host}); err != nil {
				nc.Reject(ssh.ConnectionFailed, "internal error: failed to lookup node")
				ELog(conn).Str("channel", nc.ChannelType()).Str("hostname", pl.Host).Err(err).Msg("failed to lookup node")
				continue
			}
			// check __tunnel__ user permission with given node
			var cRes *types.CheckGrantResponse
			if cRes, err = s.grantService.CheckGrant(context.Background(), &types.CheckGrantRequest{
				Account:  account,
				User:     types.GrantUserTunnel,
				Hostname: pl.Host,
			}); err != nil {
				nc.Reject(ssh.ConnectionFailed, "internal error: failed to check permission")
				ELog(conn).Str("channel", nc.ChannelType()).Str("hostname", pl.Host).Err(err).Msg("failed to lookup grant")
				continue
			}
			if !cRes.Ok {
				nc.Reject(ssh.ConnectionFailed, "error: no permission")
				ILog(conn).Str("channel", nc.ChannelType()).Str("hostname", pl.Host).Msg("trying to create tunnel on a not granted node")
				continue
			}
			// accept the new channel
			var sc ssh.Channel
			var srchan <-chan *ssh.Request
			if sc, srchan, err = nc.Accept(); err != nil {
				ELog(conn).Str("channel", nc.ChannelType()).Str("hostname", pl.Host).Err(err).Msg("failed to accept new channel")
				continue
			}
			// discard all channel-local requests
			go discardRequests(srchan)
			// dial and stream 'direct-tcpip'
			go handleLv1DirectTCPIPChannel(conn, sc, tp, nRes.Node.Address, int(pl.Port))
		} else if nc.ChannelType() == ChannelTypeSession {
			// find or create the sandbox
			var sb sandbox.Sandbox
			if sb, err = s.sandboxManager.FindOrCreate(account); err != nil {
				nc.Reject(ssh.ConnectionFailed, "internal error: failed to find or create the sandbox")
				ELog(conn).Str("channel", nc.ChannelType()).Err(err).Msg("failed to find or create the sandbox")
				continue
			}
			// load public key from sandbox /root/.ssh/id_rsa.pub
			if err = s.updateSandboxPublicKey(sb, account); err != nil {
				ELog(conn).Str("channel", nc.ChannelType()).Err(err).Msg("failed to extract sandbox public key")
			}
			// write sandbox /root/.ssh/config
			if err = s.updateSandboxSSHConfig(sb, account); err != nil {
				ELog(conn).Str("channel", nc.ChannelType()).Err(err).Msg("failed to write ssh config to sandbox")
			}
			// accept the new channel
			var sc ssh.Channel
			var srchan <-chan *ssh.Request
			if sc, srchan, err = nc.Accept(); err != nil {
				ELog(conn).Str("channel", nc.ChannelType()).Err(err).Msg("failed to accept new channel")
				continue
			}
			go handleLv1SessionChannel(conn, sc, srchan, sb, account, s.sessionService, s.replayService)
		} else {
			ELog(conn).Str("channel", nc.ChannelType()).Msg("unsupported channel type")
			nc.Reject(ssh.UnknownChannelType, "error: only channel type 'session' and 'direct-tcpip' is allowed")
			continue
		}
	}
	return
}

func (s *SSHD) handleLv2Connection(conn *ssh.ServerConn, ncchan <-chan ssh.NewChannel, grchan <-chan *ssh.Request) (err error) {
	defer conn.Close()
	// extract connection parameters
	user := conn.Permissions.Extensions[extKeyUser]
	address := conn.Permissions.Extensions[extKeyAddress]
	// no global requests is allowed in LV2 connection
	go discardRequests(grchan)
	// create ssh.Client
	var client *ssh.Client
	if client, err = ssh.Dial("tcp", fixSSHAddress(address), &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(s.clientSigners...)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}); err != nil {
		ELog(conn).Str("address", address).Msg("failed to create ssh client")
		return
	}
	defer client.Close()
	// iterate new channel requests
	for nc := range ncchan {
		// check channel type
		if nc.ChannelType() != ChannelTypeSession {
			nc.Reject(ssh.UnknownChannelType, "error: unsupported channel type")
			ILog(conn).Str("channel", nc.ChannelType()).Msg("unsupported channel type")
			continue
		}
		// open session on remote server
		var tc ssh.Channel
		var trchan <-chan *ssh.Request
		if tc, trchan, err = client.OpenChannel(ChannelTypeSession, nil); err != nil {
			nc.Reject(ssh.ConnectionFailed, "error: failed to create new session channel on remote server")
			ELog(conn).Str("channel", nc.ChannelType()).Err(err).Msg("failed to create new session channel on remote server")
			continue
		}
		// accept channel
		var sc ssh.Channel
		var srchan <-chan *ssh.Request
		if sc, srchan, err = nc.Accept(); err != nil {
			tc.Close()
			ELog(conn).Str("channel", nc.ChannelType()).Err(err).Msg("failed to accept new channel")
			continue
		}
		// bridge channels
		go handleLv2SessionChannel(conn, sc, srchan, tc, trchan, user)
	}
	return
}

func (s *SSHD) Shutdown() {
	if s.listener != nil {
		s.listener.Close()
	}
}
