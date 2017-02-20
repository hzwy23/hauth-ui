package main

import (
	"net/http"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/context"

	_ "github.com/hzwy23/hauth/start"

	_ "github.com/hzwy23/hauth/rdbms"

	"github.com/hzwy23/hauth/token/hjwt"

	"github.com/hzwy23/hauth/route"

	"github.com/hzwy23/hauth/logs"
)

func RequireAuth(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {

	route.InsertFilter("/v1/*", beego.BeforeRouter, func(ctx *context.Context) {
		cookie, err := ctx.Request.Cookie("Authorization")
		if err != nil || !hjwt.CheckToken(cookie.Value) {
			logs.Warn("have no authority. redirect to index")
			RequireAuth(ctx.ResponseWriter, ctx.Request)
		}
	})
	route.Run()
}
