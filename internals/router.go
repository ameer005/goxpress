package server

import (
	"fmt"
	"http-server/internals/httpmethod"
)

type HandlerFunc func(ctx *Context)

type Router struct {
	// Method > path > function
	// Get > / > handler
	routes map[httpmethod.Method]map[string]HandlerFunc
}

func HandleRequest(ctx *Context, router *Router) {
	method := ctx.req.method
	path := ctx.req.path

	handler, ok := router.routes[httpmethod.Method(method)][path]

	if !ok {
		fmt.Println("not found sex sux")
		return
	}

	handler(ctx)

	println(method, path)
}

func (t *Router) Route(method httpmethod.Method, path string, handler HandlerFunc) {
	if _, exists := t.routes[method]; !exists {
		t.routes[method] = make(map[string]HandlerFunc)
	}
	t.routes[method][path] = handler
}
