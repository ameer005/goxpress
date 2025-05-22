package server

type HandlerFunc func(ctx *Context)

type Router struct {
	// Method > path > function
	// Get > / > handler
	routes map[string]map[string]HandlerFunc
}

func HandleRequest(ctx *Context, router *Router) {
	method := ctx.req.method
	path := ctx.req.path

	println(method, path)
}
