package rdbms

import (
	"io/ioutil"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
)

func getEntry(ctx *context.Context) {
	defer hret.HttpPanic()
	ctx.Request.ParseForm()
	id := ctx.Request.FormValue("Id")
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	row := dbobj.QueryRow(sys_rdbms_011, jclaim.User_id, id)

	var url string
	err = row.Scan(&url)
	if err != nil {
		logs.Error("cant not fetch menu_url", err)
		url = "./views/theme/default/sysconfig.tpl"
	}

	hz, _ := ioutil.ReadFile(url)
	ctx.ResponseWriter.Write(hz)
}

func init() {
	beego.Get("/v1/auth/index/entry", getEntry)
}
