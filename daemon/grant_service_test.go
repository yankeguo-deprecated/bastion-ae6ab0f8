package daemon

import (
	"testing"
	"google.golang.org/grpc"
	"github.com/yankeguo/bastion/types"
	"context"
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
