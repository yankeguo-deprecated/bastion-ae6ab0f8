package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"strings"
)

func mountRoutes(n *nova.Nova) {
	router.Route(n).Get("/api/check").Use(routeCheck)
	router.Route(n).Get("/api/authorized_keys").Use(routeAuthorizedKeys)
	router.Route(n).Post("/api/tokens/create").Use(routeCreateToken)
	router.Route(n).Get("/api/tokens").Use(
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
	router.Route(n).Post("/api/users/current/update_password").Use(
		requiresLoggedIn(false),
		routeUpdateCurrentUserPassword,
	)
	router.Route(n).Get("/api/users/current/grant_items").Use(
		requiresLoggedIn(false),
		routeGetCurrentUserGrantItems,
	)
	router.Route(n).Get("/api/users/current/keys").Use(
		requiresLoggedIn(false),
		routeListKeys,
	)
	router.Route(n).Post("/api/users/current/keys/create").Use(
		requiresLoggedIn(false),
		routeCreateKey,
	)
	router.Route(n).Post("/api/keys/destroy").Use(
		requiresLoggedIn(false),
		routeDestroyKey,
	)
	router.Route(n).Get("/api/nodes").Use(
		requiresLoggedIn(true),
		routeListNodes,
	)
	router.Route(n).Post("/api/nodes/create").Use(
		requiresLoggedIn(true),
		routeCreateNode,
	)
	router.Route(n).Post("/api/nodes/destroy").Use(
		requiresLoggedIn(true),
		routeDestroyNode,
	)
	router.Route(n).Post("/api/nodes/update_is_key_managed").Use(
		requiresLoggedIn(true),
		routeUpdateNodeIsKeyManaged,
	)
	router.Route(n).Get("/api/users").Use(
		requiresLoggedIn(true),
		routeListUsers,
	)
	router.Route(n).Post("/api/users/create").Use(
		requiresLoggedIn(true),
		routeCreateUser,
	)
	router.Route(n).Post("/api/users/update_is_admin").Use(
		requiresLoggedIn(true),
		routeUpdateUserIsAdmin,
	)
	router.Route(n).Post("/api/users/update_is_blocked").Use(
		requiresLoggedIn(true),
		routeUpdateUserIsBlocked,
	)
	router.Route(n).Post("/api/users/update_nickname").Use(
		requiresLoggedIn(true),
		routeUpdateUserNickname,
	)
	router.Route(n).Get("/api/users/:account").Use(
		requiresLoggedIn(true),
		routeGetUser,
	)
	router.Route(n).Get("/api/users/:account/grants").Use(
		requiresLoggedIn(true),
		routeGetGrants,
	)
	router.Route(n).Post("/api/users/:account/grants/create").Use(
		requiresLoggedIn(true),
		routeCreateGrant,
	)
	router.Route(n).Post("/api/users/:account/grants/destroy").Use(
		requiresLoggedIn(true),
		routeDestroyGrant,
	)
	router.Route(n).Get("/api/sessions").Use(
		requiresLoggedIn(true),
		routeListSessions,
	)
	router.Route(n).Get("/api/sessions/:id").Use(
		requiresLoggedIn(true),
		routeGetSession,
	)
	router.Route(n).Get("/api/replays/:id/download").Use(
		requiresLoggedIn(true),
		routeDownloadReplay,
	)
	router.Route(n).Get("/replays/:id").Use(routePageReplay)
}

func routeCheck(c *nova.Context) error {
	v := view.Extract(c)
	v.Data["ok"] = true
	v.DataAsJSON()
	return nil
}

func routeAuthorizedKeys(c *nova.Context) (err error) {
	mks, v := masterKeyService(c), view.Extract(c)
	var res *types.ListMasterKeysResponse
	if res, err = mks.ListMasterKeys(c.Req.Context(), &types.ListMasterKeysRequest{}); err != nil {
		return
	}
	var out string
	for _, k := range res.MasterKeys {
		out = out + strings.TrimSpace(k.PublicKey) + "\n"
	}
	v.Text(out)
	return
}
