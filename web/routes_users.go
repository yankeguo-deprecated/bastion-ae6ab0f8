package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"strings"
	"github.com/novakit/router"
)

func routeListUsers(c *nova.Context) (err error) {
	us, v := userService(c), view.Extract(c)
	var res1 *types.ListUsersResponse
	if res1, err = us.ListUsers(c.Req.Context(), &types.ListUsersRequest{}); err != nil {
		return
	}
	v.Data["users"] = res1.Users
	v.DataAsJSON()
	return
}

func routeCreateUser(c *nova.Context) (err error) {
	us, v := userService(c), view.Extract(c)
	var res1 *types.CreateUserResponse
	if res1, err = us.CreateUser(c.Req.Context(), &types.CreateUserRequest{
		Account:  c.Req.FormValue("account"),
		Nickname: c.Req.FormValue("nickname"),
		Password: c.Req.FormValue("password"),
	}); err != nil {
		return
	}
	v.Data["user"] = res1.User
	v.DataAsJSON()
	return
}

func routeGetCurrentUser(c *nova.Context) (err error) {
	a, v := authResult(c), view.Extract(c)
	v.Data["user"] = a.User
	v.DataAsJSON()
	return
}

func routeGetCurrentUserGrantItems(c *nova.Context) (err error) {
	a, gs, v := authResult(c), grantService(c), view.Extract(c)
	var res1 *types.ListGrantItemsResponse
	if res1, err = gs.ListGrantItems(c.Req.Context(), &types.ListGrantItemsRequest{
		Account: a.User.Account,
	}); err != nil {
		return
	}
	v.Data["grant_items"] = res1.GrantItems
	v.DataAsJSON()
	return
}

func routeUpdateCurrentUserNickname(c *nova.Context) (err error) {
	a, us, v := authResult(c), userService(c), view.Extract(c)
	var res1 *types.UpdateUserResponse
	if res1, err = us.UpdateUser(c.Req.Context(), &types.UpdateUserRequest{
		Account:        a.User.Account,
		UpdateNickname: true,
		Nickname:       c.Req.FormValue("nickname"),
	}); err != nil {
		return
	}
	v.Data["user"] = res1.User
	v.DataAsJSON()
	return
}

func routeUpdateCurrentUserPassword(c *nova.Context) (err error) {
	a, us, v := authResult(c), userService(c), view.Extract(c)
	if _, err = us.AuthenticateUser(c.Req.Context(), &types.AuthenticateUserRequest{
		Account:  a.User.Account,
		Password: c.Req.FormValue("oldPassword"),
	}); err != nil {
		return
	}
	var res2 *types.UpdateUserResponse
	if res2, err = us.UpdateUser(c.Req.Context(), &types.UpdateUserRequest{
		Account:        a.User.Account,
		UpdatePassword: true,
		Password:       c.Req.FormValue("newPassword"),
	}); err != nil {
		return
	}
	v.Data["user"] = res2.User
	v.DataAsJSON()
	return
}

func routeUpdateUserIsAdmin(c *nova.Context) (err error) {
	us, v := userService(c), view.Extract(c)
	var res1 *types.UpdateUserResponse
	if res1, err = us.UpdateUser(c.Req.Context(), &types.UpdateUserRequest{
		Account:       c.Req.FormValue("account"),
		UpdateIsAdmin: true,
		IsAdmin:       strings.HasPrefix(strings.ToLower(strings.TrimSpace(c.Req.FormValue("is_admin"))), "t"),
	}); err != nil {
		return
	}
	v.Data["user"] = res1.User
	v.DataAsJSON()
	return
}

func routeUpdateUserIsBlocked(c *nova.Context) (err error) {
	us, v := userService(c), view.Extract(c)
	var res1 *types.UpdateUserResponse
	if res1, err = us.UpdateUser(c.Req.Context(), &types.UpdateUserRequest{
		Account:         c.Req.FormValue("account"),
		UpdateIsBlocked: true,
		IsBlocked:       strings.HasPrefix(strings.ToLower(strings.TrimSpace(c.Req.FormValue("is_blocked"))), "t"),
	}); err != nil {
		return
	}
	v.Data["user"] = res1.User
	v.DataAsJSON()
	return
}

func routeGetUser(c *nova.Context) (err error) {
	rp, us, v := router.PathParams(c), userService(c), view.Extract(c)
	var res1 *types.GetUserResponse
	if res1, err = us.GetUser(c.Req.Context(), &types.GetUserRequest{
		Account: rp.Get("account"),
	}); err != nil {
		return
	}
	v.Data["user"] = res1.User
	v.DataAsJSON()
	return
}
