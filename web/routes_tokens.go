package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func hideAuthenticationError(err *error) {
	if err == nil {
		return
	}
	if *err == nil {
		return
	}
	if s, ok := status.FromError(*err); ok {
		if s.Code() == codes.InvalidArgument {
			*err = errors.New("account not exists, or password is incorrect")
		}
	}
}

func routeCreateToken(c *nova.Context) (err error) {
	ts, us, v := tokenService(c), userService(c), view.Extract(c)
	if err = c.Req.ParseForm(); err != nil {
		return
	}
	var res1 *types.AuthenticateUserResponse
	if res1, err = us.AuthenticateUser(c.Req.Context(), &types.AuthenticateUserRequest{
		Account:  c.Req.FormValue("account"),
		Password: c.Req.FormValue("password"),
	}); err != nil {
		hideAuthenticationError(&err)
		return
	}
	var res2 *types.CreateTokenResponse
	if res2, err = ts.CreateToken(c.Req.Context(), &types.CreateTokenRequest{
		Account:     res1.User.Account,
		Description: c.Req.Header.Get("User-Agent"),
	}); err != nil {
		return
	}
	v.Data["user"] = res1.User
	v.Data["token"] = res2.Token
	v.DataAsJSON()
	return
}

func routeDestroyToken(c *nova.Context) (err error) {
	ts, v, a := tokenService(c), view.Extract(c), authResult(c)
	if err = c.Req.ParseForm(); err != nil {
		return
	}
	var res1 *types.GetTokenResponse
	if res1, err = ts.GetToken(c.Req.Context(), &types.GetTokenRequest{
		Token: c.Req.FormValue("token"),
	}); err != nil {
		return
	}
	if res1.Token.Account != a.User.Account {
		err = errors.New("not your token")
		return
	}
	// if it's current token, mark to clean current token
	if res1.Token.Token == a.Token.Token {
		c.Res.Header().Set(headerKeyAction, headerValueClearToken)
	}
	if _, err = ts.DeleteToken(c.Req.Context(), &types.DeleteTokenRequest{
		Token: res1.Token.Token,
	}); err != nil {
		return
	}
	v.DataAsJSON()
	return
}
