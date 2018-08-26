package daemon

import (
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errInvalidPassword = status.Error(codes.InvalidArgument, "invalid password")
	errUserBlocked     = status.Error(codes.InvalidArgument, "user blocked")
)

func (d *Daemon) ListUsers(c context.Context, req *types.ListUsersRequest) (res *types.ListUsersResponse, err error) {
	var users []models.User
	if err = d.db.All(&users); err != nil {
		return
	}
	ret := make([]*types.User, 0, len(users))
	for _, u := range users {
		ret = append(ret, u.ToGRPCUser())
	}
	res = &types.ListUsersResponse{Users: ret}
	return
}

func (d *Daemon) CreateUser(c context.Context, req *types.CreateUserRequest) (res *types.CreateUserResponse, err error) {
	// fix request
	if err = req.Validate(); err != nil {
		return
	}

	// inside a transaction
	u := models.User{}
	err = d.db.Tx(true, func(db *Node) (err error) {
		// find existing
		if err = db.CheckDuplicated("User", "account", req.Account); err != nil {
			return
		}
		// assign values
		copier.Copy(&u, req)
		// create password
		if u.PasswordDigest, err = bcryptGenerate(req.Password); err != nil {
			err = errInternal
			return
		}
		// assign created_at / updated_at and save
		u.CreatedAt = now()
		u.UpdatedAt = u.CreatedAt
		if err = db.Save(&u); err != nil {
			return
		}
		return
	})
	// return if err != nil
	if err != nil {
		return
	}
	// build response
	res = &types.CreateUserResponse{User: u.ToGRPCUser()}
	return
}

func (d *Daemon) TouchUser(c context.Context, req *types.TouchUserRequest) (res *types.TouchUserResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	u := models.User{}
	if err = d.db.Tx(true, func(db *Node) (err error) {
		// find by account
		if err = db.One("Account", req.Account, &u); err != nil {
			return
		}
		// update viewed_at
		u.ViewedAt = now()
		// save
		if err = db.Save(&u); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}
	// build response
	res = &types.TouchUserResponse{User: u.ToGRPCUser()}
	return
}

func (d *Daemon) UpdateUser(c context.Context, req *types.UpdateUserRequest) (res *types.UpdateUserResponse, err error) {
	// validate request
	if err = req.Validate(); err != nil {
		return
	}
	// find user by account
	u := models.User{}
	if err = d.db.One("Account", req.Account, &u); err != nil {
		return
	}
	// update user
	if req.UpdateIsBlocked {
		u.IsBlocked = req.IsBlocked
		// clear password failed when unblocking a user
		if !req.IsBlocked {
			u.PasswordFailed = 0
		}
	}
	if req.UpdateIsAdmin {
		u.IsAdmin = req.IsAdmin
	}
	if req.UpdateNickname {
		u.Nickname = req.Nickname
	}
	if req.UpdatePassword {
		if u.PasswordDigest, err = bcryptGenerate(req.Password); err != nil {
			err = errInternal
			return
		}
	}
	// update updated_at
	u.UpdatedAt = now()
	// save
	if err = d.db.Save(&u); err != nil {
		return
	}
	// build response
	res = &types.UpdateUserResponse{User: u.ToGRPCUser()}
	return
}

func (d *Daemon) AuthenticateUser(c context.Context, req *types.AuthenticateUserRequest) (res *types.AuthenticateUserResponse, err error) {
	u := models.User{}
	// find by account
	if err = d.db.One("Account", req.Account, &u); err != nil {
		return
	}
	// check blocked
	if u.IsBlocked {
		err = errUserBlocked
		return
	}
	// validate password
	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(req.Password)); err != nil {
		err = errInvalidPassword
		// update PasswordFailed, if failed too many times, block user
		u.PasswordFailed = u.PasswordFailed + 1
		log.Debug().Str("account", u.Account).Int64("failed", u.PasswordFailed).Msg("failed increased")
		if u.PasswordFailed > 6 {
			log.Debug().Str("account", u.Account).Int64("failed", u.PasswordFailed).Msg("blocked due to failed too much")
			u.IsBlocked = true
		}
		d.db.Save(&u)
		return
	}
	// clear PasswordFailed
	if u.PasswordFailed > 0 {
		u.PasswordFailed = 0
		log.Debug().Str("account", u.Account).Int64("failed", u.PasswordFailed).Msg("failed cleared")
		d.db.Save(&u)
	}
	// build response
	res = &types.AuthenticateUserResponse{User: u.ToGRPCUser()}
	return
}

func (d *Daemon) GetUser(c context.Context, req *types.GetUserRequest) (res *types.GetUserResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	u := models.User{}
	if err = d.db.One("Account", req.Account, &u); err != nil {
		return
	}
	res = &types.GetUserResponse{User: u.ToGRPCUser()}
	return
}
