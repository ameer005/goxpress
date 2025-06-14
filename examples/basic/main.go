package main

import (
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
		/* type formdata struct { */
		/* 	Username string `json:"username"` */
		/* 	Email    string `json:"email"` */
		/* } */
		/**/
		/* data, err := goxpress.FormBody[formdata](ctx.Req) */
		/* if err != nil { */
		/* 	fmt.Println(err) */
		/* } */
		/**/

		ctx.Res.Status(200).JSON(map[string]any{"status": "sucess"})
	})

	app.Listen()
}
