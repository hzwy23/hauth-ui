package main

import (
	"fmt"

	"net/http"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/context"

	_ "github.com/hzwy23/hauth/start"

	_ "github.com/hzwy23/hauth/rdbms"

	"github.com/hzwy23/hauth/token/hjwt"

	"github.com/hzwy23/hauth/route"

	"github.com/hzwy23/hauth/logs"
)

func redictToHtpps() {

	var redirectHandle = http.NewServeMux()

	redirectHandle.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		http.Redirect(w, r, "https://www.asofdate.com", http.StatusMovedPermanently)

	})

	err := http.ListenAndServe(":8081", redirectHandle)

	if err != nil {

		fmt.Println("start http rediect to https failed.")

	}
}

func RequireAuth(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {

	go redictToHtpps()
	route.InsertFilter("/v1/*", beego.BeforeRouter, func(ctx *context.Context) {
		cookie, err := ctx.Request.Cookie("Authorization")
		if err != nil || !hjwt.CheckToken(cookie.Value) {
			logs.Warn("have no authority. redirect to index")
			RequireAuth(ctx.ResponseWriter, ctx.Request)
		}
	})
	route.Run()
}
