package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
)

func mountRoutes(n *nova.Nova) {
	router.Route(n).Get("/api/check").Use(routeCheck)
	router.Route(n).Post("/api/tokens/create").Use(routeCreateToken)
	router.Route(n).Post("/api/tokens/:token/destroy").Use(
		requiresLoggedIn(false),
		routeDestroyToken,
	)
	router.Route(n).Get("/api/users/current").Use(
		requiresLoggedIn(false),
		routeGetCurrentUser,
	)
}

func routeCheck(c *nova.Context) error {
	v := view.Extract(c)
	v.Data["ok"] = true
	v.DataAsJSON()
	return nil
}
