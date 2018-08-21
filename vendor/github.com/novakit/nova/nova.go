package nova // import "github.com/novakit/nova"

import (
	"net/http"
)

// HandlerFunc function handles a per-request context and returns error
type HandlerFunc func(*Context) error

// ErrorHandlerFunc functions handles a pre-request context and a error
type ErrorHandlerFunc func(*Context, error)

// Nova application
type Nova struct {
	// Env environment must be one of "production", "development", "test"
	Env Env
	// Handlers slice of handler functions
	Handlers []HandlerFunc
	// ErrorHandler handler will be invoked on error returned by previous handlers
	ErrorHandler ErrorHandlerFunc
}

// New create a new instance of Nova
func New() *Nova {
	return &Nova{
		Handlers: []HandlerFunc{},
		ErrorHandler: func(c *Context, err error) {
			if c.Env.IsDevelopment() || c.Env.IsTest() {
				http.Error(c.Res, "internal server error: "+err.Error(), http.StatusInternalServerError)
			} else {
				http.Error(c.Res, "internal server error", http.StatusInternalServerError)
			}
			return
		},
	}
}

// Use add mutiple handler functions
func (n *Nova) Use(handlers ...HandlerFunc) {
	n.Handlers = append(n.Handlers, handlers...)
}

// ErrorHandler register a error handler
func (n *Nova) Error(handler ErrorHandlerFunc) {
	n.ErrorHandler = handler
}

// CreateContext create a nova.context from http.Request / http.ResponseWriter
func (n *Nova) CreateContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Env:          n.Env,
		Handlers:     n.Handlers,
		ErrorHandler: n.ErrorHandler,
		Req:          req,
		Res:          res,
		Values:       map[string]interface{}{},
	}
}

// ServeHTTP implements http.HandlerFunc
func (n *Nova) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	n.CreateContext(res, req).Next()
}
