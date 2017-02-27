package main

import (
	"net/http"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/context"

	"github.com/hzwy23/hauth/utils/token/hjwt"

	"github.com/hzwy23/hauth/utils/logs"
)

func RequireAuth(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {

	beego.InsertFilter("/v1/*", beego.BeforeRouter, func(ctx *context.Context) {
		cookie, err := ctx.Request.Cookie("Authorization")
		if err != nil || !hjwt.CheckToken(cookie.Value) {
			logs.Warn("have no authority. redirect to index")
			RequireAuth(ctx.ResponseWriter, ctx.Request)
		}
	})
	beego.Run()
}
