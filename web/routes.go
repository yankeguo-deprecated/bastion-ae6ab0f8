package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
)

func mountRoutes(n *nova.Nova) {
	router.Route(n).Get("/api/check").Use(routeCheck)
	router.Route(n).Post("/api/tokens/create").Use(routeCreateToken)
	router.Route(n).Get("/api/tokens/list").Use(
		requiresLoggedIn(false),
		routeListTokens,
	)
	router.Route(n).Post("/api/tokens/destroy").Use(
		requiresLoggedIn(false),
		routeDestroyToken,
	)
	router.Route(n).Get("/api/users/current").Use(
		requiresLoggedIn(false),
		routeGetCurrentUser,
	)
	router.Route(n).Post("/api/users/current/update_nickname").Use(
		requiresLoggedIn(false),
		routeUpdateCurrentUserNickname,
	)
	router.Route(n).Get("/api/users/current/grant_items").Use(
		requiresLoggedIn(false),
		routeGetCurrentUserGrantItems,
	)
}

func routeCheck(c *nova.Context) error {
	v := view.Extract(c)
	v.Data["ok"] = true
	v.DataAsJSON()
	return nil
}
