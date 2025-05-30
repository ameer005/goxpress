package main

import (
	"fmt"

	"github.com/ameer005/goxpress"
	"github.com/ameer005/goxpress/httpmethod"
	middleware "github.com/ameer005/goxpress/middlewares"
)

var app *goxpress.Server = goxpress.NewServer(":8080")

func routMid(ctx *goxpress.Context) {
	ctx.Data["route"] = "message from route level handler"
}

func main() {
	router := app.Router
	app.Use(middleware.Cors(middleware.CorsOptions{
		AllowOrigin: []string{"localhost:8080"},
	}))

	app.Use(func(ctx *goxpress.Context) {

		ctx.Data["global"] = "message from global function"
	})

	router.Route(httpmethod.GET, "/", routMid, func(ctx *goxpress.Context) {
		ctx.Res.Status(200).Send("yo")
	})

	err := app.Listen()
	if err != nil {
		fmt.Println("Connectin Error")
	}
}
