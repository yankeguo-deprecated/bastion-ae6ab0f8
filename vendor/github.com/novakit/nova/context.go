package nova // import "github.com/novakit/nova"

import (
	"fmt"
	"net/http"
)

// Context per-request context
type Context struct {
	// Env environment
	Env Env
	// Handlers the handlers
	Handlers []HandlerFunc
	// ErrorHandler must not be nil
	ErrorHandler ErrorHandlerFunc
	// Req the http request
	Req *http.Request
	// Res the http response writer
	Res http.ResponseWriter
	// Values pre-request associated values
	Values map[string]interface{}

	hCursor int // index of current handler invoked
}

// Next invoke the next HandlerFunc registered in application
func (c *Context) Next() {
	// recover from panic
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				c.ErrorHandler(c, err)
			} else {
				c.ErrorHandler(c, fmt.Errorf("panic: %v", r))
			}
		}
	}()

	// reached end of handlers chain
	if c.hCursor >= len(c.Handlers) {
		http.NotFound(c.Res, c.Req)
		return
	}
	// save handler index
	i := c.hCursor
	// increase handler index
	c.hCursor++
	// call handler
	if err := c.Handlers[i](c); err != nil {
		// error handling
		c.ErrorHandler(c, err)
	}
}
