package goxpress

import (
	"errors"
	"regexp"

	"github.com/ameer005/goxpress/httpmethod"
)

type HandlerFunc func(ctx *Context)

type routeEntry struct {
	path        string
	regex       *regexp.Regexp
	keys        []string
	handler     HandlerFunc
	middlewares []HandlerFunc
}

type router struct {
	// method > [router entry]
	routes            map[httpmethod.Method][]routeEntry
	globalMiddlewares []HandlerFunc
}

func (t *router) Use(middlware HandlerFunc) {
	t.globalMiddlewares = append(t.globalMiddlewares, middlware)
}

// Client side method for defining route and handler function
func (t *router) Route(method httpmethod.Method, path string, handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		return
	}

	// separating main handler function from middlewares
	mainHandler := handlers[len(handlers)-1]
	middlewares := handlers[:len(handlers)-1]

	regex, keys := parsePath(path)

	entry := routeEntry{
		path:        path,
		regex:       regex,
		keys:        keys,
		handler:     mainHandler,
		middlewares: middlewares,
	}

	if _, ok := t.routes[method]; !ok {
		t.routes[method] = []routeEntry{}

	}
	t.routes[method] = append(t.routes[method], entry)
}

// Internal method for handling incoming requests
func HandleRequest(ctx *Context, router *router) {
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

			// running middlewares
			err := router.runMiddlewares(entry.middlewares, ctx)
			if err != nil {
				return
			}

			// rnning route handler
			entry.handler(ctx)
			return
		}

	}

	// if path doesn't exist in memory
	ctx.Res.Status(404).Send("route not found")
	return

}

// Helpers
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

// running all middlewares
// update ctx or return early
func (t *router) runMiddlewares(routeMiddlewares []HandlerFunc, ctx *Context) error {
	// running global middlewares
	for _, middleware := range t.globalMiddlewares {
		middleware(ctx)

		if ctx.Res.ResponseWritten {
			return errors.New("Middleware exit early")
		}
	}

	// running route level middlewares
	for _, middleware := range routeMiddlewares {
		middleware(ctx)

		if ctx.Res.ResponseWritten {
			return errors.New("Middleware exit early")
		}
	}

	return nil
}
