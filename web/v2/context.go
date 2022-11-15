//go:build v2

package web

import "net/http"

// Context 请求上下文
type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
}
