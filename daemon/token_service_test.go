package daemon

import (
	"testing"
	"google.golang.org/grpc"
	"github.com/yankeguo/bastion/types"
	"context"
	"time"
)

func TestDaemon_CreateGetTouchListDeleteToken(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		ts := types.NewTokenServiceClient(conn)
		us := types.NewUserServiceClient(conn)

		us.CreateUser(context.Background(), &types.CreateUserRequest{
			Account:  "test",
			Password: "qwertyqwerty",
			Nickname: "testuser",
		})

		res, err := ts.CreateToken(context.Background(), &types.CreateTokenRequest{
			Account:     "test",
			Description: "some test token",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Token.Token) == 0 {
			t.Fatal("empty token generated")
		}
		t.Log(res)

		token := res.Token.Token

		res2, err := ts.GetToken(context.Background(), &types.GetTokenRequest{Token: token})
		if err != nil {
			t.Fatal(err)
		}
		if res2.Token.Account != "test" {
			t.Fatal("wrong")
		}
		t.Log(res2)

		va := res2.Token.ViewedAt

		time.Sleep(time.Second * 2)

		res3, err := ts.TouchToken(context.Background(), &types.TouchTokenRequest{Token: token})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res3)

		if res3.Token.ViewedAt == va {
			t.Fatal("failed 4")
		}

		res4, err := ts.ListTokens(context.Background(), &types.ListTokensRequest{Account: "test"})
		if err != nil {
			t.Fatal(err)
		}
		if len(res4.Tokens) != 1 {
			t.Fatal("failed 5")
		}
		if res4.Tokens[0].Token != token || res4.Tokens[0].Account != "test" {
			t.Fatal("failed 6")
		}

		t.Log(res4)

		_, err = ts.DeleteToken(context.Background(), &types.DeleteTokenRequest{Token: token})

		if err != nil {
			t.Fatal(err)
		}

		res4, err = ts.ListTokens(context.Background(), &types.ListTokensRequest{Account: "test"})
		if err != nil {
			t.Fatal(err)
		}
		if len(res4.Tokens) != 0 {
			t.Fatal("failed 5")
		}
	})
}
