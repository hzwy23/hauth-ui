package start

import (
	"html/template"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

	beego.Get("/", func(ctx *context.Context) {

		huang, _ := template.ParseFiles("./views/login.tpl")

		huang.Execute(ctx.ResponseWriter, nil)
	})

}
