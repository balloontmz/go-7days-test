//Package gee the router for framework
// 我们将和路由相关的方法和结构提取了出来，放到了一个新的文件中router.go，方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
// router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context。
package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTextRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func Test_newRouter(t *testing.T) {
	tests := []struct {
		name string
		want *router
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{args: args{"/p/:name"}, want: []string{"p", ":name"}},
		{args: args{"/p/*"}, want: []string{"p", "*"}},
		{args: args{"/p/*name"}, want: []string{"p", "*name"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePattern(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_router_getRoute(t *testing.T) {

	type args struct {
		method string
		path   string
	}
	tests := []struct {
		name  string
		args  args
		want  *node
		want1 map[string]string
	}{
		{args: args{method: "GET", path: "/hello/geekkutu"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTextRouter()
			n, ps := r.getRoute(tt.args.method, tt.args.path)
			if n == nil {
				t.Fatal("nil shouldn't be returned")
			}
			if n.pattern != "/hello/:name" {
				t.Fatal("should match /hello/:name")
			}
			if ps["name"] != "geekkutu" {
				t.Fatal("name should be equal to 'geekkutu'")
			}
			fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
		})
	}
}

func Test_router_addRoute(t *testing.T) {
	type fields struct {
		roots    map[string]*node
		handlers map[string]HandlerFunc
	}
	type args struct {
		method  string
		pattern string
		handler HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{
				roots:    tt.fields.roots,
				handlers: tt.fields.handlers,
			}
			r.addRoute(tt.args.method, tt.args.pattern, tt.args.handler)
		})
	}
}

func Test_router_handle(t *testing.T) {
	type fields struct {
		roots    map[string]*node
		handlers map[string]HandlerFunc
	}
	type args struct {
		c *Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{
				roots:    tt.fields.roots,
				handlers: tt.fields.handlers,
			}
			r.handle(tt.args.c)
		})
	}
}
