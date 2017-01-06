package rdbms

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/utils"
	"github.com/hzwy23/hauth/utils/hjwt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/hret"
)

type RoleInfo struct {
	Role_id             string
	Role_name           string
	Role_owner          string
	Role_create_date    string
	Role_status_desc    string
	Role_status         string
	Domain_id           string
	Domain_desc         string
	Role_maintance_date string
	Role_maintance_user string
}

func getRoleInfoPage(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/platform/resource/role_info_page.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func getRoleInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	rows, err := dbobj.Query(sys_rdbms_028, jclaim.User_id, jclaim.Domain_id, offset, limit)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
	}

	var rst []RoleInfo
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get role info failed.", err)
		return
	}

	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count("select count(*) from sys_role_info"), rst)
}

func postRoleInfo(ctx *context.Context) {

	ctx.Request.ParseForm()

	//取数据
	roleid := ctx.Request.FormValue("role_id")
	rolename := ctx.Request.FormValue("role_name")
	domainid := ctx.Request.FormValue("domain_id")
	rolestatus := ctx.Request.FormValue("role_status")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	//校验
	if !utils.ValidWord(roleid, 1, 30) {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "please input role id number.")
		return
	}
	//
	if !utils.ValidHanAndWord(rolename, 1, 30) {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "角色名称必须是汉字,字母,或者下划线的组合,并且长度不能小于30")
		return
	}

	err = dbobj.Exec(sys_rdbms_026, roleid, rolename, jclaim.User_id, rolestatus, domainid, jclaim.User_id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "add new role info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "add new role info successfully.")
}

func deleteRoleInfo(ctx *context.Context) {

	ctx.Request.ParseForm()

	mjson := []byte(ctx.Request.FormValue("JSON"))
	var allrole []RoleInfo
	err := json.Unmarshal(mjson, &allrole)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "json解析失败，请重新选择需要删除的角色信息", err)
		return
	}

	for _, val := range allrole {
		err := dbobj.Exec(sys_rdbms_027, val.Role_id)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "delete role info failed.", err)
			return
		}
		logs.Info("delete role info successfully. role id is :", val.Role_id)
	}
}

func updateRoleInfo(ctx *context.Context) {
	ctx.Request.ParseForm()

	Role_id := ctx.Request.FormValue("Role_id")
	Role_name := ctx.Request.FormValue("Role_name")
	Role_status := ctx.Request.FormValue("Role_status")

	err := dbobj.Exec(sys_rdbms_050, Role_name, Role_status, Role_id)
	if err != nil {
		logs.Error(err.Error())
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "update role info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "update role info successfully.")
}

func init() {
	beego.Get("/v1/auth/role/page", getRoleInfoPage)
	beego.Get("/v1/auth/role/get", getRoleInfo)
	beego.Post("/v1/auth/role/post", postRoleInfo)
	beego.Post("/v1/auth/role/delete", deleteRoleInfo)
	beego.Put("/v1/auth/role/update", updateRoleInfo)
}
