package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type RadixNode struct {
	key           string
	children      map[string]*RadixNode
	handler       map[string]*http.HandlerFunc // Changed to pointer to allow nil values
	staticHandler *http.Handler                // Added to store static handlers       // Added to store the HTTP method for the handler
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

// PrintRoutes prints all registered routes with their HTTP methods
func (r *RadixTree) PrintRoutes() {
	fmt.Println("Registered Routes:")
	fmt.Println("==================")

	routes := r.collectRoutes()

	// Sort routes for better readability
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	for _, route := range routes {
		typeInfo := ""
		if route.Type == "static" {
			typeInfo = " (static)"
		}
		fmt.Printf("%-8s %s%s\n", route.Method, route.Path, typeInfo)
	}

	fmt.Printf("\nTotal routes: %d\n", len(routes))
}

// RouteInfo holds information about a single route
type RouteInfo struct {
	Path   string
	Method string
	Type   string // "handler" or "static"
}

// collectRoutes traverses the tree and collects all routes
func (r *RadixTree) collectRoutes() []RouteInfo {
	var routes []RouteInfo
	r.traverseNode(r.root, "", &routes)
	return routes
}

// traverseNode recursively traverses the tree to collect routes
func (r *RadixTree) traverseNode(node *RadixNode, currentPath string, routes *[]RouteInfo) {
	// Build the current path
	path := currentPath
	if node.key != "" {
		if path == "" {
			path = "/" + node.key
		} else {
			path = path + "/" + node.key
		}
	}

	// Handle root path case
	if path == "" {
		path = "/"
	}

	// Check for handlers in the handler map
	if node.handler != nil {
		for method, handler := range node.handler {
			if handler != nil {
				*routes = append(*routes, RouteInfo{
					Path:   path,
					Method: method,
					Type:   "handler",
				})
			}
		}
	}

	// Check for static handler
	if node.staticHandler != nil {
		*routes = append(*routes, RouteInfo{
			Path:   path,
			Method: "GET", // Static handlers are typically GET
			Type:   "static",
		})
	}

	// Traverse children
	for _, child := range node.children {
		r.traverseNode(child, path, routes)
	}
}

// PrintRoutesDetailed prints routes with more detailed information
func (r *RadixTree) PrintRoutesDetailed() {
	fmt.Println("Detailed Route Information:")
	fmt.Println("===========================")

	routes := r.collectRoutes()

	// Group routes by path
	pathGroups := make(map[string][]RouteInfo)
	for _, route := range routes {
		pathGroups[route.Path] = append(pathGroups[route.Path], route)
	}

	// Sort paths
	var paths []string
	for path := range pathGroups {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		routes := pathGroups[path]

		// Sort methods for this path
		sort.Slice(routes, func(i, j int) bool {
			return routes[i].Method < routes[j].Method
		})

		// Print path
		fmt.Printf("\n%s\n", path)
		fmt.Printf("%s\n", strings.Repeat("-", len(path)))

		for _, route := range routes {
			typeInfo := ""
			if route.Type == "static" {
				typeInfo = " (static)"
			}
			fmt.Printf("  %-8s%s\n", route.Method, typeInfo)
		}
	}

	fmt.Printf("\nTotal unique paths: %d\n", len(paths))
	fmt.Printf("Total route handlers: %d\n", len(r.collectRoutes()))
}

// PrintRoutesTree prints routes in a tree-like structure
func (r *RadixTree) PrintRoutesTree() {
	fmt.Println("Route Tree Structure:")
	fmt.Println("=====================")
	r.printNodeTree(r.root, "", true, make(map[string]bool))
}

// printNodeTree prints the tree structure with methods
func (r *RadixTree) printNodeTree(node *RadixNode, prefix string, isLast bool, visited map[string]bool) {
	// Build current path for this node
	currentPath := ""
	if node.key != "" {
		if prefix == "" {
			currentPath = "/" + node.key
		} else {
			currentPath = prefix + "/" + node.key
		}
	} else {
		currentPath = "/"
	}

	// Avoid infinite loops
	if visited[currentPath] {
		return
	}
	visited[currentPath] = true

	// Print node
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	if node.key == "" {
		fmt.Printf("/ (root)\n")
	} else {
		fmt.Printf("%s%s", connector, node.key)

		// Show methods for this node
		var methods []string
		if node.handler != nil {
			for method, handler := range node.handler {
				if handler != nil {
					methods = append(methods, method)
				}
			}
		}
		if node.staticHandler != nil {
			methods = append(methods, "GET (static)")
		}

		if len(methods) > 0 {
			sort.Strings(methods)
			fmt.Printf(" [%s]", strings.Join(methods, ", "))
		}
		fmt.Println()
	}

	// Sort children for consistent output
	var childKeys []string
	for key := range node.children {
		childKeys = append(childKeys, key)
	}
	sort.Strings(childKeys)

	// Print children
	for i, key := range childKeys {
		child := node.children[key]
		isLastChild := i == len(childKeys)-1

		childPrefix := prefix
		if node.key != "" {
			if childPrefix == "" {
				childPrefix = "/" + node.key
			} else {
				childPrefix = childPrefix + "/" + node.key
			}
		}

		// Add proper indentation for child nodes
		nextPrefix := ""
		if !isLast {
			nextPrefix += "│   "
		} else {
			nextPrefix += "    "
		}

		fmt.Print(nextPrefix)
		r.printNodeTree(child, childPrefix, isLastChild, visited)
	}
}

// GetMethodsForPath returns all HTTP methods available for a given path
func (r *RadixTree) GetMethodsForPath(path string) []string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	current := r.root

	// Navigate to the node
	for _, part := range parts {
		if current.children[part] == nil {
			return nil
		}
		current = current.children[part]
	}

	// Collect methods
	var methods []string
	if current.handler != nil {
		for method, handler := range current.handler {
			if handler != nil {
				methods = append(methods, method)
			}
		}
	}
	if current.staticHandler != nil {
		methods = append(methods, "GET (static)")
	}

	sort.Strings(methods)
	return methods
}
