package controllers

import (
	"text/template"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/utils/hret"

	"github.com/hzwy23/hauth/utils/token/hjwt"
	"github.com/hzwy23/hauth/utils/logs"
	"github.com/hzwy23/hauth/models"
	"net/http"
	"github.com/hzwy23/hauth/utils"
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

	domain_id := ctx.Request.FormValue("domain_id")

	if domain_id =="" {
		cookie, _ := ctx.Request.Cookie("Authorization")
		jclaim, err := hjwt.ParseJwt(cookie.Value)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
			return
		}
		domain_id = jclaim.Domain_id
	}

	rst,err:=this.models.GetDefault(domain_id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 410, "查询数据库失败。")
		return
	}
	hret.WriteJson(ctx.ResponseWriter, rst)
}


func (this UserController) Post(ctx *context.Context) {

	ctx.Request.ParseForm()

	userId := ctx.Request.FormValue("userId")
	userDesc := ctx.Request.FormValue("userDesc")

	if !utils.ValidAlnumAndSymbol(userId) {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "user name must be alpha or number")
		return
	}
	//

	if !utils.ValidHanWord(userDesc) {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "user name must be words")
		return
	}
	//
	password := ctx.Request.FormValue("userPasswd")
	if !utils.ValidAlphaNumber(password, 6, 12) {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "user password must be 6-12 bits")
		return
	}

	userPasswd, err := utils.Encrypt(ctx.Request.FormValue("userPasswd"))
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "email account is not valid.please check your email address.", err)
		return
	}
	userStatus := ctx.Request.FormValue("userStatus")
	userEmail := ctx.Request.FormValue("userEmail")
	userPhone := ctx.Request.FormValue("userPhone")
	userOrgUnitId := ctx.Request.FormValue("userOrgUnitId")

	//
	if !utils.ValidEmail(userEmail) {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "email account is not valid.please check your email address.", err)
		return
	}
	//
	if !utils.ValidMobile(userPhone) {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "phone number is not valid.please check your phone number.", err)
		return
	}

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	err=this.models.Post(userId,userPasswd,userDesc,userStatus,jclaim.User_id,userEmail,userPhone,userOrgUnitId)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"插入用户信息失败。")
		return
	}else{
		logs.Info("新增用户信息",userId)
		hret.WriteHttpOkMsgs(ctx.ResponseWriter,"新增用户信息成功")
	}
}


func (this UserController) Delete(ctx *context.Context){
	ijs := []byte(ctx.Request.FormValue("JSON"))

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	err=this.models.Delete(ijs,jclaim.User_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,err.Error())
		return
	}else{
		logs.Info("delete user id successfully.")
		hret.WriteHttpOkMsgs(ctx.ResponseWriter,"delete user id successfully.")
		return
	}
}


func (this UserController)Search(ctx *context.Context){
	ctx.Request.ParseForm()
	var org_id = ctx.Request.FormValue("org_id")
	var status_id = ctx.Request.FormValue("status_id")
	var domain_id = ctx.Request.FormValue("domain_id")
	if domain_id == "" {
		cookie, _ := ctx.Request.Cookie("Authorization")
		jclaim, err := hjwt.ParseJwt(cookie.Value)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
			return
		}
		domain_id = jclaim.Domain_id
	}
	logs.Debug(org_id,status_id)
	rst,err:=this.models.Search(org_id,status_id,domain_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"操作数据库失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}