package rdbms

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
	"github.com/hzwy23/hauth/utils"
)

type mSysOrgInfo struct {
	Org_unit_id     string `json:"Org_unit_id"`
	Org_unit_desc   string `json:"Org_unit_desc"`
	Up_org_id       string `json:"Up_org_id"`
	Org_status_id   string `json:"Org_status_id"`
	Org_status_desc string `json:"Org_status_desc"`
	Domain_id       string `json:"Domain_id"`
	Domain_desc     string `json:"Domain_desc"`
	Start_date      string `json:"Start_date"`
	End_date        string `json:"End_date"`
	Create_date     string `json:"Create_date"`
	Maintance_date  string `json:"Maintance_date"`
	Create_user     string `json:"Create_user"`
	Maintance_user  string `json:"Maintance_user"`
	Code_number     string `json:"Code_number"`
	Org_dept        string `json:"Org_dept,omitempty"`
}

func getOrgPage(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/platform/resource/org_page.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func getSysOrgInfo(ctx *context.Context) {

	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	var rst []mSysOrgInfo
	rows, err := dbobj.Query(sys_rdbms_041, jclaim.Domain_id, jclaim.Org_id, offset, limit)
	if err != nil {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get org info failed.", err)
		return
	}

	err = dbobj.Scan(rows, &rst)
	if err != nil {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get org info failed.", err)
		return
	}
	tops := getOrgTops(rst)
	var ret []mSysOrgInfo
	for _, val := range tops {
		var tmp []mSysOrgInfo
		orgTree(rst, val.Org_unit_id, 2, &tmp)
		val.Org_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_042, jclaim.Domain_id, jclaim.Org_id), ret)
}

func deleteOrgInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	orgList := ctx.Request.FormValue("JSON")
	var mjs []mSysOrgInfo
	err := json.Unmarshal([]byte(orgList), &mjs)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "delete org info failed.", err)
		return
	}

	tx, _ := dbobj.Begin()
	for _, val := range mjs {

		_, err := tx.Exec(sys_rdbms_044, val.Org_unit_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "delete org info failed.", err)
			return
		}
	}
	tx.Commit()
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "delete org info successfully.")
}

func updateOrgInfo(ctx *context.Context) {
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
	org_status_id := 0
	if utils.AGTEB(start_date, end_date) {
		org_status_id = 1
	}
	err = dbobj.Exec(sys_rdbms_069, org_unit_desc, up_org_id, org_status_id,
		start_date, end_date, maintance_user, org_unit_id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "modify org info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "modify org info successfully")
}

func insertOrgInfo(ctx *context.Context) {
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
	org_status_id := 0
	if utils.AGTEB(start_date, end_date) {
		org_status_id = 1
	}

	err = dbobj.Exec(sys_rdbms_043, org_unit_id, org_unit_desc, up_org_id, org_status_id,
		domain_id, start_date, end_date, create_user, maintance_user, id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "add new org info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "add new org info successfully")
}

func getOrgTops(node []mSysOrgInfo) []mSysOrgInfo {
	var ret []mSysOrgInfo
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

func orgTree(node []mSysOrgInfo, id string, d int, result *[]mSysOrgInfo) {
	var oneline mSysOrgInfo
	for _, val := range node {
		if val.Up_org_id == id {
			oneline = val
			oneline.Org_dept = strconv.Itoa(d)
			*result = append(*result, oneline)
			orgTree(node, val.Org_unit_id, d+1, result)
		}
	}
}

func getOrgInfoByDomain(ctx *context.Context) {

	ctx.Request.ParseForm()
	domainid := ctx.Request.FormValue("domain_id")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	var rst []mSysOrgInfo
	if domainid != jclaim.Domain_id {
		rows, err := dbobj.Query(sys_rdbms_061, jclaim.Domain_id, domainid)
		if err != nil {
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get org info failed.", err)
			return
		}

		err = dbobj.Scan(rows, &rst)
		if err != nil {
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get org info failed.", err)
			return
		}
	} else {
		rows, err := dbobj.Query(sys_rdbms_060, jclaim.Org_id, jclaim.Domain_id, domainid)
		if err != nil {
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get org info failed.", err)
			return
		}

		err = dbobj.Scan(rows, &rst)
		if err != nil {
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get org info failed.", err)
			return
		}

	}
	tops := getOrgTops(rst)
	var ret []mSysOrgInfo
	for _, val := range tops {
		var tmp []mSysOrgInfo
		orgTree(rst, val.Org_unit_id, 2, &tmp)
		val.Org_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	hret.WriteJson(ctx.ResponseWriter, ret)
}

func init() {
	beego.Get("/v1/auth/relation/domain/org", getOrgInfoByDomain)
	beego.Get("/v1/auth/resource/org/page", getOrgPage)
	beego.Get("/v1/auth/resource/org/get", getSysOrgInfo)
	beego.Post("/v1/auth/resource/org/insert", insertOrgInfo)
	beego.Post("/v1/auth/resource/org/delete", deleteOrgInfo)
	beego.Put("/v1/auth/resource/org/update", updateOrgInfo)
}
