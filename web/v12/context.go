package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Context 请求上下文
type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter

	// 主要是为了给 middleware 读写使用
	RespData       []byte
	RespStatusCode int

	PathParams map[string]string

	// 缓存的数据
	queryValues url.Values

	// 命中的路由
	MatchRoute string
}

func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) RespJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	//c.Resp.WriteHeader(status)
	c.RespData = data
	c.RespStatusCode = status
	//c.Resp.Header().Set("Content-Type", "application/json")
	//c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	//n, err := c.Resp.Write(data)
	//if n != len(data) {
	//	return errors.New("web: 未写入全部数据")
	//}
	return err
}

func (c *Context) BindJSON(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body 为 nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: c.Req.FormValue(key)}
}

func (c *Context) QueryValue(key string) StringValue {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}
	values, ok := c.queryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: values[0]}
}

func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: val}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
