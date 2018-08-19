package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"strconv"
)

func routeListSessions(c *nova.Context) (err error) {
	skip, _ := strconv.ParseInt(c.Req.FormValue("skip"), 10, 8)
	limit, _ := strconv.ParseInt(c.Req.FormValue("limit"), 10, 8)
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
