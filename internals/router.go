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
	method := ctx.Req.RequestMethod()
	path := ctx.Req.RequestPath()

	handler, ok := router.routes[httpmethod.Method(method)][path]

	if !ok {
		fmt.Println("Page not found")
		return
	}

	handler(ctx)
}

func (t *Router) Route(method httpmethod.Method, path string, handler HandlerFunc) {
	if _, exists := t.routes[method]; !exists {
		t.routes[method] = make(map[string]HandlerFunc)
	}
	t.routes[method][path] = handler
}
