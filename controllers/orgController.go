package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"

	"github.com/astaxie/beego/context"

	"github.com/hzwy23/hauth/models"
	"github.com/hzwy23/hauth/utils"
	"github.com/hzwy23/hauth/utils/hret"
	"github.com/hzwy23/hauth/utils/logs"
	"github.com/hzwy23/hauth/utils/token/hjwt"
)

type OrgController struct {
	models *models.OrgModel
}

var OrgCtl = &OrgController{
	models: new(models.OrgModel),
}

func (OrgController) GetOrgPage(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/hauth/org_page.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func (this OrgController) GetSysOrgInfo(ctx *context.Context) {

	domain_id:=ctx.Request.FormValue("domain_id")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	if domain_id==""{
		domain_id = jclaim.Domain_id
	}
	rst, err := this.models.Get(domain_id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 417, "操作数据库失败")
	}
	tops := this.getOrgTops(rst)
	var ret []models.SysOrgInfo
	for _, val := range tops {
		var tmp []models.SysOrgInfo
		this.orgTree(rst, val.Org_unit_id, 2, &tmp)
		val.Org_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	hret.WriteJson(ctx.ResponseWriter,ret)
}

func (this OrgController) DeleteOrgInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	orgList := ctx.Request.FormValue("JSON")
	var mjs []models.SysOrgInfo
	err := json.Unmarshal([]byte(orgList), &mjs)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "delete org info failed.", err)
		return
	}

	err = this.models.Delete(mjs)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 418, "操作数据库失败")
		return
	} else {
		hret.WriteHttpOkMsgs(ctx.ResponseWriter, "delete org info successfully.")
		return
	}
}

func (this OrgController) UpdateOrgInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	org_unit_id := ctx.Request.FormValue("Id")
	org_unit_desc := ctx.Request.FormValue("Org_unit_desc")
	up_org_id := ctx.Request.FormValue("Up_org_id")
	start_date := ctx.Request.FormValue("Start_date")
	end_date := ctx.Request.FormValue("End_date")

	maintance_user := jclaim.User_id
	org_status_id := "0"
	if utils.AGTEB(start_date, end_date) {
		org_status_id = "1"
	}
	err = this.models.Update(org_unit_desc, up_org_id, org_status_id,
		start_date, end_date, maintance_user, org_unit_id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "modify org info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "modify org info successfully")
}

func (this OrgController) InsertOrgInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	org_unit_id := ctx.Request.FormValue("Org_unit_id")
	org_unit_desc := ctx.Request.FormValue("Org_unit_desc")
	up_org_id := ctx.Request.FormValue("Up_org_id")
	domain_id := ctx.Request.FormValue("Domain_id")
	start_date := ctx.Request.FormValue("Start_date")
	end_date := ctx.Request.FormValue("End_date")
	id := domain_id + "_join_" + org_unit_id
	create_user := jclaim.User_id
	maintance_user := jclaim.User_id
	org_status_id := "0"
	if utils.AGTEB(start_date, end_date) {
		org_status_id = "1"
	}

	err = this.models.Post(org_unit_id, org_unit_desc, up_org_id, org_status_id,
		domain_id, start_date, end_date, create_user, maintance_user, id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "add new org info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "add new org info successfully")
}

func (OrgController) getOrgTops(node []models.SysOrgInfo) []models.SysOrgInfo {
	var ret []models.SysOrgInfo
	for _, val := range node {
		flag := true
		for _, iv := range node {
			if val.Up_org_id == iv.Org_unit_id {
				flag = false
			}
		}
		if flag {
			ret = append(ret, val)
		}
	}
	return ret
}

func (this OrgController) orgTree(node []models.SysOrgInfo, id string, d int, result *[]models.SysOrgInfo) {
	var oneline models.SysOrgInfo
	for _, val := range node {
		if val.Up_org_id == id {
			oneline = val
			oneline.Org_dept = strconv.Itoa(d)
			*result = append(*result, oneline)
			this.orgTree(node, val.Org_unit_id, d+1, result)
		}
	}
}

func (this OrgController) GetOrgInfoByDomain(ctx *context.Context) {

	ctx.Request.ParseForm()
	domainid := ctx.Request.FormValue("domain_id")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	rst, err := this.models.GetOrgByDomainId(jclaim.Org_id, jclaim.Domain_id, domainid)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 419, "操作数据库失败")
		return
	}
	tops := this.getOrgTops(rst)
	var ret []models.SysOrgInfo
	for _, val := range tops {
		var tmp []models.SysOrgInfo
		this.orgTree(rst, val.Org_unit_id, 2, &tmp)
		val.Org_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	hret.WriteJson(ctx.ResponseWriter, ret)
}
