package rdbms

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/hzwy23/dbobj"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
)

type RoleUserRel struct {
	Uuid           string
	User_id        string
	Role_id        string
	Role_name      string
	Maintance_date string
	Maintance_user string
}

func getRoleUserPage(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/platform/resource/role_user_rel_page.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func getRoleUserRel(ctx *context.Context) {

	ctx.Request.ParseForm()

	userId := ctx.Request.FormValue("UserId")

	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))

	rows, err := dbobj.Query(sys_rdbms_022, userId, offset, limit)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get role user relation info failed.", err)
		return
	}
	var one RoleUserRel
	var rst []RoleUserRel
	for rows.Next() {
		err := rows.Scan(&one.Uuid,
			&one.User_id,
			&one.Role_id,
			&one.Role_name,
			&one.Maintance_date,
			&one.Maintance_user)
		if err != nil {
			logs.Error(err)
			return
		}
		rst = append(rst, one)
	}
	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_023, userId), rst)
}

func postRoleUserRel(ctx *context.Context) {

	ctx.Request.ParseForm()

	ijs := []byte(ctx.Request.FormValue("JSON"))
	var rst []RoleUserRel
	err := json.Unmarshal(ijs, &rst)
	if err != nil {
		logs.Error(err.Error())
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 319, "json data unmarshal failed.", err)
	}

	maintanceDate := time.Now().Format("2006-01-02")
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 318, "No Auth", err)
		return
	}

	for _, val := range rst {
		err := dbobj.Exec(sys_rdbms_024, val.Role_id, val.User_id, maintanceDate, jclaim.User_id)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "add role and user relation faild.", err)
		}
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "add role and user relation successfully.")
}

func deleteRoleUserRel(ctx *context.Context) {
	ctx.Request.ParseForm()
	mjson := []byte(ctx.Request.FormValue("JSON"))
	var allrole []RoleUserRel
	err := json.Unmarshal(mjson, &allrole)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "json data unmarshal failed.", err)
		return
	}

	for _, val := range allrole {
		err := dbobj.Exec(sys_rdbms_025, val.Uuid)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "delete role and user relation failed.", err)
			return
		}
		logs.Info("delete role of user,role code uuid is :", val.Uuid)
	}
}

func init() {
	beego.Get("/v1/auth/roleuser/page", getRoleUserPage)
	beego.Get("/v1/auth/roleuser/get", getRoleUserRel)
	beego.Post("/v1/auth/roleuser/post", postRoleUserRel)
	beego.Post("/v1/auth/roleuser/delete", deleteRoleUserRel)
}
