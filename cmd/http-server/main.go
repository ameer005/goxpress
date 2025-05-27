package main

import (
	"fmt"
	server "http-server/internals"
	"http-server/internals/httpmethod"
)

var app *server.Server = server.NewServer(":8080")

func routMid(ctx *server.Context) {
	ctx.Data["route"] = "message from route level handler"
}

func main() {
	router := app.Router

	app.Use(func(ctx *server.Context) {

		ctx.Data["global"] = "message from global function"
	})

	router.Route(httpmethod.GET, "/", routMid, func(ctx *server.Context) {
		fmt.Println(ctx.Data["global"])
		fmt.Println(ctx.Data["route"])

		ctx.Res.Status(200).Send("yo")
	})

	err := app.Listen()
	if err != nil {
		fmt.Println("Connectin Error")
	}
}
