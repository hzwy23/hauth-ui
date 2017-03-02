package controllers

import (
	"github.com/hzwy23/hauth/utils/hret"
	"io/ioutil"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/models"
	"github.com/hzwy23/hauth/utils/token/hjwt"
	"github.com/hzwy23/hauth/utils/logs"
	"html/template"
	"fmt"
)

type DomainShareControll struct{
	models *models.DomainShareModel
}

var DomainShareCtl = DomainShareControll{
	models:new(models.DomainShareModel),
}

func (DomainShareControll)Page(ctx *context.Context){
	defer hret.HttpPanic()
	ctx.Request.ParseForm()
	var domain_id = ctx.Request.FormValue("domain_id")
	logs.Debug("domain_id is :",domain_id)

	rst,err:=DomainCtl.models.GetRow(domain_id)

	if err!=nil{
		logs.Error(err)
		file,_:=ioutil.ReadFile("./views/hauth/domain_share_info.tpl")
		ctx.ResponseWriter.Write(file)
	} else {
		hz,_:=template.ParseFiles("./views/hauth/domain_share_info.tpl")
		hz.Execute(ctx.ResponseWriter,rst)
	}

}


func (this DomainShareControll)Get(ctx *context.Context){
	defer hret.HttpPanic()
	domain_id :=ctx.Request.FormValue("domain_id")
	if domain_id==""{
		cookie, _ := ctx.Request.Cookie("Authorization")
		jclaim, err := hjwt.ParseJwt(cookie.Value)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 410, "No Auth")
			return
		}
		domain_id = jclaim.Domain_id
	}
	fmt.Println(domain_id)
	rst,err:=this.models.Get(domain_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"操作数据库失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}

func (this DomainShareControll)UnAuth(ctx *context.Context){
	ctx.Request.ParseForm()
	domain_id :=ctx.Request.FormValue("domain_id")
	rst,err:=this.models.UnAuth(domain_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"查询未授权共享域信息失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}

func (this DomainShareControll)Post(ctx *context.Context){
	ctx.Request.ParseForm()
	domain_id := ctx.Request.FormValue("domain_id")
	target_domain_id :=ctx.Request.FormValue("target_domain_id")
	auth_level :=ctx.Request.FormValue("auth_level")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 410, "No Auth")
		return
	}
	err = this.models.Post(domain_id,target_domain_id,auth_level,jclaim.User_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"授权失败，操作数据库时出现异常")
		return
	}else{
		hret.WriteHttpOkMsgs(ctx.ResponseWriter,"域信息共享成功")
	}
}

func (this DomainShareControll)Delete(ctx *context.Context){
	ctx.Request.ParseForm()
	js :=ctx.Request.FormValue("JSON")
	err:=this.models.Delete(js)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"删除授权域失败",err)
		return
	}else{
		hret.WriteHttpOkMsgs(ctx.ResponseWriter,"删除授权信息成功")
	}
}

func (this DomainShareControll)Put(ctx *context.Context){
	ctx.Request.ParseForm()
	uuid:=ctx.Request.FormValue("uuid")
	level:=ctx.Request.FormValue("auth_level")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 410, "No Auth")
		return
	}

	err = this.models.Update(uuid,jclaim.User_id,level)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"更新域共享信息傻逼爱。")
		return
	}else{
		hret.WriteHttpOkMsgs(ctx.ResponseWriter,"更新域共享模式成功")
	}
}