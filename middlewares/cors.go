package middleware

import (
	"slices"
	"strings"

	"github.com/ameer005/goxpress"
	"github.com/ameer005/goxpress/httpmethod"
)

type CorsOptions struct {
	AllowOrigin    []string
	Credentials    bool
	AllowedMethods []httpmethod.Method
}

func Cors(options CorsOptions) goxpress.HandlerFunc {
	return func(ctx *goxpress.Context) {

		// setting allow origin headers
		origin := ctx.Req.Headers("Origin")
		reqHeaders := ctx.Req.Headers("Access-Control-Request-Headers")

		if len(options.AllowOrigin) == 0 {
			// Allow all origins
			ctx.Res.SetHeader("Access-Control-Allow-Origin", "*")
		} else {
			// If origin header is present and allowed
			if origin != "" && slices.Contains(options.AllowOrigin, origin) {
				ctx.Res.SetHeader("Access-Control-Allow-Origin", origin)
			} else {
				// Origin not allowed or missing, reject
				ctx.Res.Status(403).Send("Forbidden")
				return
			}
		}

		if options.Credentials {
			ctx.Res.SetHeader("Access-Control-Allow-Credentials", "true")
		}

		var allowedMethods strings.Builder
		if len(options.AllowedMethods) > 0 {
			for i, method := range options.AllowedMethods {
				if i > 0 {
					allowedMethods.WriteString(", ")
				}
				allowedMethods.WriteString(string(method))
			}
		} else {
			allowedMethods.WriteString("POST, GET, PUT, PATCH, DELETE, OPTIONS")
		}

		ctx.Res.SetHeader("Access-Control-Allow-Methods", allowedMethods.String())

		if reqHeaders != "" {
			ctx.Res.SetHeader("Access-Control-Allow-Headers", reqHeaders)
		}

		// Handling preflight request
		if ctx.Req.RequestMethod() == string(httpmethod.OPTIONS) {
			ctx.Res.Status(204).Send("")
			return
		}
	}
}
