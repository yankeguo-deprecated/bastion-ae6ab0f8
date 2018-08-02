package daemon

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
)

// Daemon daemon instance
type Daemon struct {
	DB     *storm.DB
	Server *grpc.Server

	opts types.DaemonOptions
}

func New(opts types.DaemonOptions) *Daemon {
	return &Daemon{opts: opts}
}

func (d *Daemon) Tx(writable bool, cb func(storm.Node) error) (err error) {
	var tx storm.Node
	if tx, err = d.DB.Begin(writable); err != nil {
		err = errFromStorm(err)
		return
	}
	defer tx.Rollback()
	if err = cb(tx); err != nil {
		return
	}
	if writable {
		if err = tx.Commit(); err != nil {
			err = errFromStorm(err)
			return
		}
	}
	return
}

func (d *Daemon) Run() (err error) {
	// open db
	if d.DB, err = d.openDB(); err != nil {
		return
	}
	defer d.DB.Close()

	// create listener
	var l net.Listener
	if l, err = d.createListener(); err != nil {
		return
	}
	defer l.Close()

	// create server
	d.Server = d.createGRPCServer()

	// run server
	if err = d.Server.Serve(l); err != nil {
		if err == grpc.ErrServerStopped {
			err = nil
		}
	}
	return
}

func (d *Daemon) openDB() (db *storm.DB, err error) {
	// ensure database directory
	os.MkdirAll(filepath.Dir(d.opts.DB), 0640)
	// open db
	if db, err = storm.Open(d.opts.DB); err != nil {
		return
	}
	// migrate database
	for _, m := range models.AllModels {
		if err = db.Init(m); err != nil {
			db.Close()
			return
		}
	}
	return
}

func (d *Daemon) createListener() (l net.Listener, err error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", d.opts.Host, d.opts.Port))
}

func (d *Daemon) createGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	types.RegisterUserServiceServer(s, d)
	types.RegisterNodeServiceServer(s, d)
	types.RegisterKeyServiceServer(s, d)
	types.RegisterGrantServiceServer(s, d)
	types.RegisterSessionServiceServer(s, d)
	return s
}

func (d *Daemon) Stop() {
	if d.Server != nil {
		d.Server.GracefulStop()
	}
	return
}
