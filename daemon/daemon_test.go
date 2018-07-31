package daemon

import (
	"path/filepath"
	"os"
	"encoding/hex"
	"crypto/rand"
	"testing"
	"github.com/yankeguo/bastion/types"
	"time"
	"google.golang.org/grpc"
	"context"
	"github.com/yankeguo/bastion/daemon/db"
	"github.com/yankeguo/bastion/utils"
)

func temporaryFile() string {
	buf := make([]byte, 8, 8)
	rand.Read(buf)
	return filepath.Join(os.TempDir(), "bnktestdb"+hex.EncodeToString(buf)+".sqlite3")
}

func TestDaemon_Run(t *testing.T) {
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

	_, err = c.AuthenticateUser(context.Background(), &types.AuthenticateUserRequest{
		Account:  "test1",
		Password: "qwerty",
	})
	if err == nil {
		t.Fatal("not failed")
	}

	nu := &db.User{Account: "test1", Nickname: "test user1"}
	nu.PasswordDigest, _ = utils.BcryptGenerate("qwerty")
	d.DB.Create(&nu)

	res, err := c.AuthenticateUser(context.Background(), &types.AuthenticateUserRequest{
		Account:  "test1",
		Password: "qwerty",
	})
	if err != nil {
		t.Fatal("failed", err)
	}
	if res.User.Nickname != "test user1" {
		t.Fatal("failed")
	}
	if time.Unix(res.User.CreatedAt, 0).In(time.UTC).Day() != time.Now().In(time.UTC).Day() {
		t.Fatal("failed created at")
	}
}
