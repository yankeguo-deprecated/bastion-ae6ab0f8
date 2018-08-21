package binfs // import "github.com/novakit/binfs"

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// dummy io.ReadSeeker for directory
type dirReadSeeker struct{}

func (dirReadSeeker) Read(p []byte) (n int, err error) {
	return 0, ErrIsDirectory
}

func (dirReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, ErrIsDirectory
}

// fileInfo implements os.FileInfo
type fileInfo struct {
	name  string
	size  int64
	date  time.Time
	isDir bool
}

func (f fileInfo) Name() string {
	return f.name
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() os.FileMode {
	if f.isDir {
		return os.FileMode(0777)
	}
	return os.FileMode(0666)
}

func (f fileInfo) ModTime() time.Time {
	return f.date
}

func (f fileInfo) IsDir() bool {
	return f.isDir
}

func (fileInfo) Sys() interface{} {
	return nil
}

// Chunk a file in a binfs
type Chunk struct {
	Path []string
	Date time.Time
	Data []byte
}

// Node represents a internal node in file tree
type Node struct {
	Path     []string
	Name     string
	Children map[string]*Node
	Chunk    *Chunk
}

// NodeWalker function to walk over all nodes
type NodeWalker func(n *Node)

// Walk walk over all nodes
func (n *Node) Walk(fn NodeWalker) {
	if n != nil {
		fn(n)
		if n.Children != nil {
			for _, v := range n.Children {
				v.Walk(fn)
			}
		}
	}
}

// Load load a file into zone
func (n *Node) Load(c *Chunk) {
	n.Ensure(c.Path...).Chunk = c
}

// Open open a file, a partial mocking of *os.File
func (n *Node) Open(name string) (File, error) {
	comps := strings.Split(name, "/")
	c := n.Find(comps...)
	if c == nil {
		return nil, os.ErrNotExist
	}
	return newFile(c), nil
}

// Child Find or create a child
func (n *Node) Child(name string) *Node {
	if n.Path == nil {
		n.Path = []string{}
	}
	if n.Children == nil {
		n.Children = map[string]*Node{}
	}
	c := n.Children[name]
	if c == nil {
		nPath := make([]string, len(n.Path))
		copy(nPath, n.Path)
		nPath = append(nPath, name)
		c = &Node{
			Name: name,
			Path: nPath,
		}
		n.Children[name] = c
	}
	return c
}

// Ensure find or create a deep child
func (n *Node) Ensure(name ...string) *Node {
	t := n
	for _, v := range name {
		if v == "" {
			continue
		}
		t = t.Child(v)
	}
	return t
}

// Find a deep child
func (n *Node) Find(name ...string) *Node {
	t := n
	for _, v := range name {
		// ignoring invalid name components
		if len(v) == 0 || v == "." || v == ".." {
			continue
		}
		if t != nil && t.Children != nil {
			t = t.Children[v]
		} else {
			return nil
		}
	}
	return t
}

// SortedChildren returns children sorted by name
func (n *Node) SortedChildren() []*Node {
	if n == nil || n.Children == nil {
		return []*Node{}
	}
	out := []*Node{}
	keys := []string{}
	for k := range n.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		out = append(out, n.Children[k])
	}
	return out
}

// FileInfo creates a related os.FileInfo
func (n *Node) FileInfo() os.FileInfo {
	info := fileInfo{}
	if len(n.Path) > 0 {
		info.name = n.Path[len(n.Path)-1]
	}
	if n.Chunk != nil {
		info.date = n.Chunk.Date
		info.size = int64(len(n.Chunk.Data))
	} else {
		info.isDir = true
	}
	return info
}

// ReadSeeker creates a related io.ReadSeeker
func (n *Node) ReadSeeker() io.ReadSeeker {
	if n.Chunk != nil {
		return bytes.NewReader(n.Chunk.Data)
	}
	return dirReadSeeker{}
}

// nodeWrapper wraps Node to http.FileSystem
type nodeWrapper struct {
	n *Node
}

func (n nodeWrapper) Open(file string) (http.File, error) {
	return n.n.Open(file)
}

// FileSystem creates http.FileSystem implementation
func (n *Node) FileSystem() http.FileSystem {
	return nodeWrapper{n: n}
}
