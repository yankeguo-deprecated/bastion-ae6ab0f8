package binfs // import "github.com/novakit/binfs"

import (
	"net/http"
)

// DefaultRoot the default root of BinFS
var DefaultRoot = &Node{}

// Load load a file into zone
func Load(c *Chunk) {
	DefaultRoot.Load(c)
}

// Open open a file, a partial mocking of *os.File
func Open(name string) (File, error) {
	return DefaultRoot.Open(name)
}

// Find find a deep child node
func Find(name ...string) *Node {
	return DefaultRoot.Find(name...)
}

// Walk walk the default root
func Walk(fn NodeWalker) {
	DefaultRoot.Walk(fn)
}

// FileSystem creates http.FileSystem implementation
func FileSystem() http.FileSystem {
	return DefaultRoot.FileSystem()
}
