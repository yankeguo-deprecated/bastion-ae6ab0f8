package daemon

import (
	"github.com/asdine/storm"
	"github.com/jinzhu/copier"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/net/context"
	"strings"
	"time"
)

func (d *Daemon) ListKeys(c context.Context, req *types.ListKeysRequest) (res *types.ListKeysResponse, err error) {
	var keys []models.Key
	if err = d.DB.Find("Account", strings.TrimSpace(req.Account), &keys); err != nil {
		err = errFromStorm(err)
		return
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
		k.CreatedAt = time.Now().Unix()
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
		return
	}
	res = &types.GetKeyResponse{Key: k.ToGRPCKey()}
	return
}
