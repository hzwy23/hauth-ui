package rdbms

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
)

type SysResThemeInfo struct {
	Theme_id     string `json:"theme_id"`
	Theme_desc   string `json:"theme_name"`
	Res_id       string `json:"res_id"`
	Res_url      string `json:"res_url"`
	Res_type     string `json:"res_type"`
	Res_bg_color string `json:"res_bg_color"`
	Res_class    string `json:"res_class"`
	Group_id     string `json:"group_id"`
	Res_icon     string `json:"res_img"`
	Sort_id      string `json:"sort_id"`
}

func GetThemesInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	var res_id = ctx.Request.FormValue("res_id")
	var rst []SysResThemeInfo
	rows, err := dbobj.Query(sys_rdbms_070, jclaim.User_id, res_id)

	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "query database failed.")
		return
	}
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "query database failed.")
		return
	}
	hret.WriteJson(ctx.ResponseWriter, rst)
}

func init() {
	beego.Get("/v1/auth/user/theme/resource", GetThemesInfo)
}
