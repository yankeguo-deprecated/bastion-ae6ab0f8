package daemon

import (
	"github.com/yankeguo/bastion/daemon/db"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/yankeguo/bastion/utils"
)

var (
	errInvalidAuthentication = status.Error(codes.InvalidArgument, "user not found or invalid password")
)

func whereAccount(a string) map[string]interface{} {
	return map[string]interface{}{"account": a}
}

func (d *Daemon) ListUsers(context.Context, *types.ListUsersRequest) (*types.ListUsersResponse, error) {
	panic("implement me")
}

func (d *Daemon) CreateUser(context.Context, *types.CreateUserRequest) (*types.CreateUserResponse, error) {
	panic("implement me")
}

func (d *Daemon) TouchUser(context.Context, *types.TouchUserRequest) (*types.TouchUserResponse, error) {
	panic("implement me")
}

func (d *Daemon) UpdateUser(c context.Context, req *types.UpdateUserRequest) (res *types.UpdateUserResponse, err error) {
	update := map[string]interface{}{}
	if req.UpdateNickname {
		update["nickname"] = req.Nickname
	}
	if req.UpdatePassword {
		var hs string
		if hs, err = utils.BcryptGenerate(req.Password); err != nil {
			err = errInternal
			return
		}
		update["password_digest"] = hs
	}
	if req.UpdateIsAdmin {
		if req.IsAdmin {
			update["is_admin"] = 1
		} else {
			update["is_admin"] = 0
		}
	}
	if req.UpdateIsBlocked {
		if req.IsBlocked {
			update["is_blocked"] = 1
		} else {
			update["is_blocked"] = 0
		}
	}
	q := d.DB.Model(new(db.User)).Where(whereAccount(req.Account)).Update(update)
	if err = q.Error; err != nil {
		err = errDatabase
		return
	}
	if q.RowsAffected == 0 {
		err = errRecordNotFound
		return
	}
	u := db.User{}
	if err = d.DB.Find(&u, whereAccount(req.Account)).Error; err != nil {
		err = DatabaseErrorToGRPCError(err)
		return
	}
	res = &types.UpdateUserResponse{User: u.ToRPCUser()}
	return
}

func (d *Daemon) AuthenticateUser(c context.Context, req *types.AuthenticateUserRequest) (res *types.AuthenticateUserResponse, err error) {
	u := db.User{}
	if err = d.DB.Find(&u, whereAccount(req.Account)).Error; err != nil {
		if IsRecordNotFound(err) {
			// not using record not found, hide the existence of the user
			err = errInvalidAuthentication
		} else {
			err = errDatabase
		}
		return
	}
	if u.ID == 0 || !utils.BcryptValidate(u.PasswordDigest, req.Password) {
		err = errInvalidAuthentication
		return
	}
	res = &types.AuthenticateUserResponse{User: u.ToRPCUser()}
	return
}
