package static // import "github.com/novakit/static"

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/novakit/binfs"
	"github.com/novakit/nova"
)

// Options options for static
type Options struct {
	// Prefix path prefix of url
	Prefix string
	// Directory directory to serve, will use WEBROOT envvar
	Directory string
	// BinFS using binfs, not suggested with WEBROOT envvar
	BinFS bool
	// Index serving index.html
	Index bool
}

func sanitizeOptions(opts ...Options) (opt Options) {
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = os.Getenv("WEBROOT")
	}
	if len(opt.Directory) == 0 {
		opt.Directory = "public"
	}
	return
}

func buildFileSystem(opt Options) (fs http.FileSystem) {
	if opt.BinFS {
		c := strings.Split(opt.Directory, "/")
		n := binfs.Find(c...)
		if n == nil {
			panic("directory not find in binfs")
		}
		fs = n.FileSystem()
	} else {
		fs = http.Dir(opt.Directory)
	}
	return
}

func trimPathPrefix(pfx string, path string) (ret string, ok bool) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(pfx, "/") {
		pfx = "/" + pfx
	}
	if strings.HasPrefix(path, pfx) {
		ret = "/" + path[len(pfx):]
		ok = true
		return
	}
	ret = path
	return
}

func setContentType(stat os.FileInfo, res http.ResponseWriter) {
	mt := mime.TypeByExtension(path.Ext(stat.Name()))
	if len(mt) > 0 {
		res.Header().Set("Content-Type", mt)
	}
}

func setContentLength(stat os.FileInfo, res http.ResponseWriter) {
	res.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
}

func setLastModified(stat os.FileInfo, res http.ResponseWriter) {
	res.Header().Set("Last-Modified", stat.ModTime().Format(http.TimeFormat))
}

// Handler create a nova.HandlerFunc
func Handler(opts ...Options) nova.HandlerFunc {
	opt := sanitizeOptions(opts...)
	fs := buildFileSystem(opt)
	return func(c *nova.Context) (err error) {
		// must be GET/HEAD method
		if c.Req.Method != http.MethodGet && c.Req.Method != http.MethodHead {
			c.Next()
			return
		}

		// index fixed flag
		var indexRewriteTried bool

		// trim and validate prefix path
		fsPath := path.Clean(c.Req.URL.Path)
		if len(opt.Prefix) > 0 {
			var ok bool
			if fsPath, ok = trimPathPrefix(opt.Prefix, fsPath); !ok {
				c.Next()
				return
			}
		}

		// open file
	OPEN:
		var file http.File
		if file, err = fs.Open(fsPath); err != nil {
			// bypass 404, 403
			if os.IsNotExist(err) || os.IsPermission(err) {
				err = nil
				c.Next()
			}
			return
		}
		defer file.Close()

		// stat file
		var stat os.FileInfo
		if stat, err = file.Stat(); err != nil {
			return
		}

		// skip dir
		if stat.IsDir() {
			if opt.Index && !indexRewriteTried {
				// set index fixed mark to true, prevent loop
				indexRewriteTried = true
				// close the dir
				file.Close()
				// change fs path
				fsPath = path.Join(fsPath, "index.html")
				// goto open
				goto OPEN
			} else {
				c.Next()
				return
			}
		}

		// set content-type/content-length/last-modified
		setContentType(stat, c.Res)
		setContentLength(stat, c.Res)
		setLastModified(stat, c.Res)

		// check if-modified-since
		nma := c.Req.Header.Get("If-Modified-Since")
		if len(nma) > 0 {
			var t time.Time
			if t, err = http.ParseTime(nma); err != nil {
				// ignore error, continue sending file
				err = nil
			} else {
				if !stat.ModTime().After(t) {
					c.Res.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}

		// send 200
		c.Res.WriteHeader(http.StatusOK)

		// send body if GET
		if c.Req.Method == http.MethodGet {
			_, err = io.Copy(c.Res, file)
		}
		return
	}
}
