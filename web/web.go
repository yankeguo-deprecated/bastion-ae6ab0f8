package web

import (
	"fmt"
	"github.com/novakit/nova"
	"github.com/novakit/static"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"net/http"
)

func NewServer(opts types.WebOptions) *http.Server {
	n := nova.New()
	// mount static module
	n.Use(static.Handler(static.Options{
		Directory: "public",
		BinFS:     !opts.Dev,
	}))
	// mount view module for json rendering only
	n.Use(view.Handler(view.Options{
		Directory: "views",
		BinFS:     !opts.Dev,
	}))
	// mount rpc module
	n.Use(rpcModule(opts))
	// mount all routes
	mountRoutes(n)
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler: n,
	}
}
