package main

import (
	"net/http"
	"strings"
)

type Route struct {
	handlers map[string]http.HandlerFunc
}

type StaticRoute struct {
	handlers map[string]http.Handler
}

type Router struct {
	routes  map[string]*Route
	statics map[string]*StaticRoute
}

func NewRouter() *Router {
	return &Router{
		routes:  make(map[string]*Route),
		statics: make(map[string]*StaticRoute),
	}
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = &Route{
			handlers: make(map[string]http.HandlerFunc),
		}
	}
	r.routes[path].handlers[method] = handler
}

func (r *Router) HandleStatic(path string, handler http.Handler) {
	if r.statics[path] == nil {
		r.statics[path] = &StaticRoute{
			handlers: make(map[string]http.Handler),
		}
	}
	r.statics[path].handlers[http.MethodGet] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if strings.HasPrefix(req.URL.Path, "/static/") {
		handler := r.statics["/static/"].handlers[http.MethodGet]
		handler.ServeHTTP(w, req)
		return
	} else {
		if route, exists := r.routes[req.URL.Path]; exists {
			if handler, methodExists := route.handlers[req.Method]; methodExists {
				handler(w, req)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		} else {
			for path, route := range r.routes {
				if strings.HasSuffix(path, "/") && strings.HasPrefix(req.URL.Path, path) {
					if handler, methodExists := route.handlers[req.Method]; methodExists {
						handler(w, req)
						return
					}
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					return
				}
			}
		}
	}

	http.NotFound(w, req)
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}
