package controllers

import (
	"text/template"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/utils/hret"

	"github.com/hzwy23/hauth/utils/token/hjwt"
	"github.com/hzwy23/hauth/utils/logs"
	"github.com/hzwy23/hauth/models"
)

type UserController struct {
	models *models.UserModel
}

var UserCtl = &UserController{
	new(models.UserModel),
}

func (UserController) Page(ctx *context.Context) {
	defer hret.HttpPanic()
	hz, _ := template.ParseFiles("./views/hauth/UserInfoPage.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}


func (this UserController)Get(ctx *context.Context){
	ctx.Request.ParseForm()
	offset:=ctx.Request.FormValue("offset")
	limit:=ctx.Request.FormValue("limit")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	rst,err:=this.models.GetDefault(jclaim.Domain_id,offset,limit)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 410, "查询数据库失败。")
		return
	}
	hret.WriteJson(ctx.ResponseWriter, rst)
}
