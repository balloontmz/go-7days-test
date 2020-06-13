//Package gee the router for framework
// 我们将和路由相关的方法和结构提取了出来，放到了一个新的文件中router.go，方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
// router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context。
package gee

import (
	"net/http"
	"strings"
)

type router struct {
	// roots key eg, roots['GET'] roots['POST']
	// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//only one * allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//getRoute get pattern(node) and params use request method and path
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // the params after * may not being parsed

	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// fmt.Print("before search", searchParts, "\n")
	n := root.search(searchParts, 0)
	if n != nil { // 匹配到节点的情况
		//还需要赋值参数
		//匹配成功 n.pattern 必定存在
		parts := parsePattern(n.pattern)
		// fmt.Print("当前参数为", parts)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

//addRoute add route for pattern
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

//handle handle the request context
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusOK, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
