package main

import (
	"fmt"
	server "http-server/internals"
)

func main() {
	app := server.NewServer(":8080")
	err := app.Listen()
	if err != nil {
		fmt.Println("Connectin Error")
	}
}
