package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/models"
	"github.com/hzwy23/hauth/utils/hret"
	"github.com/hzwy23/hauth/utils/logs"
	"github.com/hzwy23/hauth/utils/token/hjwt"
)

type userRolesController struct {
	models *models.UserRolesModel
}

var UserRolesController = &userRolesController{
	models: new(models.UserRolesModel),
}

func (this userRolesController) CleanUserRoles(ctx *context.Context) {
	ctx.Request.ParseForm()
	var rst []models.UserRolesModel
	err := json.Unmarshal([]byte(ctx.Request.FormValue("JSON")), &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "unmarshal failed.", err)
		return
	}
	err = this.models.CleanRoles(rst)
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "clean roles of user successfully.")
}


func (this userRolesController)GetRolesByUuser(ctx *context.Context){
	ctx.Request.ParseForm()
	user_id := ctx.Request.FormValue("user_id")
	logs.Debug(user_id)
	rst,err:=this.models.GetRolesByUser(user_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"操作数据库失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}

func (this userRolesController)GetOtherRoles(ctx *context.Context){
	ctx.Request.ParseForm()
	user_id :=ctx.Request.FormValue("user_id")
	logs.Debug(user_id)
	if user_id==""{
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"请选择需要查询的用户")
		return
	}
	rst,err:=this.models.GetOtherRoles(user_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"查询未授权角色信息失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}
func (this userRolesController)Auth(ctx *context.Context){
	ctx.Request.ParseForm()
	ijs := ctx.Request.FormValue("JSON")
	logs.Error(ijs)

	cok, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cok.Value)
	if err != nil {
		logs.Error(err)
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/", http.StatusMovedPermanently)
		return
	}

	err=this.models.Auth(jclaim.User_id,ijs)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"授权用户角色信息失败",err)
		return
	}else{
		hret.WriteHttpOkMsgs(ctx.ResponseWriter,"授权用户角色信息成功")
		return
	}
}

func (this userRolesController)Revoke(ctx *context.Context){
	ctx.Request.ParseForm()
	user_id:=ctx.Request.FormValue("user_id")
	role_id:=ctx.Request.FormValue("role_id")
	logs.Debug(user_id,role_id)
	err:=this.models.Revoke(user_id,role_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"删除数据失败",err)
		return
	}else{
		hret.WriteHttpOkMsgs(ctx.ResponseWriter,"撤销用户角色授权成功")
		return
	}
}