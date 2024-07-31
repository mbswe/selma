package web

import (
	"net/http"
)

// Router struct holds the registered routes
type Router struct {
	routes map[string]func(http.ResponseWriter, *http.Request)
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
}

// Route adds a new route to the router with a specific handler
func (r *Router) Route(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.routes[path] = handler
}

// ServeHTTP implements the http.Handler interface for Router
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := r.routes[req.URL.Path]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}
