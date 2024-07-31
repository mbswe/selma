// `selma/router.go`
package selma

import (
	"net/http"
)

// Middleware type represents a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Route struct holds the handler and its middleware
type Route struct {
	handler    http.Handler
	middleware []Middleware
}

// Router struct holds the registered routes
type Router struct {
	routes map[string]map[string]*Route
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]*Route),
	}
}

// addRoute adds a new route to the router with a specific method, path, handler, and optional middleware
func (r *Router) addRoute(method, path string, handler http.HandlerFunc, middleware ...Middleware) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]*Route)
	}
	r.routes[path][method] = &Route{handler: handler, middleware: middleware}
}

// Get adds a new GET route to the router
func (r *Router) Get(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.addRoute(http.MethodGet, path, handler, middleware...)
}

// Post adds a new POST route to the router
func (r *Router) Post(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.addRoute(http.MethodPost, path, handler, middleware...)
}

// Put adds a new PUT route to the router
func (r *Router) Put(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.addRoute(http.MethodPut, path, handler, middleware...)
}

// Delete adds a new DELETE route to the router
func (r *Router) Delete(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.addRoute(http.MethodDelete, path, handler, middleware...)
}

// ServeHTTP implements the http.Handler interface for Router
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if routes, ok := r.routes[req.URL.Path]; ok {
		if route, ok := routes[req.Method]; ok {
			handler := route.handler
			for i := len(route.middleware) - 1; i >= 0; i-- {
				handler = route.middleware[i](handler)
			}
			handler.ServeHTTP(w, req)
			return
		}
	}
	http.NotFound(w, req)
}
