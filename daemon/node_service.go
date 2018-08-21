package daemon

import (
	"github.com/jinzhu/copier"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/net/context"
	"strings"
)

func (d *Daemon) ListNodes(c context.Context, req *types.ListNodesRequest) (res *types.ListNodesResponse, err error) {
	var ns []models.Node
	if err = d.db.All(&ns); err != nil {
		return
	}
	ret := make([]*types.Node, 0, len(ns))

	for _, n := range ns {
		ret = append(ret, n.ToGRPCNode())
	}
	res = &types.ListNodesResponse{Nodes: ret}
	return
}

func (d *Daemon) PutNode(c context.Context, req *types.PutNodeRequest) (res *types.PutNodeResponse, err error) {
	// fix request
	if err = req.Validate(); err != nil {
		return
	}
	// create node
	n := models.Node{}
	copier.Copy(&n, req)
	n.CreatedAt = now()
	if err = d.db.Save(&n); err != nil {
		return
	}
	// build response
	res = &types.PutNodeResponse{Node: n.ToGRPCNode()}
	return
}

func (d *Daemon) DeleteNode(c context.Context, req *types.DeleteNodeRequest) (res *types.DeleteNodeResponse, err error) {
	req.Hostname = strings.TrimSpace(req.Hostname)
	res = &types.DeleteNodeResponse{}
	err = d.db.DeleteStruct(&models.Node{Hostname: req.Hostname})
	return
}

func (d *Daemon) GetNode(c context.Context, req *types.GetNodeRequest) (res *types.GetNodeResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	n := models.Node{}
	if err = d.db.One("Hostname", req.Hostname, &n); err != nil {
		return
	}
	res = &types.GetNodeResponse{Node: n.ToGRPCNode()}
	return
}

func (d *Daemon) TouchNode(c context.Context, req *types.TouchNodeRequest) (res *types.TouchNodeResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}
	n := models.Node{}
	if err = d.db.Tx(true, func(db *Node) (err error) {
		if err = db.One("Hostname", req.Hostname, &n); err != nil {
			return
		}
		n.ViewedAt = now()
		if err = db.Save(&n); err != nil {
			return
		}
		return
	}); err != nil {
		return
	}
	res = &types.TouchNodeResponse{Node: n.ToGRPCNode()}
	return
}
