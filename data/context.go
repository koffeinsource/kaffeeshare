package data

import (
	"net/http"

	"github.com/koffeinsource/go-klogger"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// Context is our internal context and is passed to every function
type Context struct {
	C   context.Context
	Log klogger.GAELogger
}

// MakeContext creates a default context
func MakeContext(r *http.Request) *Context {
	var ret Context
	ret.C = appengine.NewContext(r)
	ret.Log.Context = ret.C
	return &ret
}
