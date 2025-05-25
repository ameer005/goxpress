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

	router.Route(httpmethod.POST, "/", func(ctx *server.Context) {
		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		body, err := server.JSONBody[User](ctx.Req)
		if err != nil {
			ctx.Res.Status(400).Send("Incorrect request data")
			return
		}

		fmt.Println(body.Age)
		fmt.Println(body.Name)

		ctx.Res.Status(200).Send("yo")
	})

	err := app.Listen()
	if err != nil {
		fmt.Println("Connectin Error")
	}
}
