package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/context"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/models"
	"github.com/hzwy23/hauth/utils"
	"github.com/hzwy23/hauth/utils/hret"
	"github.com/hzwy23/hauth/utils/logs"
	"github.com/hzwy23/hauth/utils/token/hjwt"
)

type DomainController struct {
	models *models.ProjectMgr
}

var DomainCtl = &DomainController{models: &models.ProjectMgr{}}

func (DomainController) GetDomainInfoPage(ctx *context.Context) {
	defer hret.HttpPanic()
	file, _ := ioutil.ReadFile("./views/hauth/domain_info.tpl")
	ctx.ResponseWriter.Write(file)
}

func (DomainController)SharePage(ctx *context.Context){
	defer hret.HttpPanic()
	file,_:=ioutil.ReadFile("./views/hauth/domain_share_info.tpl")
	ctx.ResponseWriter.Write(file)
}

func (this DomainController) GetDomainInfo(ctx *context.Context) {

	ctx.Request.ParseForm()
	offset := ctx.Request.FormValue("offset")
	limit := ctx.Request.FormValue("limit")

	rst, err := this.models.Get(offset, limit)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "查询数据库失败")
	}

	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count("select count(*) from SYS_domain_info"), rst)
}

func (this DomainController) PostDomainInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	domainId := ctx.Request.FormValue("domainId")
	domainDesc := ctx.Request.FormValue("domainDesc")
	domainStatus := ctx.Request.FormValue("domainStatus")
	//校验
	if !utils.ValidAlnumAndSymbol(domainId, 1, 30) {
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

	err = this.models.Post(domainId, domainDesc, domainStatus, jclaim.User_id)

	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "add domain info failed.", err)
		return
	}
}

func (this DomainController) DeleteDomainInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	ijs := []byte(ctx.Request.FormValue("JSON"))
	var js []models.ProjectMgr
	err := json.Unmarshal(ijs, &js)
	if err != nil {
		logs.Error(err)
		ctx.ResponseWriter.WriteHeader(http.StatusExpectationFailed)
		ctx.ResponseWriter.Write([]byte("域编码格式错误,无法删除" + string(ijs)))
		return
	}
	fmt.Println("delete info is :", js)
	err = this.models.Delete(js)
	if err != nil {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "get subdomain info failed.", err)
	} else {
		hret.WriteHttpOkMsgs(ctx.ResponseWriter, "删除域信息成功")
	}
}

func (this DomainController) UpdateDomainInfo(ctx *context.Context) {
	ctx.Request.ParseForm()

	domainId := ctx.Request.FormValue("domainId")
	domainDesc := ctx.Request.FormValue("domainDesc")
	domainStatus := ctx.Request.FormValue("domainStatus")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	err = this.models.Update(domainDesc, domainStatus, jclaim.User_id, domainId)
	if err != nil {
		logs.Error(err)
		ctx.ResponseWriter.WriteHeader(http.StatusExpectationFailed)
		ctx.ResponseWriter.Write([]byte("更新域信息失败" + domainId))
		return
	}
}

func (DomainController) getDomainTops(node []models.ProjectMgr) []models.ProjectMgr {
	var ret []models.ProjectMgr
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

func (this DomainController) dtree(node []models.ProjectMgr, id string, d int, result *[]models.ProjectMgr) {
	var oneline models.ProjectMgr
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
			this.dtree(node, val.Project_id, d+1, result)
		}
	}
}

func (this DomainController) GetDomainByUserInfo(ctx *context.Context) {

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	rst, err := this.models.GetDomainInfoByUser(jclaim.User_id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "查询数据库失败")
		return
	}
	var ret []models.ProjectMgr
	for _, val := range this.getDomainTops(rst) {
		var tmp []models.ProjectMgr
		this.dtree(rst, val.Project_id, 2, &tmp)
		val.Domain_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	hret.WriteJson(ctx.ResponseWriter, ret)
}

func (this DomainController) GetSubDomainInfo(ctx *context.Context) {
	ctx.Request.ParseForm()
	domainid := ctx.Request.FormValue("domain_id")
	rst, err := this.models.GetDomainInfoByUpId(domainid)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 333, "查询数据库失败")
		return
	}
	hret.WriteJson(ctx.ResponseWriter, rst)
}
