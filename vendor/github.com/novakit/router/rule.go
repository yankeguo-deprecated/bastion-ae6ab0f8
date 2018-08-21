package router // import "github.com/novakit/router"

import (
	"net/url"
	"strings"

	"github.com/novakit/nova"
)

// ContextPathParamsKey key for path params in request.Context()
const ContextPathParamsKey = "_path_params"

// PathParams extract path parameters from nova.Context
func PathParams(c *nova.Context) (u url.Values) {
	u, _ = c.Values[ContextPathParamsKey].(url.Values)
	return
}

// Rule is a single rule in a route
type Rule interface {
	// Match match the http.Request, returns if matched
	Match(c *nova.Context) bool
}

// Rules a slice of Rule, with nil-safe methods
type Rules []Rule

// Add create a copy of Rules with new rule added
func (rs Rules) Add(rule Rule) Rules {
	if rs == nil {
		return Rules{rule}
	}
	return Rules{rs, rule}
}

// Match match all rules against request
func (rs Rules) Match(c *nova.Context) bool {
	for _, r := range rs {
		if !r.Match(c) {
			return false
		}
	}
	return true
}

// MethodRule route rule with method restriction
type MethodRule struct {
	Method []string
}

// Match implements Rule
func (r MethodRule) Match(c *nova.Context) bool {
	for _, m := range r.Method {
		if m == c.Req.Method {
			return true
		}
	}
	return false
}

// PathRule route rule with path restriction and params extraction
type PathRule struct {
	Path string
}

// Match implements Rule
func (r PathRule) Match(c *nova.Context) bool {
	pp := url.Values{}
	ns := sanitizePathComponents(strings.Split(r.Path, "/"))
	hs := sanitizePathComponents(strings.Split(c.Req.URL.Path, "/"))
	// length mismatch
	if len(ns) != len(hs) {
		if len(hs) > len(ns) && len(ns) > 0 && strings.HasPrefix(ns[len(ns)-1], "*") {
			// continue if path components longer than pattern components and pattern components has a wildcard ending
		} else {
			return false
		}
	}
	// iterate pattern components
	for i, n := range ns {
		h := hs[i]
		if strings.HasPrefix(n, ":") {
			// capture single parameter
			pp.Set(n[1:], h)
		} else if strings.HasPrefix(n, "*") {
			// capture wildcard parameter
			pp.Set(n[1:], strings.Join(hs[i:], "/"))
			break
		} else {
			// match path component and pattern component
			if n != h {
				return false
			}
		}
	}
	// assign path params
	c.Values[ContextPathParamsKey] = pp
	return true
}

func sanitizePathComponents(in []string) []string {
	ret := make([]string, 0, len(in))
	for _, c := range in {
		if len(c) > 0 {
			ret = append(ret, c)
		}
	}
	return ret
}

// HeaderRule route rule with header restriction
type HeaderRule struct {
	Name  string
	Value []string
}

// Match implements Rule
func (r HeaderRule) Match(c *nova.Context) bool {
	v := c.Req.Header.Get(r.Name)
	for _, vv := range r.Value {
		if vv == v {
			return true
		}
	}
	return false
}

// HostRule route rule with host restriction
type HostRule struct {
	Host []string
}

// Match implements Rule
func (r HostRule) Match(c *nova.Context) bool {
	for _, h := range r.Host {
		if h == c.Req.Host {
			return true
		}
	}
	return false
}
