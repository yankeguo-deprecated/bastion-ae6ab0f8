package daemon

import (
	"context"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestDaemon_ListUsers(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		c := types.NewUserServiceClient(conn)
		c.CreateUser(context.Background(), &types.CreateUserRequest{Account: "test1", Password: "qwerty"})
		c.CreateUser(context.Background(), &types.CreateUserRequest{Account: "test2", Password: "qwerty"})
		c.CreateUser(context.Background(), &types.CreateUserRequest{Account: "test3", Password: "qwerty"})

		res, err := c.ListUsers(context.Background(), &types.ListUsersRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Users) != 3 {
			t.Fatal("failed 1")
		}
		if res.Users[2].Nickname != "test3" {
			t.Fatal("failed 2")
		}
		t.Log(res)
	})
}

func TestDaemon_CreateGetUser(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		c := types.NewUserServiceClient(conn)
		res, err := c.CreateUser(context.Background(), &types.CreateUserRequest{
			Account:  "testuser",
			Password: "qwerty",
		})
		if err != nil {
			t.Fatal(err)
		}
		if res.User.Nickname != "testuser" {
			t.Fatal("failed 1")
		}
		t.Log(res)
		res, err = c.CreateUser(context.Background(), &types.CreateUserRequest{
			Account:  "testuser",
			Password: "qwerty",
		})
		if err == nil {
			t.Fatal("failed 2")
		}
		t.Log(res, err)
		res, err = c.CreateUser(context.Background(), &types.CreateUserRequest{
			Account:  "testuser2",
			Password: "",
		})
		if err == nil {
			t.Fatal("failed 3")
		}
		t.Log(res, err)
		res1, err := c.GetUser(context.Background(), &types.GetUserRequest{
			Account: "testuser",
		})
		if err != nil {
			t.Fatal(err)
		}
		if res1.User.CreatedAt == 0 {
			t.Fatal("failed 4")
		}
		t.Log(res1)
	})
}

func TestDaemon_TouchUser(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		c := types.NewUserServiceClient(conn)
		u := models.User{
			Account:   "testuser",
			IsAdmin:   false,
			IsBlocked: true,
		}
		u.PasswordDigest, _ = bcryptGenerate("qwerty")
		daemon.db.Save(&u)
		res, err := c.TouchUser(context.Background(), &types.TouchUserRequest{
			Account: "testuser",
		})
		if err != nil || res.User.ViewedAt == 0 {
			t.Fatal(err)
		}
		t.Log(res)
		ov := res.User.ViewedAt
		time.Sleep(time.Second * 2)
		res, err = c.TouchUser(context.Background(), &types.TouchUserRequest{
			Account: "testuser",
		})
		if ov == res.User.ViewedAt {
			t.Fatal("failed 4")
		}
		t.Log(res)
	})
}

func TestDaemon_UpdateUser(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		c := types.NewUserServiceClient(conn)
		c.CreateUser(context.Background(), &types.CreateUserRequest{
			Account:  "testuser",
			Password: "qwerty",
			IsAdmin:  true,
		})
		time.Sleep(time.Second + time.Millisecond*100)
		res, err := c.UpdateUser(context.Background(), &types.UpdateUserRequest{
			Account:         "testuser",
			UpdateIsAdmin:   true,
			IsAdmin:         false,
			UpdateIsBlocked: true,
			IsBlocked:       true,
			UpdateNickname:  true,
			Nickname:        "test user ",
		})
		if err != nil || res.User.UpdatedAt == res.User.CreatedAt {
			t.Fatal(err)
		}
		if res.User.IsAdmin {
			t.Fatal("failed 1")
		}
		if !res.User.IsBlocked {
			t.Fatal("failed 2")
		}
		if res.User.Nickname != "test user" {
			t.Fatal("failed 3")
		}
		t.Log(res)
	})
}

func TestDaemon_AuthenticateUser(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		c := types.NewUserServiceClient(conn)
		_, err := c.AuthenticateUser(context.Background(), &types.AuthenticateUserRequest{
			Account:  "testuser",
			Password: "qwerty",
		})
		if err == nil {
			t.Fatal("failed 1")
		}
		u := models.User{
			Account: "testuser",
		}
		u.PasswordDigest, _ = bcryptGenerate("qwerty")
		daemon.db.Save(&u)
		_, err = c.AuthenticateUser(context.Background(), &types.AuthenticateUserRequest{
			Account:  "testuser",
			Password: "abcdef",
		})
		if err == nil {
			t.Fatal("failed 2")
		}
		res, err := c.AuthenticateUser(context.Background(), &types.AuthenticateUserRequest{
			Account:  "testuser",
			Password: "qwerty",
		})
		if err != nil {
			t.Fatal("failed 3")
		}
		if res.User == nil {
			t.Fatal("failed 4")
		}
		t.Log(res)
	})
}
