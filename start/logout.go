package start

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/hret"
)

func init() {
	beego.Any("/logout", func(ctx *context.Context) {
		cookie := http.Cookie{Name: "Authorization", Value: "", Path: "/", MaxAge: 3600}
		http.SetCookie(ctx.ResponseWriter, &cookie)
		hret.WriteHttpOkMsgs(ctx.ResponseWriter, "logout system safely.")
	})
}
