package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"github.com/novakit/router"
	"strconv"
	"time"
)

func routeGetGrants(c *nova.Context) (err error) {
	gs, v, rp := grantService(c), view.Extract(c), router.PathParams(c)
	var res1 *types.ListGrantsResponse
	if res1, err = gs.ListGrants(c.Req.Context(), &types.ListGrantsRequest{
		Account: rp.Get("account"),
	}); err != nil {
		return
	}
	v.Data["grants"] = res1.Grants
	v.DataAsJSON()
	return
}

func routeCreateGrant(c *nova.Context) (err error) {
	gs, v, rp := grantService(c), view.Extract(c), router.PathParams(c)
	expiresIn, _ := strconv.ParseInt(c.Req.FormValue("expires_in"), 10, 64)
	var expiresAt int64
	if expiresIn == 0 {
		expiresAt = 0
	} else {
		expiresAt = time.Now().Unix() + expiresIn
	}
	var res1 *types.PutGrantResponse
	if res1, err = gs.PutGrant(c.Req.Context(), &types.PutGrantRequest{
		Account:         rp.Get("account"),
		User:            c.Req.FormValue("user"),
		HostnamePattern: c.Req.FormValue("hostname_pattern"),
		ExpiredAt:       expiresAt,
	}); err != nil {
		return
	}
	v.Data["grant"] = res1.Grant
	v.DataAsJSON()
	return
}

func routeDestroyGrant(c *nova.Context) (err error) {
	gs, v, rp := grantService(c), view.Extract(c), router.PathParams(c)
	if _, err = gs.DeleteGrant(c.Req.Context(), &types.DeleteGrantRequest{
		Account:         rp.Get("account"),
		User:            c.Req.FormValue("user"),
		HostnamePattern: c.Req.FormValue("hostname_pattern"),
	}); err != nil {
		return
	}
	v.DataAsJSON()
	return
}
