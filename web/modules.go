package web

import (
	"github.com/novakit/nova"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"context"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

const (
	contextKeyGRPCConn = "_rpc_conn"
	contextKeyAuth     = "_auth"

	headerKeyToken        = "X-Bastion-Token"
	headerKeyAction       = "X-Bastion-Action"
	headerValueClearToken = "clear-token"
)

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

// Auth result
type Auth struct {
	Token *types.Token
	User  *types.User
}

func (a Auth) IsLoggedIn() bool {
	return a.Token != nil && a.User != nil && !a.User.IsBlocked
}

func (a Auth) IsBlocked() bool {
	return a.User != nil && a.User.IsBlocked
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
			// get token
			var res *types.GetTokenResponse
			if res, err = ts.GetToken(context.Background(), &types.GetTokenRequest{Token: token}); err != nil {
				markClearTokenIfNeeded(c, err)
				return
			}
			a.Token = res.Token
			// get user
			var res2 *types.GetUserResponse
			if res2, err = us.GetUser(context.Background(), &types.GetUserRequest{Account: a.Token.Account}); err != nil {
				markClearTokenIfNeeded(c, err)
				return
			}
			a.User = res2.User
		}
		c.Values[contextKeyAuth] = a
		c.Next()
		return
	}
}

func authResult(c *nova.Context) (a Auth) {
	a, _ = c.Values[contextKeyAuth].(Auth)
	return
}
