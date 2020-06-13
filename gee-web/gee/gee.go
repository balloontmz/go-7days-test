//Package gee the entrance for framework
package gee

import (
	"log"
	"net/http"
	"strings"
)

//HandlerFunc the function that implement handler interface
type HandlerFunc func(*Context)

//Engine the Engine for framework
type Engine struct {
	*RouteGroup // 拥有路由组所有的方法和属性
	router      *router
	groups      []*RouteGroup // 当前应用所有路由组,用于干啥???
}

//RouteGroup the route group for route
type RouteGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouteGroup //为了追溯 prefix,已经通过 预设 prefix 的方式实现
	engine      *Engine
}

//ServeHTTP implement http.Handler interface
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

//New is the constructor for gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouteGroup = &RouteGroup{engine: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
}

//Group is defined to create a RouteGroup
//remember all the group share the engine instance
func (group *RouteGroup) Group(prefix string) *RouteGroup {
	engine := group.engine
	newGroup := &RouteGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//AddRoute add route to framework
func (group *RouteGroup) AddRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//GET defines the method to add GET request
func (group *RouteGroup) GET(pattern string, handler HandlerFunc) {
	group.AddRoute("GET", pattern, handler)
}

//POST defines the method to add POST request
func (group *RouteGroup) POST(pattern string, handler HandlerFunc) {
	group.AddRoute("POST", pattern, handler)
}

//Use middleware for group
func (group *RouteGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
