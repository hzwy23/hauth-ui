package rdbms

import (
	"html/template"
)

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func getAdvancePage(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/platform/resource/sys_advance_page.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func init() {
	beego.Get("/v1/auth/advance/page", getAdvancePage)
}
