package daemon

import (
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/daemon/db"
	"google.golang.org/grpc"
	"net"
	"fmt"
	"github.com/pkg/errors"
)

// Daemon daemon instance
type Daemon struct {
	DB *db.DB

	Server *grpc.Server

	opts types.DaemonOptions
}

func New(opts types.DaemonOptions) *Daemon {
	return &Daemon{opts: opts}
}

func (d *Daemon) Run() (err error) {
	// open database
	if d.DB, err = db.Open(d.opts.DB); err != nil {
		return
	}
	// migrate
	if err = d.DB.Migrate(); err != nil {
		return
	}
	// create listener
	var l net.Listener
	if l, err = net.Listen("tcp", fmt.Sprintf("%s:%d", d.opts.Host, d.opts.Port)); err != nil {
		return
	}
	// create d.Server
	if d.Server != nil {
		err = errors.New("daemon is already started")
		return
	}
	d.Server = grpc.NewServer()
	types.RegisterUserServiceServer(d.Server, d)
	return d.Server.Serve(l)
}

func (d *Daemon) Shutdown() (err error) {
	if d.Server != nil {
		d.Server.GracefulStop()
		d.Server = nil
	}
	return
}
