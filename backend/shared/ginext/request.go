package ginext

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Request wraps gin.Context with additional functionality
type Request struct {
	GinCtx *gin.Context
	ctx    context.Context
}

// NewRequest creates a new handler request
func NewRequest(c *gin.Context) *Request {
	return &Request{
		GinCtx: c,
		ctx:    c.Request.Context(),
	}
}

// Context returns the request context
func (r *Request) Context() context.Context {
	if r.ctx == nil {
		r.ctx = context.Background()
	}
	return r.ctx
}

// MustBind does a binding on v with incoming request data
// it'll panic if any invalid data (and by design, it should be recovered by error handler middleware)
func (r *Request) MustBind(v interface{}) {
	r.MustNoError(r.GinCtx.ShouldBindJSON(v))
}

// MustBindUri binds URI parameters
func (r *Request) MustBindUri(v interface{}) {
	r.MustNoError(r.GinCtx.ShouldBindUri(v))
}

// MustNoError makes an ASSERT on err variable, panic when it's not nil
// then it must be recovered by WrapHandler
func (r *Request) MustNoError(err error) {
	if err != nil {
		panic(NewError(http.StatusBadRequest, err.Error()))
	}
}

// Param gets URL parameter
func (r *Request) Param(key string) string {
	return r.GinCtx.Param(key)
}

// Query gets query parameter
func (r *Request) Query(key string) string {
	return r.GinCtx.Query(key)
}

// DefaultQuery gets query parameter with default value
func (r *Request) DefaultQuery(key, defaultValue string) string {
	return r.GinCtx.DefaultQuery(key, defaultValue)
}
