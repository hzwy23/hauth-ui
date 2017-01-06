package rdbms

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/utils/hjwt"
)

type gettedRoles struct {
	Role_id   string `json:"role_id"`
	Role_name string `json:"role_name"`
}

func getGettedRoles(ctx *context.Context) {
	ctx.Request.ParseForm()
	var rst []gettedRoles
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	rows, err := dbobj.Query(sys_rdbms_046, jclaim.User_id, jclaim.User_id)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query getted roles info failed.", err)
		return
	}
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query getted roles info failed.", err)
		return
	}
	hret.WriteJson(ctx.ResponseWriter, rst)
}

func canGrantRoles(ctx *context.Context) {
	ctx.Request.ParseForm()
	userid := ctx.Request.FormValue("user_id")
	var rst []gettedRoles
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	rows, err := dbobj.Query(sys_rdbms_047, jclaim.User_id, jclaim.User_id, userid)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query getted roles info failed.", err)
		return
	}
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query getted roles info failed.", err)
		return
	}
	hret.WriteJson(ctx.ResponseWriter, rst)
}

func init() {
	beego.Get("/v1/auth/roles/getted", getGettedRoles)
	beego.Get("/v1/auth/roles/canGrant", canGrantRoles)
}
