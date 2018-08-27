package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"strconv"
)

func routeGetSession(c *nova.Context) (err error) {
	v, ss, us, pr := view.Extract(c), sessionService(c), userService(c), router.PathParams(c)
	id, _ := strconv.ParseInt(pr.Get("id"), 10, 64)
	var res1 *types.GetSessionResponse
	if res1, err = ss.GetSession(c.Req.Context(), &types.GetSessionRequest{Id: id}); err != nil {
		return
	}
	var res2 *types.GetUserResponse
	if res2, err = us.GetUser(c.Req.Context(), &types.GetUserRequest{Account: res1.Session.Account}); err != nil {
		return
	}
	v.Data["session"] = res1.Session
	v.Data["user"] = res2.User
	v.DataAsJSON()
	return
}

func routeListSessions(c *nova.Context) (err error) {
	skip, _ := strconv.ParseInt(c.Req.FormValue("skip"), 10, 64)
	limit, _ := strconv.ParseInt(c.Req.FormValue("limit"), 10, 64)
	v, ss := view.Extract(c), sessionService(c)
	var res *types.ListSessionsResponse
	if res, err = ss.ListSessions(c.Req.Context(), &types.ListSessionsRequest{
		Skip:  int32(skip),
		Limit: int32(limit),
	}); err != nil {
		return
	}
	v.Data["sessions"] = res.Sessions
	v.Data["skip"] = res.Skip
	v.Data["limit"] = res.Limit
	v.Data["total"] = res.Total
	v.DataAsJSON()
	return
}
