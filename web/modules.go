package web

import (
	"github.com/novakit/nova"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
)

const (
	contextKeyGRPCConn = "_rpc_conn"
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

func rpcUserService(c *nova.Context) types.UserServiceClient {
	return types.NewUserServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func rpcNodeService(c *nova.Context) types.NodeServiceClient {
	return types.NewNodeServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func rpcGrantService(c *nova.Context) types.GrantServiceClient {
	return types.NewGrantServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func rpcKeyService(c *nova.Context) types.KeyServiceClient {
	return types.NewKeyServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}

func rpcSessionService(c *nova.Context) types.SessionServiceClient {
	return types.NewSessionServiceClient(c.Values[contextKeyGRPCConn].(*grpc.ClientConn))
}
