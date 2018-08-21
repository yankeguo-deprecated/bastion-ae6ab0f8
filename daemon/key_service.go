package daemon

import (
	"github.com/asdine/storm/q"
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
	if err = d.db.Find("Account", req.Account, &keys); err != nil {
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
	if err = d.db.Tx(true, func(db *Node) (err error) {
		// check duplicated
		if err = db.CheckDuplicated("Key", "fingerprint", req.Fingerprint); err != nil {
			return
		}
		// delete existed sandbox keys
		if req.Source == types.KeySourceSandbox {
			if err = db.Select(q.Eq("Source", types.KeySourceSandbox), q.Eq("Account", req.Account)).Delete(new(models.Key)); err != nil {
				return
			}
		}
		// copy and save new key
		copier.Copy(&k, req)
		k.CreatedAt = now()
		if err = db.Save(&k); err != nil {
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
	err = d.db.DeleteStruct(&models.Key{Fingerprint: req.Fingerprint})
	res = &types.DeleteKeyResponse{}
	return
}

func (d *Daemon) GetKey(c context.Context, req *types.GetKeyRequest) (res *types.GetKeyResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	k := models.Key{}
	if err = d.db.One("Fingerprint", req.Fingerprint, &k); err != nil {
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
	if err = d.db.Tx(true, func(db *Node) (err error) {
		if err = db.One("Fingerprint", req.Fingerprint, &k); err != nil {
			return
		}
		k.ViewedAt = now()
		if err = db.Save(&k); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}
	res = &types.TouchKeyResponse{Key: k.ToGRPCKey()}
	return
}
