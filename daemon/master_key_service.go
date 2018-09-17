package daemon

import (
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/net/context"
)

func (d *Daemon) ListMasterKeys(ctx context.Context, req *types.ListMasterKeysRequest) (res *types.ListMasterKeysResponse, err error) {
	var mKeys []models.MasterKey
	if err = d.db.All(&mKeys); err != nil {
		return
	}
	ret := make([]*types.MasterKey, 0, len(mKeys))
	for _, k := range mKeys {
		ret = append(ret, k.ToGRPCModel())
	}
	res = &types.ListMasterKeysResponse{MasterKeys: ret}
	return
}

func (d *Daemon) UpdateAllMasterKeys(ctx context.Context, req *types.UpdateAllMasterKeysRequest) (res *types.UpdateAllMasterKeysResponse, err error) {
	if err = d.db.Tx(true, func(db *Node) (err error) {
		// find existing master keys
		var mKeys []models.MasterKey
		if err = db.All(&mKeys); err != nil {
			return
		}
		// delete all master keys
		for _, k := range mKeys {
			if err = db.DeleteStruct(&k); err != nil {
				return
			}
		}
		// save new master keys
		for _, k := range req.MasterKeys {
			if err = db.Save(&models.MasterKey{
				Fingerprint: k.Fingerprint,
				PublicKey:   k.PublicKey,
			}); err != nil {
				return
			}
		}
		return
	}); err != nil {
		return
	}
	res = &types.UpdateAllMasterKeysResponse{}
	return
}
