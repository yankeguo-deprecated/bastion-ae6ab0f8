package daemon

import (
	"context"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestDaemon_PutListDeleteGrant(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		s := types.NewGrantServiceClient(conn)

		res, err := s.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local.*",
			User:            "root",
			ExpiredAt:       0,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res)
		res, err = s.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local.*",
			User:            "root",
			ExpiredAt:       10,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res)
		res2, err := s.ListGrants(context.Background(), &types.ListGrantsRequest{Account: "test"})
		if err != nil {
			t.Fatal(err)
		}
		if len(res2.Grants) != 1 {
			t.Fatal("bad count")
		}
		t.Log(res2)
		res, err = s.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local.*",
			User:            "work",
			ExpiredAt:       10,
		})
		if err != nil {
			t.Fatal(err)
		}
		res, err = s.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local2.*",
			User:            "work",
			ExpiredAt:       10,
		})
		if err != nil {
			t.Fatal(err)
		}
		res2, err = s.ListGrants(context.Background(), &types.ListGrantsRequest{Account: "test"})
		if err != nil {
			t.Fatal(err)
		}
		if len(res2.Grants) != 3 {
			t.Fatal("bad count")
		}
		t.Log(res2)
		_, err = s.DeleteGrant(context.Background(), &types.DeleteGrantRequest{Account: "test", HostnamePattern: "local.*", User: "work"})
		if err != nil {
			t.Fatal(err)
		}
		res2, err = s.ListGrants(context.Background(), &types.ListGrantsRequest{Account: "test"})
		if err != nil {
			t.Fatal(err)
		}
		if len(res2.Grants) != 2 {
			t.Fatal("bad count")
		}
		t.Log(res2)
	})
}

func TestDaemon_CheckGrantListGrantItems(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		rs := types.NewGrantServiceClient(conn)
		ns := types.NewNodeServiceClient(conn)

		// add nodes
		ns.PutNode(context.Background(), &types.PutNodeRequest{
			Hostname: "local.host1",
			Address:  "127.0.0.1:2222",
		})
		ns.PutNode(context.Background(), &types.PutNodeRequest{
			Hostname: "local.host2",
			Address:  "127.0.0.1:2223",
		})
		ns.PutNode(context.Background(), &types.PutNodeRequest{
			Hostname: "local2.host1",
			Address:  "127.0.0.1:2224",
		})
		ns.PutNode(context.Background(), &types.PutNodeRequest{
			Hostname: "local2.host2",
			Address:  "127.0.0.1:2225",
		})

		// add grants
		rs.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local.*",
			User:            "root",
		})
		rs.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local.*",
			User:            "work",
		})
		rs.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local2.*",
			User:            "work",
			ExpiredAt:       time.Now().Unix() + 1,
		})
		rs.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local2.host2",
			User:            "work",
			ExpiredAt:       time.Now().Unix() + 100,
		})
		rs.PutGrant(context.Background(), &types.PutGrantRequest{
			Account:         "test",
			HostnamePattern: "local2.host2",
			User:            "work",
			ExpiredAt:       time.Now().Unix() + 1000,
		})

		// check
		res, err := rs.CheckGrant(context.Background(), &types.CheckGrantRequest{
			Account:  "test",
			Hostname: "local2.host1",
			User:     "root",
		})
		if err != nil {
			t.Fatal(err)
		}
		if res.Ok {
			t.Fatal("failed 1")
		}
		res, err = rs.CheckGrant(context.Background(), &types.CheckGrantRequest{
			Account:  "test",
			Hostname: "local2.host1",
			User:     "work",
		})
		if err != nil {
			t.Fatal(err)
		}
		if !res.Ok {
			t.Fatal("failed 2")
		}
		// wait 2 second
		time.Sleep(time.Second * 2)
		res, err = rs.CheckGrant(context.Background(), &types.CheckGrantRequest{
			Account:  "test",
			Hostname: "local2.host1",
			User:     "work",
		})
		if err != nil {
			t.Fatal(err)
		}
		if res.Ok {
			t.Fatal("failed 3")
		}
		// list all
		res2, err := rs.ListGrantItems(context.Background(), &types.ListGrantItemsRequest{
			Account: "test",
		})
		if len(res2.GrantItems) != 5 {
			t.Fatal("failed 4")
		}
		if res2.GrantItems[len(res2.GrantItems)-1].ExpiredAt < time.Now().Unix()+500 {
			t.Fatal("failed 5")
		}
		for _, i := range res2.GrantItems {
			t.Log(i)
		}
	})
}
