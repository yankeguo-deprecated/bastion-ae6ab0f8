package daemon

import (
	"context"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"testing"
)

func TestDaemon_ListPutDeleteNode(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		ns := types.NewNodeServiceClient(conn)
		res, err := ns.ListNodes(context.Background(), &types.ListNodesRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Nodes) != 0 {
			t.Fatal("failed 1")
		}
		_, err = ns.PutNode(context.Background(), &types.PutNodeRequest{
			Hostname: "localhost1",
			Address:  "127.0.0.1",
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = ns.PutNode(context.Background(), &types.PutNodeRequest{
			Hostname: "localhost2",
			Address:  "127.0.0.2",
		})
		if err != nil {
			t.Fatal(err)
		}
		res, err = ns.ListNodes(context.Background(), &types.ListNodesRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Nodes) != 2 {
			t.Fatal("failed 2", res.Nodes)
		}
		if res.Nodes[1].Hostname != "localhost2" {
			t.Fatal("failed 3")
		}
		t.Log(res)
		_, err = ns.DeleteNode(context.Background(), &types.DeleteNodeRequest{Hostname: "localhost1"})
		if err != nil {
			t.Fatal(err)
		}
		res, err = ns.ListNodes(context.Background(), &types.ListNodesRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Nodes) != 1 {
			t.Fatal("failed 4", res.Nodes)
		}
		if res.Nodes[0].Hostname != "localhost2" {
			t.Fatal("failed 5")
		}
	})
}
