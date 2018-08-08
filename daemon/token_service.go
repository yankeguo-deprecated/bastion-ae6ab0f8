package daemon

import (
	"github.com/yankeguo/bastion/types"
	"golang.org/x/net/context"
	"github.com/yankeguo/bastion/daemon/models"
)

func (d *Daemon) CreateToken(c context.Context, req *types.CreateTokenRequest) (res *types.CreateTokenResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	t := models.Token{}
	if err = d.DB.Tx(true, func(db *Node) (err error) {
		var u models.User
		if err = db.One("Account", req.Account, &u); err != nil {
			return
		}
		t.Account = u.Account
		t.Description = req.Description
		t.Token = newToken()
		t.CreatedAt = now()
		if err = db.Save(&t); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}
	res = &types.CreateTokenResponse{Token: t.ToGRPCToken()}
	return
}

func (d *Daemon) GetToken(c context.Context, req *types.GetTokenRequest) (res *types.GetTokenResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	t := models.Token{}
	if len(req.Token) > 0 {
		if err = d.DB.One("Token", req.Token, &t); err != nil {
			return
		}
	} else {
		if err = d.DB.One("Id", req.Id, &t); err != nil {
			return
		}
	}
	res = &types.GetTokenResponse{Token: t.ToGRPCToken()}
	return
}

func (d *Daemon) TouchToken(c context.Context, req *types.TouchTokenRequest) (res *types.TouchTokenResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	t := models.Token{}
	if err = d.DB.Tx(true, func(db *Node) (err error) {
		if err = db.One("Token", req.Token, &t); err != nil {
			return
		}
		t.ViewedAt = now()
		if err = db.Save(&t); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	res = &types.TouchTokenResponse{Token: t.ToGRPCToken()}
	return
}

func (d *Daemon) ListTokens(c context.Context, req *types.ListTokensRequest) (res *types.ListTokensResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	var ts []models.Token
	if err = d.DB.Find("Account", req.Account, &ts); err != nil {
		return
	}
	ret := make([]*types.Token, 0, len(ts))
	for _, t := range ts {
		// hide the actual token
		nt := t.ToGRPCToken()
		nt.Token = ""
		ret = append(ret, nt)
	}
	res = &types.ListTokensResponse{Tokens: ret}
	return
}

func (d *Daemon) DeleteToken(c context.Context, req *types.DeleteTokenRequest) (res *types.DeleteTokenResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	if err = d.DB.DeleteStruct(&models.Token{Id: req.Id}); err != nil {
		return
	}
	res = &types.DeleteTokenResponse{}
	return
}
