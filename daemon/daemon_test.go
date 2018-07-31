package daemon

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
)

var d *Daemon

func temporaryFile() string {
	buf := make([]byte, 8, 8)
	rand.Read(buf)
	return filepath.Join(os.TempDir(), "bnktestdb"+hex.EncodeToString(buf)+".bolt")
}

func withDaemon(t *testing.T, cb func(*testing.T, *Daemon, *grpc.ClientConn)) {
	d = New(types.DaemonOptions{
		DB:   temporaryFile(),
		Host: "127.0.0.1",
		Port: 1997,
	})
	go d.Run()
	defer d.Shutdown()
	time.Sleep(time.Second / 2)
	var err error
	var c *grpc.ClientConn
	if c, err = grpc.Dial("127.0.0.1:1997", grpc.WithInsecure()); err != nil {
		t.Error(err)
		return
	}
	defer c.Close()
	cb(t, d, c)
}
