package binfs // import "github.com/novakit/binfs"

import (
	"errors"
	"io"
	"net/http"
	"os"
)

// ErrIsDirectory error returned while trying read/seek a directory
var ErrIsDirectory = errors.New("is a directory")

// File abstracts a binfs file
type File interface {
	http.File
}

type file struct {
	io.ReadSeeker
	info os.FileInfo
	node *Node
	// idx defines current cursor while executing Readdir(n int)
	idx int
}

// Close close implements io.Closer
func (f file) Close() error {
	return nil
}

func (f file) Readdir(n int) ([]os.FileInfo, error) {
	// out
	out := []os.FileInfo{}
	// Find children
	children := f.node.SortedChildren()
	// handle n > 0
	if n > 0 {
		var err error
		// if empty
		if len(children) == 0 {
			return out, io.EOF
		}
		// determine iteration max
		max := f.idx + n
		if max > len(children) {
			max = len(children)
			err = io.EOF
		}
		// output
		for i := f.idx; i < max; i++ {
			out = append(out, children[i].FileInfo())
		}
		f.idx = max - 1
		return out, err
	}
	// dir all children
	for _, sub := range children {
		out = append(out, sub.FileInfo())
	}
	return out, nil
}

func (f file) Stat() (os.FileInfo, error) {
	return f.info, nil
}

// newFile creates a file from a node
func newFile(n *Node) *file {
	return &file{
		ReadSeeker: n.ReadSeeker(),
		info:       n.FileInfo(),
		node:       n,
	}
}
