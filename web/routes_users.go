package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/view"
)

func routeGetCurrentUser(c *nova.Context) (err error) {
	a, v := authResult(c), view.Extract(c)
	v.Data["user"] = a.User
	v.DataAsJSON()
	return
}
