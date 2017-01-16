package rdbms

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
)

func userDomainupdate(ctx *context.Context) {

	ctx.Request.ParseForm()
	user_id := ctx.Request.FormValue("userId")
	domain_id := ctx.Request.FormValue("domainId")
	org_id := ctx.Request.FormValue("orgId")
	orgid := domain_id + "_join_" + org_id
	if org_id == "" || domain_id == "" {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 313, "domain_id or org_id is empty.please check values.")
		return
	}

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	if jclaim.User_id == user_id {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "you can't change domain info yourself")
		return
	}

	tx, _ := dbobj.Begin()

	_, err = tx.Exec(sys_rdbms_039, orgid, jclaim.User_id, user_id)
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		ctx.ResponseWriter.Write([]byte("修改用户所属于失败"))
		return
	}
	tx.Commit()
}

func init() {
	beego.Put("/v1/auth/user/domain/update", userDomainupdate)
}
