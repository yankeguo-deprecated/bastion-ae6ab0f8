package daemon

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/yankeguo/bastion/daemon/models"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"net"
	"os"
	"path/filepath"
)

// Daemon daemon instance
type Daemon struct {
	opts   types.DaemonOptions
	db     *DB
	server *grpc.Server
}

func New(opts types.DaemonOptions) *Daemon {
	return &Daemon{opts: opts}
}

func (d *Daemon) initDB() (err error) {
	// ensure database directory
	os.MkdirAll(filepath.Dir(d.opts.DB), 0640)
	// open db
	var stormDB *storm.DB
	if stormDB, err = storm.Open(d.opts.DB); err != nil {
		return
	}
	// migrate database
	for _, m := range models.AllModels {
		if err = stormDB.Init(m); err != nil {
			stormDB.Close()
			return
		}
	}
	d.db = &DB{db: stormDB}
	return
}

func (d *Daemon) createListener() (l net.Listener, err error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", d.opts.Host, d.opts.Port))
}

func (d *Daemon) createGRPCServer() *grpc.Server {
	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(zerologUnaryInterceptor, grpc_recovery.UnaryServerInterceptor()),
		grpc_middleware.WithStreamServerChain(zerologStreamInterceptor, grpc_recovery.StreamServerInterceptor()),
	)
	types.RegisterUserServiceServer(s, d)
	types.RegisterNodeServiceServer(s, d)
	types.RegisterKeyServiceServer(s, d)
	types.RegisterGrantServiceServer(s, d)
	types.RegisterSessionServiceServer(s, d)
	types.RegisterTokenServiceServer(s, d)
	types.RegisterReplayServiceServer(s, d)
	return s
}

func (d *Daemon) Run() (err error) {
	// open db
	if err = d.initDB(); err != nil {
		return
	}
	defer d.db.Close()

	// create listener
	var l net.Listener
	if l, err = d.createListener(); err != nil {
		return
	}
	defer l.Close()

	// create server
	d.server = d.createGRPCServer()

	// run server
	if err = d.server.Serve(l); err != nil {
		if err == grpc.ErrServerStopped {
			err = nil
		}
	}
	return
}

func (d *Daemon) Stop() {
	if d.server != nil {
		d.server.GracefulStop()
	}
	return
}
