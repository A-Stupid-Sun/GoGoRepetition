package gee

import (
	"fmt"
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
		fmt.Fprint(w, "404 ERROR FROM URL  %q ", req.URL)
	}

}
