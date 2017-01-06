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

func userDomainupdate(ctx *context.Context) {

	ctx.Request.ParseForm()
	user_id := ctx.Request.FormValue("userId")
	domain_id := ctx.Request.FormValue("domainId")
	org_id := ctx.Request.FormValue("orgId")

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
	sql := `update sys_user_domain_rel set domain_id = ?,maintance_date = NOW(),grant_user_id = ? where user_id = ?`

	_, err = tx.Exec(sql, domain_id, jclaim.User_id, user_id)
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		ctx.ResponseWriter.Write([]byte("修改用户所属于失败"))
		return
	}

	sql = `update sys_user_info set org_unit_id = ?,user_maintance_date = now(),user_maintance_user = ? where user_id = ?`

	_, err = tx.Exec(sql, org_id, jclaim.User_id, user_id)
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
