package router // import "github.com/novakit/router"

import (
	"net/http"

	"github.com/novakit/nova"
)

// Route create a Router on nova application
func Route(n *nova.Nova) Router {
	return Router{Nova: n}
}

// Router router
type Router struct {
	Nova  *nova.Nova
	Rules Rules
}

// Rule make a clone of Router and append a Rule
func (h Router) Rule(r Rule) Router {
	return Router{Nova: h.Nova, Rules: h.Rules.Add(r)}
}

// Method add method restriction
func (h Router) Method(method ...string) Router {
	return h.Rule(MethodRule{Method: method})
}

// Get shortcut for .Method("GET").Path(...)
func (h Router) Get(path string) Router {
	return h.Method(http.MethodGet).Path(path)
}

// Post shortcut for .Method("POST").Path(...)
func (h Router) Post(path string) Router {
	return h.Method(http.MethodPost).Path(path)
}

// Put shortcut for .Method("PUT").Path(...)
func (h Router) Put(path string) Router {
	return h.Method(http.MethodPut).Path(path)
}

// Patch shortcut for .Method("PATCH").Path(...)
func (h Router) Patch(path string) Router {
	return h.Method(http.MethodPatch).Path(path)
}

// Delete shortcut for .Method("DELETE").Path(...)
func (h Router) Delete(path string) Router {
	return h.Method(http.MethodDelete).Path(path)
}

// Path add path restriction
func (h Router) Path(path string) Router {
	return h.Rule(PathRule{Path: path})
}

// Host add host restriction
func (h Router) Host(host ...string) Router {
	return h.Rule(HostRule{Host: host})
}

// Header add header restriction
func (h Router) Header(name string, value ...string) Router {
	return h.Rule(HeaderRule{Name: name, Value: value})
}

// Use complete the handler and register to *arc.Arc, then create a empty Router
func (h Router) Use(handlers ...nova.HandlerFunc) Router {
	for _, handler0 := range handlers {
		handler := handler0
		h.Nova.Use(func(c *nova.Context) (err error) {
			if h.Rules.Match(c) {
				err = handler(c)
			} else {
				c.Next()
			}
			return
		})
	}
	return Route(h.Nova)
}
