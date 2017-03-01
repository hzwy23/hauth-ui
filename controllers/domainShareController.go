package controllers

import (
	"github.com/hzwy23/hauth/utils/hret"
	"io/ioutil"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/models"
	"github.com/hzwy23/hauth/utils/token/hjwt"
	"github.com/hzwy23/hauth/utils/logs"
	"html/template"
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
	rst,err:=this.models.Get(domain_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"操作数据库失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}