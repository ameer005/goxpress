package server

import (
	"http-server/internals/httpmethod"
	"regexp"
)

type HandlerFunc func(ctx *Context)

type RouteEntry struct {
	path    string
	regex   *regexp.Regexp
	keys    []string
	hanlder HandlerFunc
}

type Router struct {
	// method > [router entry]
	routes map[httpmethod.Method][]RouteEntry
}

// Client side method for defining route and handler function
func (t *Router) Route(method httpmethod.Method, path string, handler HandlerFunc) {
	regex, keys := parsePath(path)

	entry := RouteEntry{
		path:    path,
		regex:   regex,
		keys:    keys,
		hanlder: handler,
	}

	if _, ok := t.routes[method]; !ok {
		t.routes[method] = []RouteEntry{}

	}
	t.routes[method] = append(t.routes[method], entry)
}

// Internal method for handling incoming requests
func HandleRequest(ctx *Context, router *Router) {
	method := ctx.Req.RequestMethod()
	path := ctx.Req.RequestPath()

	// route Entry
	entries, ok := router.routes[httpmethod.Method(method)]

	if !ok {
		ctx.Res.Status(404).Send("route not found")
		return
	}

	for _, entry := range entries {

		// matching cur path with string
		// match = [full string, extracted values...]
		matches := entry.regex.FindStringSubmatch(path)

		// finding the path in routes
		if matches != nil {
			// building params
			for i, key := range entry.keys {
				ctx.Req.SetRequestParam(key, matches[i+1])
			}

			// funning route handler
			entry.hanlder(ctx)
			return
		}

	}

	// if path doesn't exist in memory
	ctx.Res.Status(404).Send("route not found")
	return

}

func parsePath(path string) (*regexp.Regexp, []string) {
	// for storing params key like id
	var keys []string

	// regex for replacing : with regex patter
	re := regexp.MustCompile(`:([^/]+)`)

	regexPattern := re.ReplaceAllStringFunc(path, func(s string) string {
		key := s[1:]
		keys = append(keys, key)
		return "([^/]+)"
	})
	regexPattern = "^" + regexPattern + "$"

	return regexp.MustCompile(regexPattern), keys
}
