package rdbms

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/dbobj/utils"
	"strings"
)

type ResInfo struct {
	Res_id        string
	Res_name      string
	Res_attr      string
	Res_attr_desc string
	Res_up_id     string
	Res_type      string
	Res_type_desc string
	Res_dept      string
}

func getResourcePage(ctx *context.Context) {
	hz, _ := ioutil.ReadFile("./views/platform/resource/res_info_page.tpl")
	ctx.ResponseWriter.Write(hz)
}

func getResourceInfo(ctx *context.Context) {
	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))

	rows, err := dbobj.Query(sys_rdbms_040, offset, limit)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
	}

	var rst []ResInfo
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
	}

	tops := getTops(rst)

	var ret []ResInfo
	for _, val := range tops {
		var tmp []ResInfo
		resTree(rst, val.Res_id, 2, &tmp)
		val.Res_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count("select count(*) from sys_resource_info"), ret)
}

func getTops(node []ResInfo) []ResInfo {
	var ret []ResInfo
	for _, val := range node {
		flag := true
		for _, iv := range node {
			if val.Res_up_id == iv.Res_id {
				flag = false
			}
		}
		if flag {
			ret = append(ret, val)
		}
	}
	return ret
}

func resTree(node []ResInfo, id string, d int, result *[]ResInfo) {
	var oneline ResInfo
	for _, val := range node {
		if val.Res_up_id == id {
			oneline.Res_id = val.Res_id
			oneline.Res_name = val.Res_name
			oneline.Res_attr = val.Res_attr
			oneline.Res_attr_desc = val.Res_attr_desc
			oneline.Res_up_id = val.Res_up_id
			oneline.Res_dept = strconv.Itoa(d)
			oneline.Res_type = val.Res_type
			oneline.Res_type_desc = val.Res_type_desc
			*result = append(*result, oneline)
			resTree(node, val.Res_id, d+1, result)
		}
	}
}

func getResourceUpInfo(ctx *context.Context) {

	rows, err := dbobj.Query(sys_rdbms_071)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
	}

	var rst []ResInfo
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
	}

	tops := getTops(rst)

	var ret []ResInfo
	for _, val := range tops {
		var tmp []ResInfo
		resTree(rst, val.Res_id, 2, &tmp)
		val.Res_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	hret.WriteJson(ctx.ResponseWriter, ret)
}

func postSysResourceInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	res_id := ctx.Request.FormValue("resId")
	res_desc := ctx.Request.FormValue("resDesc")
	res_url := ctx.Request.FormValue("resUrl")
	res_up_id := ctx.Request.FormValue("resUpId")
	res_class := ctx.Request.FormValue("resClass")
	res_attr := ctx.Request.FormValue("resAttr")
	res_typ := ctx.Request.FormValue("resType")
	res_images := ctx.Request.FormValue("resImages")
	res_group_id := ctx.Request.FormValue("resGroupid")
	res_sort_id := ctx.Request.FormValue("resSortid")
	if !utils.ValidAlphaNumber(res_id,1,30){
		logs.Error("资源编码必须由1,30位字母或数字组成")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"资源编码必须由1,30位字母或数字组成")
		return
	}

	if strings.TrimSpace(res_desc)==""{
		logs.Error("菜单名称不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单名称不能为空")
		return
	}

	if strings.TrimSpace(res_url)==""{
		logs.Error("菜单路由地址不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单路由地址不能为空")
		return
	}

	if strings.TrimSpace(res_up_id)==""{
		logs.Error("菜单上级编码不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单上级编码不能为空")
		return
	}

	if strings.TrimSpace(res_class)==""{
		logs.Error("菜单样式类型不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单样式类型不能为空")
		return
	}

	if strings.TrimSpace(res_attr)==""{
		logs.Error("菜单属性值不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单属性值不能为空")
		return
	}

	if strings.TrimSpace(res_typ)==""{
		logs.Error("菜单类别不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单类别不能为空")
		return
	}

	if strings.TrimSpace(res_images)==""{
		logs.Error("菜单图标不能为空")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单图标不能为空")
		return
	}

	if !utils.ValidNumber(res_group_id,1,2){
		logs.Error("菜单分组信息必须是数字")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单分组信息必须是数字")
		return
	}

	if !utils.ValidNumber(res_sort_id,1,2){
		logs.Error("菜单排序号必须是数字")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter,333,"菜单排序号必须是数字")
		return
	}

	logs.Debug("res_id:", res_id, "res_desc:", res_desc, "res_url:", res_url,
		"res_up_id:", res_up_id, "res_class:", res_class, "res_attr:", res_attr,
		"res_type:", res_typ, "res_images:", res_images, "res_group_id:", res_group_id,
		"res_sort_id:", res_sort_id)

	tx, err := dbobj.Begin()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 319, "启动数据库事务处理失败")
		return
	}

	_, err = tx.Exec(sys_rdbms_072, res_id, res_desc, res_attr, res_up_id, res_typ)
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 321, "写入资源信息失败")
		return
	}

	_, err = tx.Exec(sys_rdbms_073, "1001", res_id, res_url, res_typ, "#339999", res_class, res_group_id, res_images, res_sort_id)
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 322, "写入资源信息失败")
		return
	}

	_, err = tx.Exec(sys_rdbms_074, "vertex_root_join_sysadmin", res_id)
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 323, "写入资源信息失败")
		return
	}
	tx.Commit()
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "add resource info successfully.")
}

func deleteResourceInfo(ctx *context.Context) {

	ctx.Request.ParseForm()

	ijs := ctx.Request.FormValue("JSON")

	var rst []ResInfo

	err := json.Unmarshal([]byte(ijs), &rst)

	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "解析前端参数失败")
		return
	}
	logs.Debug("info is :", rst)
	tx, err := dbobj.Begin()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 345, "开启事物处理失败")
		return
	}
	for _, val := range rst {

		if !HaveRightsById(ctx,val.Res_id){
			logs.Error("没有权限删除这个菜单",val.Res_id)
			tx.Rollback()
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 349, "没有权限删除这个菜单"+val.Res_id)
			return
		}

		_, err = tx.Exec(sys_rdbms_075, val.Res_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 346, "删除角色资源关系信息失败")
			return
		}
		_, err = tx.Exec(sys_rdbms_076, val.Res_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 347, "删除主题中资源信息失败")
			return
		}

		_, err = tx.Exec(sys_rdbms_077, val.Res_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 348, "删除资源列表中的资源信息失败")
			return
		}
	}
	tx.Commit()
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "remove resource successfully.")
}

func init() {
	beego.Get("/v1/auth/resource/page", getResourcePage)
	beego.Get("/v1/auth/resource/get", getResourceInfo)
	beego.Post("/v1/auth/resource/delete", deleteResourceInfo)
	beego.Post("/v1/auth/resource/post", postSysResourceInfo)
	beego.Get("/v1/auth/resource/get/upid", getResourceUpInfo)
}
