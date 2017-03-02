package controllers

import (
	"io/ioutil"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/models"
	"github.com/astaxie/beego/logs"
	"github.com/hzwy23/hauth/utils/hret"
	"github.com/hzwy23/hauth/utils"
	"strings"
)

type ResourceController struct {
	models *models.ResourceModel
}

var ResourceCtl = &ResourceController{
	new(models.ResourceModel),
}

func (ResourceController) Page(ctx *context.Context) {
	hz, _ := ioutil.ReadFile("./views/hauth/res_info_page.tpl")
	ctx.ResponseWriter.Write(hz)
}

func (this ResourceController)Get(ctx *context.Context){
	rst,err:=this.models.Get()
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"查询菜单信息失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}

func (this ResourceController)Query(ctx *context.Context){
	ctx.Request.ParseForm()
	res_id := ctx.Request.FormValue("res_id")
	rst,err:=this.models.Query(res_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"查询菜单信息失败")
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}

func (this ResourceController)QueryTheme(ctx *context.Context){
	ctx.Request.ParseForm()
	res_id :=ctx.Request.FormValue("res_id")
	theme_id :=ctx.Request.FormValue("theme_id")
	rst,err:=this.models.QueryTheme(res_id,theme_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"查询主题信息失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter,rst)
}

func (this ResourceController)Post(ctx *context.Context){
	ctx.Request.ParseForm()
	theme_id :=ctx.Request.FormValue("theme_id")
	res_type :=ctx.Request.FormValue("res_type")
	res_id :=ctx.Request.FormValue("res_id")
	res_name :=ctx.Request.FormValue("res_name")
	res_up_id :=ctx.Request.FormValue("res_up_id")
	res_url :=ctx.Request.FormValue("res_url")
	res_class :=ctx.Request.FormValue("res_class")
	res_img :=ctx.Request.FormValue("res_img")
	res_bg_color :=ctx.Request.FormValue("res_bg_color")
	group_id :=ctx.Request.FormValue("group_id")
	sort_id :=ctx.Request.FormValue("sort_id")
	res_attr := "1"
	if res_type == "0" || res_type == "4" {
		res_attr = "0"
	}
	if res_type == "0"{
		res_up_id = "-1"
	}

	if !utils.ValidAlphaNumber(res_id, 1, 30) {
		logs.Error("资源编码必须由1,30位字母或数字组成")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "资源编码必须由1,30位字母或数字组成")
		return
	}

	if strings.TrimSpace(res_name) == "" {
		logs.Error("菜单名称不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单名称不能为空")
		return
	}

	if strings.TrimSpace(res_url) == "" {
		logs.Error("菜单路由地址不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单路由地址不能为空")
		return
	}

	if strings.TrimSpace(res_up_id) == "" {
		logs.Error("菜单上级编码不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单上级编码不能为空")
		return
	}

	if strings.TrimSpace(res_class) == "" {
		logs.Error("菜单样式类型不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单样式类型不能为空")
		return
	}

	if strings.TrimSpace(res_attr) == "" {
		logs.Error("菜单属性值不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单属性值不能为空")
		return
	}

	if strings.TrimSpace(res_type) == "" {
		logs.Error("菜单类别不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单类别不能为空")
		return
	}

	if strings.TrimSpace(res_img) == "" {
		logs.Error("菜单图标不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单图标不能为空")
		return
	}

	if !utils.ValidNumber(group_id, 1, 2) {
		logs.Error("菜单分组信息必须是数字")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单分组信息必须是数字")
		return
	}

	if !utils.ValidNumber(sort_id, 1, 2) {
		logs.Error("菜单排序号必须是数字")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "菜单排序号必须是数字")
		return
	}

	logs.Debug(theme_id,res_id,res_type,res_id,res_name,res_up_id,res_url,res_class,res_img,res_bg_color,group_id,sort_id,res_attr)

	err:=this.models.Post(res_id,res_name,res_attr,res_up_id,res_type,theme_id,res_url,res_bg_color,res_class,group_id,res_img,sort_id)
	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"新增菜单资源信息失败",err)
		return
	}else{
		hret.WriteHttpOkMsgs(ctx.ResponseWriter,"新增资源信息成功")
		return
	}
}

func (this ResourceController)Delete(ctx *context.Context){

	ctx.Request.ParseForm()

	res_id := ctx.Request.FormValue("res_id")

	err:=this.models.Delete(res_id)

	if err!=nil{
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,419,"删除菜单信息失败")
		return
	}else{
		hret.WriteHttpOkMsgs(ctx.ResponseWriter, "remove resource successfully.")
	}
}