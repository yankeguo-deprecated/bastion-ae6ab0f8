package web

import (
	"github.com/novakit/nova"
	"github.com/pkg/errors"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	contextKeyGRPCConn   = "_rpc_conn"
	contextKeyAuth       = "_auth"
	contextKeyWebOptions = "_web_options"

	headerKeyToken        = "X-Bastion-Token"
	headerKeyAction       = "X-Bastion-Action"
	headerValueClearToken = "clear-token"
)

func optsModule(opts types.WebOptions) nova.HandlerFunc {
	return func(c *nova.Context) error {
		c.Values[contextKeyWebOptions] = opts
		c.Next()
		return nil
	}
}

func webOptions(c *nova.Context) (opts types.WebOptions) {
	opts, _ = c.Values[contextKeyWebOptions].(types.WebOptions)
	return
}

func rpcModule(opts types.WebOptions) nova.HandlerFunc {
	return func(c *nova.Context) (err error) {
		var conn *grpc.ClientConn
		if conn, err = grpc.Dial(opts.DaemonEndpoint, grpc.WithInsecure()); err != nil {
			return
		}
		defer conn.Close()
		c.Values[contextKeyGRPCConn] = conn
		c.Next()
		return
	}
}

func userService(c *nova.Context) types.UserServiceClient {
	return types.NewUserServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func nodeService(c *nova.Context) types.NodeServiceClient {
	return types.NewNodeServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func grantService(c *nova.Context) types.GrantServiceClient {
	return types.NewGrantServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func keyService(c *nova.Context) types.KeyServiceClient {
	return types.NewKeyServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func sessionService(c *nova.Context) types.SessionServiceClient {
	return types.NewSessionServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func tokenService(c *nova.Context) types.TokenServiceClient {
	return types.NewTokenServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func replayService(c *nova.Context) types.ReplayServiceClient {
	return types.NewReplayServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func masterKeyService(c *nova.Context) types.MasterKeyServiceClient {
	return types.NewMasterKeyServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

// Auth result
type Auth struct {
	Token *types.Token
	User  *types.User
}

func (a Auth) IsLoggedIn() bool {
	return a.Token != nil && a.User != nil
}

func (a Auth) IsBlocked() bool {
	return a.IsLoggedIn() && a.User.IsBlocked
}

func (a Auth) IsLoggedInAsAdmin() bool {
	return a.IsLoggedIn() && a.User.IsAdmin
}

func markClearTokenIfNeeded(c *nova.Context, err error) {
	if err == nil {
		return
	}
	// mark the header X-Bastion-Action: clear-token if it's a invalid argument error
	if s, ok := status.FromError(err); ok {
		if s.Code() == codes.InvalidArgument {
			c.Res.Header().Set(headerKeyAction, headerValueClearToken)
		}
	}
}

func authModule() nova.HandlerFunc {
	return func(c *nova.Context) (err error) {
		ts, us := tokenService(c), userService(c)
		a := Auth{}
		token := c.Req.Header.Get(headerKeyToken)
		if len(token) != 0 {
			// get token by token value
			var res1 *types.GetTokenResponse
			if res1, err = ts.GetToken(c.Req.Context(), &types.GetTokenRequest{Token: token}); err != nil {
				markClearTokenIfNeeded(c, err)
				return
			}
			a.Token = res1.Token
			// get user
			var res2 *types.GetUserResponse
			if res2, err = us.GetUser(c.Req.Context(), &types.GetUserRequest{Account: a.Token.Account}); err != nil {
				markClearTokenIfNeeded(c, err)
				return
			}
			a.User = res2.User
			// touch token by token id, touch user by user account
			ts.TouchToken(c.Req.Context(), &types.TouchTokenRequest{Id: res1.Token.Id})
			us.TouchUser(c.Req.Context(), &types.TouchUserRequest{Account: res2.User.Account})
		}
		c.Values[contextKeyAuth] = a
		c.Next()
		return
	}
}

func requiresLoggedIn(admin bool) nova.HandlerFunc {
	return func(c *nova.Context) (err error) {
		a := authResult(c)
		if !a.IsLoggedIn() {
			err = errors.New("not logged in")
			return
		}
		if admin {
			if !a.IsLoggedInAsAdmin() {
				err = errors.New("not admin")
				return
			}
		}
		c.Next()
		return
	}
}

func authResult(c *nova.Context) (a Auth) {
	a, _ = c.Values[contextKeyAuth].(Auth)
	return
}
