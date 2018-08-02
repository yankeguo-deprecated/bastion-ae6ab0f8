package daemon

import (
	"context"
	"fmt"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"testing"
)

func TestDaemon_CreateFinishListSessions(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		ss := types.NewSessionServiceClient(conn)

		res, err := ss.CreateSession(context.Background(), &types.CreateSessionRequest{
			Account:    "test",
			IsRecorded: true,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res)

		res2, err := ss.FinishSession(context.Background(), &types.FinishSessionRequest{Id: res.Session.Id})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res2)
		if res2.Session.FinishedAt == 0 {
			t.Fatal("failed 1")
		}

		for i := 0; i < 1000; i++ {
			res, err = ss.CreateSession(context.Background(), &types.CreateSessionRequest{
				Account:    fmt.Sprintf("test%04d", i%10),
				IsRecorded: true,
			})
			if err != nil {
				t.Fatal(err)
			}
			res2, err = ss.FinishSession(context.Background(), &types.FinishSessionRequest{Id: res.Session.Id})
			if err != nil {
				t.Fatal(err)
			}
		}

		res3, err := ss.ListSessions(context.Background(), &types.ListSessionsRequest{
			Skip:  10,
			Limit: 4,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(res3.Sessions) != 4 {
			t.Fatal("failed 2")
		}
		if res3.Sessions[3].Account != "test0006" {
			t.Fatal("failed 3")
		}
		t.Log(res3)
	})
}
