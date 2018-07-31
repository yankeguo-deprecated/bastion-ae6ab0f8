package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/yankeguo/bastion/daemon/db"
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/utils"
	"google.golang.org/grpc"
)

func TestDaemon_UpdateUser(t *testing.T) {
	d := New(types.DaemonOptions{
		DB:   temporaryFile(),
		Host: "127.0.0.1",
		Port: 10089,
	})
	go d.Run()
	defer d.Shutdown()
	time.Sleep(time.Second)
	d.DB.LogMode(true)

	// Set up a connection to the server.
	conn, err := grpc.Dial("127.0.0.1:10089", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := types.NewUserServiceClient(conn)

	nu := &db.User{Account: "test1", Nickname: "test user1", IsAdmin: 1}
	nu.PasswordDigest, _ = utils.BcryptGenerate("qwerty")
	d.DB.Create(&nu)

	u, err := c.UpdateUser(context.Background(), &types.UpdateUserRequest{
		Account:         "test1",
		UpdateIsAdmin:   true,
		IsAdmin:         false,
		UpdateIsBlocked: true,
		IsBlocked:       true,
		UpdatePassword:  true,
		Password:        "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if u.User.IsAdmin {
		t.Fatal("failed is_admin")
	}
	if !u.User.IsBlocked {
		t.Fatal("failed is_blocked")
	}
}
