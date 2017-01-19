package start

import (
	"html/template"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/route"
)

func init() {

	route.Get("/", func(ctx *context.Context) {

		huang, _ := template.ParseFiles("./views/login.tpl")

		huang.Execute(ctx.ResponseWriter, nil)
	})

}
