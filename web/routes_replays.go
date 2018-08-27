package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/utils"
	"io"
	"net/http"
	"strconv"
)

func routeDownloadReplay(c *nova.Context) (err error) {
	rs, qp := replayService(c), router.PathParams(c)
	id, _ := strconv.ParseInt(qp.Get("id"), 10, 64)
	var sess types.ReplayService_ReadReplayClient
	if sess, err = rs.ReadReplay(c.Req.Context(), &types.ReadReplayRequest{SessionId: id}); err != nil {
		return
	}
	c.Res.Header().Set("Content-Type", "application/octet-stream")
	c.Res.WriteHeader(http.StatusOK)
	for {
		var f *types.ReplayFrame
		if f, err = sess.Recv(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		if err = utils.WriteReplayFrame(f, c.Res); err != nil {
			return
		}
	}
	return
}

func routePageReplay(c *nova.Context) (err error) {
	v, ar := view.Extract(c), router.PathParams(c)
	v.Data["SessionId"] = ar.Get("id")
	v.Data["ViewKey"] = c.Req.FormValue("viewKey")
	v.HTML("replay")
	return
}
