package main

import (
	"fmt"

	"github.com/ameer005/goxpress"
	"github.com/ameer005/goxpress/httpmethod"
)

var app *goxpress.Server = goxpress.NewServer(":8080")

func routMid(ctx *goxpress.Context) {
	ctx.Data["route"] = "message from route level handler"
}

func main() {
	router := app.Router

	app.Use(func(ctx *goxpress.Context) {

		ctx.Data["global"] = "message from global function"
	})

	router.Route(httpmethod.GET, "/", routMid, func(ctx *goxpress.Context) {
		type getQuery struct {
			Page int    `json:"page,string"`
			Q    string `json:"q"`
		}

		q, err := goxpress.QueryData[getQuery](ctx.Req)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(q.Page)
		fmt.Println(q.Q)

		ctx.Res.Status(200).Send("yo")
	})

	err := app.Listen()
	if err != nil {
		fmt.Println("Connectin Error")
	}
}
