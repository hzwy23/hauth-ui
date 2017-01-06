package rdbms

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
	"github.com/hzwy23/hauth/utils"
)

func getDomainInfoPage(ctx *context.Context) {
	defer hret.HttpPanic()
	file, _ := ioutil.ReadFile("./views/platform/resource/projectPage.tpl")
	ctx.ResponseWriter.Write(file)
}

type ProjectMgr struct {
	Project_id            string
	Project_name          string
	Domain_up_id          string
	Project_status        string
	Maintance_date        string
	User_id               string
	Domain_maintance_date string
	Domain_maintance_user string
	Domain_dept           string
}

func getDomainInfo(ctx *context.Context) {

	ctx.Request.ParseForm()
	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	rows, err := dbobj.Query(sys_rdbms_034, jclaim.Domain_id, offset, limit)
	defer rows.Close()
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query domain info failed.", err)
		return
	}

	//	var oneLine ProjectMgr
	var rst []ProjectMgr
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query domain info failed.", err)
		return
	}

	var ret []ProjectMgr
	for _, val := range getDomainTops(rst) {
		var tmp []ProjectMgr
		dtree(rst, val.Project_id, 2, &tmp)
		val.Domain_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}

	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count("select count(*) from SYS_domain_info"), ret)
}

func postDomainInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	domainId := ctx.Request.FormValue("domainId")
	domainDesc := ctx.Request.FormValue("domainDesc")
	domainUpId := ctx.Request.FormValue("domainUpId")
	domainStatus := ctx.Request.FormValue("domainStatus")
	//校验
	if !utils.ValidAlnumAndSymbol(domainId, 3, 30) {
		ctx.ResponseWriter.WriteHeader(http.StatusExpectationFailed)
		ctx.ResponseWriter.Write([]byte("域名编码格式错误,应为字母或数字组合，不为空"))
		return
	}

	//
	if !utils.ValidBool(domainStatus) {
		ctx.ResponseWriter.WriteHeader(http.StatusExpectationFailed)
		ctx.ResponseWriter.Write([]byte("域名状态❌"))
		return
	}

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	err = dbobj.Exec(sys_rdbms_036, domainId, domainDesc, domainUpId, domainStatus, jclaim.User_id, jclaim.User_id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "add domain info failed.", err)
		return
	}
}

func deleteDomainInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	ijs := []byte(ctx.Request.FormValue("JSON"))
	var js []ProjectMgr
	err := json.Unmarshal(ijs, &js)
	if err != nil {
		logs.Error(err)
		ctx.ResponseWriter.WriteHeader(http.StatusExpectationFailed)
		ctx.ResponseWriter.Write([]byte("域编码格式错误,无法删除" + string(ijs)))
		return
	}

	for _, val := range js {

		var tmp []ProjectMgr
		rows, err := dbobj.Query(sys_rdbms_067, val.Project_id)
		defer rows.Close()
		if err != nil {
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "get subdomain info failed.", err)
			return
		}

		err = dbobj.Scan(rows, &tmp)
		if err != nil {
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "get subdomain info failed.", err)
			return
		}

		for _, id := range tmp {
			err := dbobj.Exec(sys_rdbms_037, id.Project_id)
			if err != nil {
				logs.Error(err)
				ctx.ResponseWriter.WriteHeader(http.StatusExpectationFailed)
				ctx.ResponseWriter.Write([]byte("删除域失败" + val.Project_id))
			}
		}
	}
}

func updateDomainInfo(ctx *context.Context) {
	ctx.Request.ParseForm()

	domainId := ctx.Request.FormValue("domainId")
	domainDesc := ctx.Request.FormValue("domainDesc")
	domainUpId := ctx.Request.FormValue("domainUpId")
	domainStatus := ctx.Request.FormValue("domainStatus")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	err = dbobj.Exec(sys_rdbms_038, domainDesc, domainUpId, domainStatus, jclaim.User_id, domainId)
	if err != nil {
		logs.Error(err)
		ctx.ResponseWriter.WriteHeader(http.StatusExpectationFailed)
		ctx.ResponseWriter.Write([]byte("更新域信息失败" + domainId))
		return
	}
}

func getDomainTops(node []ProjectMgr) []ProjectMgr {
	var ret []ProjectMgr
	for _, val := range node {
		flag := true
		for _, iv := range node {
			if val.Domain_up_id == iv.Project_id {
				flag = false
			}
		}
		if flag {
			ret = append(ret, val)
		}
	}
	return ret
}

func dtree(node []ProjectMgr, id string, d int, result *[]ProjectMgr) {
	var oneline ProjectMgr
	for _, val := range node {
		if val.Domain_up_id == id {
			oneline.Project_id = val.Project_id
			oneline.Project_name = val.Project_name
			oneline.Domain_up_id = val.Domain_up_id
			oneline.Project_status = val.Project_status
			oneline.Maintance_date = val.Maintance_date
			oneline.User_id = val.User_id
			oneline.Domain_dept = strconv.Itoa(d)
			oneline.Domain_maintance_date = val.Domain_maintance_date
			oneline.Domain_maintance_user = val.Domain_maintance_user
			*result = append(*result, oneline)
			dtree(node, val.Project_id, d+1, result)
		}
	}
}

func getDomainByUserInfo(ctx *context.Context) {

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	rows, err := dbobj.Query(sys_rdbms_035, jclaim.Domain_id)
	defer rows.Close()
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query domain info failed.", err)
		return
	}

	//	var oneLine ProjectMgr
	var rst []ProjectMgr
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query domain info failed.", err)
		return
	}

	var ret []ProjectMgr
	for _, val := range getDomainTops(rst) {
		var tmp []ProjectMgr
		dtree(rst, val.Project_id, 2, &tmp)
		val.Domain_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	hret.WriteJson(ctx.ResponseWriter, ret)
}

func getSubDomainInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	domainid := ctx.Request.FormValue("domain_id")
	rows, err := dbobj.Query(sys_rdbms_067, domainid)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 345, "get subdomain info failed.", err)
		return
	}
	var rst []ProjectMgr
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 345, "get subdomain info failed.", err)
		return
	}
	hret.WriteJson(ctx.ResponseWriter, rst)
}

func init() {
	beego.Get("/v1/auth/domain/page", getDomainInfoPage)
	beego.Get("/v1/auth/domain/get", getDomainInfo)
	beego.Get("/v1/auth/user/domain", getDomainByUserInfo)
	beego.Post("/v1/auth/domain/post", postDomainInfo)
	beego.Post("/v1/auth/domain/delete", deleteDomainInfo)
	beego.Put("/v1/auth/domain/update", updateDomainInfo)
	beego.Get("/v1/auth/domain/subdomain", getSubDomainInfo)
}
