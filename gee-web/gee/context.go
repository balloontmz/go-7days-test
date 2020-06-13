//Package gee the context for gee
// 对Web服务来说，无非是根据请求*http.Request，构造响应http.ResponseWriter。但是这两个对象提供的接口粒度太细，比如我们要构造一个完整的响应，
// 需要考虑消息头(Header)和消息体(Body)，而 Header 包含了状态码(StatusCode)，消息类型(ContentType)等几乎每次请求都需要设置的信息。因此，如
// 果不进行有效的封装，那么框架的用户将需要写大量重复，繁杂的代码，而且容易出错。针对常用场景，能够高效地构造出 HTTP 响应是一个好的框架必须考虑的点。
package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//H simple format type
type H map[string]interface{}

//Context the context for the full request
type Context struct {
	//origin
	Writer http.ResponseWriter
	Req    *http.Request
	//request meta
	Path   string
	Method string
	Params map[string]string
	//response info
	StatusCode int
	//middleware
	handlers []HandlerFunc
	index    int
}

//Param get the path param
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

//PostForm get the form value use key
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//Query get the query value use key
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//Status set status for request
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//SetHeader set the header for response
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//there are some function for response below

//String set the string data to response
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//JSON set the JSON data to response
func (c *Context) JSON(code int, data interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

//Data set the Data data to response
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//HTML set the HTML data to response
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

//Next handler request -- for loop for some func that don't have Next inside
func (c *Context) Next() {
	c.index++ // 这一行不能忽略,对每个next 调用,都必须加 1
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

//Fail the request
func (c *Context) Fail(code int, msg string) {
	c.index = len(c.handlers) // 这里会报错
	c.JSON(code, H{"error": msg})
}
