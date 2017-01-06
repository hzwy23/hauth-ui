package rdbms

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/utils/hjwt"

	"github.com/astaxie/beego"
	"github.com/hzwy23/dbobj"
)

type defaultMenu struct {
	Res_id       string
	Res_name     string
	Res_bg_color string
	Res_class    string
	Res_url      string
	Res_img      string
	Group_id     string
}

func getDefaultMenuInfo(ctx *context.Context) {
	defer hret.HttpPanic()
	typeId := ctx.Request.FormValue("TypeId")
	Id := ctx.Request.FormValue("Id")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	var one defaultMenu
	var rst []defaultMenu

	rows, err := dbobj.Query(sys_rdbms_012, Id, typeId, jclaim.User_id)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "please contact sysadmin.query main menu failed.")
		return
	}
	for rows.Next() {
		err := rows.Scan(&one.Res_id,
			&one.Res_name,
			&one.Res_url,
			&one.Res_bg_color,
			&one.Res_class,
			&one.Res_img,
			&one.Group_id)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "please contact sysadmin.query main menu failed.")
			return
		}
		rst = append(rst, one)
	}

	ojs, err := json.Marshal(rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "please contact sysadmin.query main menu failed.")
		return
	}
	ctx.ResponseWriter.Write(ojs)
}

func init() {
	beego.Get("/v1/auth/main/menu", getDefaultMenuInfo)
}
