package daemon

import (
	"github.com/yankeguo/bastion/types"
	"golang.org/x/net/context"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/jinzhu/copier"
	"time"
	"github.com/yankeguo/bastion/utils"
)

func (d *Daemon) PutGrant(c context.Context, req *types.PutGrantRequest) (res *types.PutGrantResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	n := models.Grant{}
	copier.Copy(&n, req)
	n.ID = n.BuildID()
	n.CreatedAt = time.Now().Unix()
	if err = d.DB.Save(&n); err != nil {
		err = errFromStorm(err)
		return
	}
	return
}

func (d *Daemon) ListGrants(c context.Context, req *types.ListGrantsRequest) (res *types.ListGrantsResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	var rs []models.Grant
	if err = d.DB.Find("Account", req.Account, &rs); err != nil {
		err = errFromStorm(err)
		return
	}
	ret := make([]*types.Grant, 0, len(rs))
	for _, n := range rs {
		ret = append(ret, n.ToGRPCGrant())
	}
	res = &types.ListGrantsResponse{Grants: ret, Now: time.Now().Unix()}
	return
}

func (d *Daemon) DeleteGrant(c context.Context, req *types.DeleteGrantRequest) (res *types.DeleteGrantResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	n := models.Grant{}
	copier.Copy(&n, req)
	if err = d.DB.Delete("Grant", n.BuildID()); err != nil {
		err = errFromStorm(err)
		return
	}
	return
}

func (d *Daemon) CheckGrant(c context.Context, req *types.CheckGrantRequest) (res *types.CheckGrantResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	var rs []models.Grant
	if err = d.DB.Find("Account", req.Account, &rs); err != nil {
		err = errFromStorm(err)
		return
	}
	var ok bool
	now := time.Now().Unix()
	for _, n := range rs {
		if n.User == req.User && (n.ExpiredAt == 0 || n.ExpiredAt > now) {
			if utils.MatchAsterisk(n.HostnamePattern, req.Hostname) {
				ok = true
				break
			}
		}
	}
	res = &types.CheckGrantResponse{Ok: ok}
	return
}

func (d *Daemon) ListGrantItems(c context.Context, req *types.ListGrantItemsRequest) (res *types.ListGrantItemsResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	var rs []models.Grant
	if err = d.DB.Find("Account", req.Account, &rs); err != nil {
		err = errFromStorm(err)
		return
	}
	var ns []models.Node
	if err = d.DB.All(&ns); err != nil {
		err = errFromStorm(err)
		return
	}
	now := time.Now().Unix()
	ret := make([]*types.GrantItem, 0)
	for _, n := range ns {
		for _, r := range rs {
			if utils.MatchAsterisk(r.HostnamePattern, n.Hostname) && (r.ExpiredAt == 0 || r.ExpiredAt > now) {
				ret = append(ret, &types.GrantItem{
					Hostname:  n.Hostname,
					User:      r.User,
					ExpiredAt: r.ExpiredAt,
				})
			}
		}
	}
	res = &types.ListGrantItemsResponse{GrantItems: compactGrantItems(ret)}
	return
}

func compactGrantItems(is []*types.GrantItem) []*types.GrantItem {
	ret := make([]*types.GrantItem, 0, len(is))
	for _, i := range is {
		var found bool
		// find and update existing value in ret
		for _, r := range ret {
			if r.Hostname == i.Hostname && r.User == i.User {
				if i.ExpiredAt == 0 {
					r.ExpiredAt = 0
				} else if i.ExpiredAt > r.ExpiredAt {
					r.ExpiredAt = i.ExpiredAt
				}
				found = true
			}
		}
		// append ret if not found
		if !found {
			ret = append(ret, i)
		}
	}
	return ret
}
