package main

import (
	"net/http"
	"strings"
)

type RadixNode struct {
	key           string
	children      map[string]*RadixNode
	handler       map[string]*http.HandlerFunc
	staticHandler *http.Handler // Added to store static handlers
	isParameter   bool          // Added to indicate if the node is a parameter node
	args          []string
}

type RadixTree struct {
	root *RadixNode
}

type Route struct {
	handlers map[string]http.HandlerFunc
}

type StaticRoute struct {
	handlers map[string]http.Handler
}

func NewRadixTree() *RadixTree {
	return &RadixTree{
		root: &RadixNode{
			key:      "",
			children: make(map[string]*RadixNode),
			handler:  nil,
		},
	}
}

func (t *RadixTree) Insert(path []string, handler *http.HandlerFunc, staticHandler *http.Handler, method string) {
	current := t.root

	for _, part := range path {
		if current.children[part] == nil {
			current.children[part] = &RadixNode{
				key:           part,
				children:      make(map[string]*RadixNode),
				handler:       make(map[string]*http.HandlerFunc),
				staticHandler: nil,
			}
		}
		current = current.children[part]
	}
	if part := path[len(path)-1]; strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
		current.isParameter = true // Mark as a parameter node
	}
	if handler == nil {
		current.staticHandler = staticHandler
	} else {
		current.handler[method] = handler
	}
}

func (t *RadixTree) Find(path string, method string) *RadixNode {
	current := t.root
	parts := strings.Split(strings.Trim(path, "/"), "/")

	for _, part := range parts {

		if current.children[part] == nil {
			for _, childNode := range current.children {
				if childNode.isParameter { // Check for parameter nodes
					current = childNode
					current.args = append(current.args, part) // Store the parameter value
					return current
				}
			}
			return nil // No se encontró el nodo
		}
		current = current.children[part]
	}

	if current.handler != nil {
		if current.handler[method] != nil {
			return current
		}
	}

	if current.staticHandler != nil {
		return current
	}
	return nil
}

func NewRouter() *RadixTree {
	return NewRadixTree()
}

func (r *RadixTree) Handle(method, path string, handler http.HandlerFunc) {
	path_split := strings.Split(strings.Trim(path, "/"), "/")
	r.Insert(path_split, &handler, nil, method)
}

func (r *RadixTree) HandleStatic(path string, handler http.Handler) {
	path_split := strings.Split(strings.Trim(path, "/"), "/")
	r.Insert(path_split, nil, &handler, http.MethodGet)
}

func (r *RadixTree) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	node := r.root
	if strings.Contains(req.URL.Path, "/static/") {
		node = r.Find("static", http.MethodGet)
	} else {
		node = r.Find(req.URL.Path, req.Method)
	}
	if node == nil {
		http.NotFound(w, req)
		return
	}
	handler, exists := node.handler[req.Method]
	if exists && handler != nil {
		(*handler)(w, req)
		return
	}

	if node.staticHandler != nil {
		(*node.staticHandler).ServeHTTP(w, req)
		return
	}
}

func (r *RadixTree) GET(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *RadixTree) POST(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *RadixTree) PUT(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}

func (r *RadixTree) DELETE(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}
