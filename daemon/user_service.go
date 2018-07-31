package daemon

import (
	"golang.org/x/net/context"
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/daemon/db"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func (d *Daemon) ListUsers(context.Context, *types.ListUsersRequest) (*types.ListUsersResponse, error) {
	panic("implement me")
}

func (d *Daemon) CreateUser(context.Context, *types.CreateUserRequest) (*types.CreateUserResponse, error) {
	panic("implement me")
}

func (d *Daemon) TouchUser(context.Context, *types.TouchUserRequest) (*types.TouchUserResponse, error) {
	panic("implement me")
}

func (d *Daemon) UpdateUser(context.Context, *types.UpdateUserRequest) (*types.UpdateUserResponse, error) {
	panic("implement me")
}

func (d *Daemon) AuthenticateUser(c context.Context, req *types.AuthenticateUserRequest) (res *types.AuthenticateUserResponse, err error) {
	const notFound = "user not found or password invalid"
	u := db.User{}
	if err = d.DB.Find(&u, map[string]interface{}{"account": req.Account}).Error; err != nil {
		err = status.Error(codes.InvalidArgument, notFound)
		return
	}
	if u.ID == 0 || !u.CheckPassword(req.Password) {
		err = status.Error(codes.InvalidArgument, notFound)
		return
	}
	res = &types.AuthenticateUserResponse{User: u.ToRPCUser()}
	return
}
