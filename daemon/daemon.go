package daemon

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/pkg/errors"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
)

var (
	ErrDaemonAlreadyRunning = errors.New("daemon is already running")
)

// Daemon daemon instance
type Daemon struct {
	DB       *storm.DB
	Listener net.Listener
	Server   *grpc.Server

	opts types.DaemonOptions
}

func New(opts types.DaemonOptions) *Daemon {
	return &Daemon{opts: opts}
}

func (d *Daemon) Transaction(writable bool, cb func(storm.Node) error) (err error) {
	var tx storm.Node
	if tx, err = d.DB.Begin(writable); err != nil {
		err = errFromStorm(err)
		return
	}
	defer tx.Rollback()
	if err = cb(tx); err != nil {
		return
	}
	if err = tx.Commit(); err != nil {
		err = errFromStorm(err)
		return
	}
	return
}

func (d *Daemon) Run() (err error) {
	defer d.cleanup()
	// ensure database directory
	os.MkdirAll(filepath.Dir(d.opts.DB), 0640)
	// open database
	if d.DB != nil {
		err = ErrDaemonAlreadyRunning
		return
	}
	if d.DB, err = storm.Open(d.opts.DB); err != nil {
		return
	}
	// migrate database
	for _, m := range models.AllModels {
		if err = d.DB.Init(m); err != nil {
			return
		}
	}
	// create listener
	if d.Listener != nil {
		err = ErrDaemonAlreadyRunning
		return
	}
	if d.Listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", d.opts.Host, d.opts.Port)); err != nil {
		return
	}
	// create server
	if d.Server != nil {
		err = ErrDaemonAlreadyRunning
		return
	}
	d.Server = grpc.NewServer()
	types.RegisterUserServiceServer(d.Server, d)
	types.RegisterNodeServiceServer(d.Server, d)
	// serve
	return d.Server.Serve(d.Listener)
}

func (d *Daemon) cleanup() {
	if d.Server != nil {
		d.Server.GracefulStop()
		d.Server = nil
	}
	if d.Listener != nil {
		d.Listener.Close()
		d.Listener = nil
	}
	if d.DB != nil {
		d.DB.Close()
		d.DB = nil
	}
}

func (d *Daemon) Shutdown() (err error) {
	if d.Server != nil {
		d.Server.GracefulStop()
		d.Server = nil
	}
	return
}
