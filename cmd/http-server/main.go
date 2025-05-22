package main

import (
	"fmt"
	server "http-server/internals"
)

var app *server.Server = server.NewServer(":8080")

func main() {

	err := app.Listen()
	if err != nil {
		fmt.Println("Connectin Error")
	}
}
