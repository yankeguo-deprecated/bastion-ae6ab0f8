package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/view"
	"github.com/pkg/errors"
	"github.com/yankeguo/bastion/types"
)

func routeListNodes(c *nova.Context) (err error) {
	ns, v := nodeService(c), view.Extract(c)
	var res1 *types.ListNodesResponse
	if res1, err = ns.ListNodes(c.Req.Context(), &types.ListNodesRequest{}); err != nil {
		return
	}
	v.Data["nodes"] = res1.Nodes
	v.DataAsJSON()
	return
}

func routeCreateNode(c *nova.Context) (err error) {
	ns, v := nodeService(c), view.Extract(c)
	var res1 *types.PutNodeResponse
	if res1, err = ns.PutNode(
		c.Req.Context(),
		&types.PutNodeRequest{
			Hostname: c.Req.FormValue("hostname"),
			Address:  c.Req.FormValue("address"),
			Source:   types.NodeSourceManual,
		}); err != nil {
		return
	}
	v.Data["node"] = res1.Node
	v.DataAsJSON()
	return
}

func routeDestroyNode(c *nova.Context) (err error) {
	ns, v := nodeService(c), view.Extract(c)
	var res1 *types.GetNodeResponse
	if res1, err = ns.GetNode(c.Req.Context(), &types.GetNodeRequest{
		Hostname: c.Req.FormValue("hostname"),
	}); err != nil {
		return
	}
	if res1.Node.Source != types.NodeSourceManual {
		return errors.New("only manually added node can be deleted")
		return
	}
	if _, err = ns.DeleteNode(c.Req.Context(), &types.DeleteNodeRequest{
		Hostname: res1.Node.Hostname,
	}); err != nil {
		return
	}
	v.DataAsJSON()
	return
}
