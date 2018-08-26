package web

import (
	"fmt"
	"github.com/felixge/httpsnoop"
	_ "github.com/novakit/binfs"
	"github.com/novakit/nova"
	"github.com/novakit/static"
	"github.com/novakit/view"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc/status"
	"net/http"
)

func NewServer(opts types.WebOptions) *http.Server {
	n := nova.New()
	if opts.Dev {
		n.Env = nova.Development
	} else {
		n.Env = nova.Production
	}
	n.Error(func(c *nova.Context, err error) {
		log.Error().Err(err).Str("method", c.Req.Method).Str("path", c.Req.URL.Path).Msg("error occurred")
		if s, ok := status.FromError(err); ok {
			// if it's a grpc status error, extract description
			http.Error(c.Res, s.Message(), http.StatusInternalServerError)
		} else {
			// just render any error as 500 and expose the message
			http.Error(c.Res, err.Error(), http.StatusInternalServerError)
		}
	})
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
	// mount opts module
	n.Use(optsModule(opts))
	// mount rpc module
	n.Use(rpcModule(opts))
	// mount auth module
	n.Use(authModule())
	// mount all routes
	mountRoutes(n)
	// build the http.Server
	return &http.Server{
		Addr: fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(n, w, r)
			log.Info().Str("method", r.Method).Str("url", r.URL.String()).Int("code", m.Code).Dur("duration", m.Duration).Int64("written", m.Written).Msg("request handled")
		}),
	}
}
