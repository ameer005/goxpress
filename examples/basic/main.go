package main

import (
	"fmt"

	"github.com/ameer005/goxpress"
	"github.com/ameer005/goxpress/httpmethod"
)

var app *goxpress.Server = goxpress.NewServer(":8080")

func main() {
	router := app.Router

	router.Use(func(ctx *goxpress.Context) {
		fmt.Printf("Request: %s %s\n", ctx.Req.RequestMethod())
		ctx.Req.Headers("yo")

	})

	router.Route(httpmethod.GET, "/", func(ctx *goxpress.Context) {
		ctx.Res.Status(200).JSON(map[string]any{"status": "sucess"})
	})

	app.Listen()
}
