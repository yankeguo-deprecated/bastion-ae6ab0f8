package daemon

import (
	"github.com/asdine/storm"
	"github.com/jinzhu/copier"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/net/context"
)

func (d *Daemon) ListKeys(c context.Context, req *types.ListKeysRequest) (res *types.ListKeysResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	var keys []models.Key
	if err = d.DB.Find("Account", req.Account, &keys); err != nil {
		if err == storm.ErrNotFound {
			err = nil
		} else {
			err = errFromStorm(err)
			return
		}
	}
	ret := make([]*types.Key, 0, len(keys))
	for _, k := range keys {
		ret = append(ret, k.ToGRPCKey())
	}
	res = &types.ListKeysResponse{Keys: ret}
	return
}

func (d *Daemon) CreateKey(c context.Context, req *types.CreateKeyRequest) (res *types.CreateKeyResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	k := models.Key{}
	if err = d.Tx(true, func(db storm.Node) (err error) {
		if err = checkDuplicated(db, "Key", "fingerprint", req.Fingerprint); err != nil {
			return
		}
		copier.Copy(&k, req)
		k.CreatedAt = now()
		if err = db.Save(&k); err != nil {
			err = errFromStorm(err)
			return
		}
		return
	}); err != nil {
		return
	}
	res = &types.CreateKeyResponse{Key: k.ToGRPCKey()}
	return
}

func (d *Daemon) DeleteKey(c context.Context, req *types.DeleteKeyRequest) (res *types.DeleteKeyResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	err = errFromStorm(d.DB.DeleteStruct(&models.Key{Fingerprint: req.Fingerprint}))
	res = &types.DeleteKeyResponse{}
	return
}

func (d *Daemon) GetKey(c context.Context, req *types.GetKeyRequest) (res *types.GetKeyResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	k := models.Key{}
	if err = d.DB.One("Fingerprint", req.Fingerprint, &k); err != nil {
		err = errFromStorm(err)
		return
	}
	res = &types.GetKeyResponse{Key: k.ToGRPCKey()}
	return
}

func (d *Daemon) TouchKey(c context.Context, req *types.TouchKeyRequest) (res *types.TouchKeyResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	k := models.Key{}
	if err = d.Tx(true, func(db storm.Node) (err error) {
		if err = db.One("Fingerprint", req.Fingerprint, &k); err != nil {
			err = errFromStorm(err)
			return
		}
		k.ViewedAt = now()
		if err = db.Save(&k); err != nil {
			err = errFromStorm(err)
			return
		}
		return
	}); err != nil {
		return
	}
	res = &types.TouchKeyResponse{Key: k.ToGRPCKey()}
	return
}
