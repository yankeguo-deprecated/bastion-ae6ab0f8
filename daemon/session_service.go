package daemon

import (
	"github.com/asdine/storm"
	"github.com/jinzhu/copier"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/net/context"
)

func (d *Daemon) CreateSession(c context.Context, req *types.CreateSessionRequest) (res *types.CreateSessionResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	s := models.Session{}
	copier.Copy(&s, req)
	s.CreatedAt = now()
	if err = d.DB.Save(&s); err != nil {
		err = errFromStorm(err)
		return
	}
	res = &types.CreateSessionResponse{Session: s.ToGRPCSession()}
	return
}

func (d *Daemon) FinishSession(c context.Context, req *types.FinishSessionRequest) (res *types.FinishSessionResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	s := models.Session{}
	if err = d.DB.One("Id", req.Id, &s); err != nil {
		err = errFromStorm(err)
		return
	}
	s.FinishedAt = now()
	if err = d.DB.Save(&s); err != nil {
		err = errFromStorm(err)
		return
	}
	res = &types.FinishSessionResponse{Session: s.ToGRPCSession()}
	return
}

func (d *Daemon) ListSessions(c context.Context, req *types.ListSessionsRequest) (res *types.ListSessionsResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	var sessions []models.Session
	var total int
	if err = d.Tx(false, func(db storm.Node) (err error) {
		if total, err = db.Count(new(models.Session)); err != nil {
			err = errFromStorm(err)
			return
		}
		if err = db.All(&sessions, storm.Reverse(), storm.Skip(int(req.Skip)), storm.Limit(int(req.Limit))); err != nil {
			err = errFromStorm(err)
			return
		}
		return
	}); err != nil {
		return
	}
	ret := make([]*types.Session, 0, len(sessions))
	for _, s := range sessions {
		ret = append(ret, s.ToGRPCSession())
	}
	res = &types.ListSessionsResponse{
		Skip:     req.Skip,
		Limit:    req.Limit,
		Total:    int32(total),
		Sessions: ret,
	}
	return
}
