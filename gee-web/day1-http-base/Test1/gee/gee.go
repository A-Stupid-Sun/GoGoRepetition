package gee

import (
	"fmt"
	"log"
	"net/http"
)

type RouteHandler func(w http.ResponseWriter, req *http.Request)
type Engine struct {
	router map[string]RouteHandler
}

func New() *Engine {
	return &Engine{router: make(map[string]RouteHandler)}
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.Method + "-" + req.URL.Path
	if routeHandler, ok := engine.router[url]; ok {
		routeHandler(w, req)
	} else {
		fmt.Fprintf(w, "404 ERROR FROM URL  %q ", req.URL)
	}
}
func (engine *Engine) addRoute(method string, pattern string, handle RouteHandler) {
	key := method + "-" + pattern
	engine.router[key] = handle
	log.Printf("New Route Added: %q", key)
}
func (engine *Engine) GET(pattern string, handler RouteHandler) {
	engine.addRoute("GET", pattern, handler)
}
func (engine *Engine) POST(pattern string, handler RouteHandler) {
	engine.addRoute("POST", pattern, handler)
}
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
