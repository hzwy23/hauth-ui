package rdbms

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/utils/hjwt"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/utils"
)

func postModifyPasswd(ctx *context.Context) {

	ctx.Request.ParseForm()

	oriPasswd := ctx.Request.FormValue("orapasswd")
	newPasswd := ctx.Request.FormValue("newpasswd")
	surePasswd := ctx.Request.FormValue("surepasswd")

	if newPasswd != surePasswd {
		logs.Error("new passwd confirm failed. please check your new password and confirm password")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 309, "new passwd confirm failed. please check your new password and confirm password")
		return
	}

	oriEn, err := utils.Encrypt(oriPasswd)
	if err != nil {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "new passwd confirm failed. please check your new password and confirm password")
		return
	}

	newPd, err := utils.Encrypt(newPasswd)
	if err != nil {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "new passwd confirm failed. please check your new password and confirm password")
		return
	}
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	err = dbobj.Exec(sys_rdbms_014, newPd, jclaim.User_id, oriEn)
	if err != nil {
		logs.Error(dbobj.GetErrorMsg(err))
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "new passwd confirm failed. please check your new password and confirm password")
		return
	}
	http.SetCookie(ctx.ResponseWriter, &http.Cookie{Name: "Authorization", Value: "", Path: "/", MaxAge: 3600})
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "modify user password successfully.")
}

func adminModifyPasswd(ctx *context.Context) {
	ctx.Request.ParseForm()

	userid := ctx.Request.FormValue("userid")
	newPasswd := ctx.Request.FormValue("newpasswd")
	surePasswd := ctx.Request.FormValue("surepasswd")

	if newPasswd != surePasswd {
		logs.Error("new passwd confirm failed. please check your new password and confirm password")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 309, "new passwd confirm failed. please check your new password and confirm password")
		return
	}

	newPd, err := utils.Encrypt(newPasswd)
	if err != nil {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "new passwd confirm failed. please check your new password and confirm password")
		return
	}

	err = dbobj.Exec(sys_rdbms_015, newPd, userid)
	if err != nil {
		logs.Error(dbobj.GetErrorMsg(err))
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "new passwd confirm failed. please check your new password and confirm password")
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "modify user password successfully.")
}

func init() {
	beego.Post("/v1/auth/passwd/update", postModifyPasswd)
	beego.Post("/v1/auth/passwd/modify", adminModifyPasswd)
}
