package main

import (
	"fmt"

	"github.com/ameer005/goxpress"
	"github.com/ameer005/goxpress/httpmethod"
)

var app *goxpress.Server = goxpress.NewServer(":8080")

func main() {
	router := app.Router

	router.Route(httpmethod.GET, "/", func(ctx *goxpress.Context) {
		ctx.Res.Status(200).JSON(map[string]any{"status": "sucess"})
	})

	router.Route(httpmethod.POST, "/", func(ctx *goxpress.Context) {

		metadata, _ := ctx.Req.FormFile("file")
		fmt.Println("yo")
		fmt.Println(metadata.Path)

		ctx.Res.Status(200).JSON(map[string]any{"status": "sucess"})
	})

	app.Listen()
}
