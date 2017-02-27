package controllers

import (
	"text/template"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/utils/token/hjwt"
	"github.com/hzwy23/hauth/utils/hret"
	"net/http"
	"encoding/json"
	"github.com/hzwy23/hauth/models"
	"github.com/astaxie/beego/logs"
	"github.com/hzwy23/hauth/utils"
	"io/ioutil"
)

type RoleController struct {
	models.RoleModel
}

var RoleCtl = &RoleController{
	models.RoleModel{},
}


func (RoleController) Page(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/hauth/role_info_page.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func (RoleController) ResourcePage(ctx *context.Context){
	defer hret.HttpPanic()
	file,_:=ioutil.ReadFile("./views/hauth/res_role_rel_page.tpl")
	ctx.ResponseWriter.Write(file)
}

func (this RoleController)GetRoleInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	offset := ctx.Request.FormValue("offset")
	limit := ctx.Request.FormValue("limit")
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	rst,err:=this.Get(jclaim.User_id,jclaim.Domain_id,offset,limit)

	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get role info failed.", err)
		return
	}

	hret.WriteJson(ctx.ResponseWriter, rst)
}

func (this RoleController)PostRoleInfo(ctx *context.Context) {

	ctx.Request.ParseForm()

	//取数据
	roleid := ctx.Request.FormValue("role_id")
	rolename := ctx.Request.FormValue("role_name")
	domainid := ctx.Request.FormValue("domain_id")
	rolestatus := ctx.Request.FormValue("role_status")
	id := domainid + "_join_" + roleid
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

	err = this.Post(id, rolename, jclaim.User_id, rolestatus, domainid, roleid)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "add new role info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "add new role info successfully.")
}

func (this RoleController)DeleteRoleInfo(ctx *context.Context) {

	ctx.Request.ParseForm()

	mjson := []byte(ctx.Request.FormValue("JSON"))
	var allrole []models.RoleInfo
	err := json.Unmarshal(mjson, &allrole)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "json解析失败，请重新选择需要删除的角色信息", err)
		return
	}
	err=this.Delete(allrole)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,418,"操作数据库失败。")
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter,"删除角色信息成功。")
}

func (this RoleController)UpdateRoleInfo(ctx *context.Context) {
	ctx.Request.ParseForm()

	Role_id := ctx.Request.FormValue("Role_id")
	Role_name := ctx.Request.FormValue("Role_name")
	Role_status := ctx.Request.FormValue("Role_status")

	err := this.Update(Role_name, Role_status, Role_id)
	if err != nil {
		logs.Error(err.Error())
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "update role info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "update role info successfully.")
}
