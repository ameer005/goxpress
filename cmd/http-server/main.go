package main

import (
	"fmt"
	server "http-server/internals"
	"http-server/internals/httpmethod"
)

var app *server.Server = server.NewServer(":8080")

func main() {
	router := app.Router
	router.Route(httpmethod.GET, "/", func(ctx *server.Context) {
		ctx.Res.Status(200).Send("yo")
	})

	err := app.Listen()
	if err != nil {
		fmt.Println("Connectin Error")
	}
}
