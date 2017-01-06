package rpc

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
)

type DomainInfo struct {
	Domain_id        string
	Domain_name      string
	Domain_up_id     string
	Domain_status_id string
	Create_date      string
	Create_user_id   string
	Modify_date      string
	Modify_user_id   string
	Domain_dept      string
}

func FindDomain(d []DomainInfo, id string) bool {
	for _, val := range d {
		if val.Domain_id == id {
			return true
		}
	}
	return false
}

func GetDomainsOfUser(ctx *context.Context) ([]DomainInfo, error) {
	r := ctx.Request
	r.ParseForm()
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return nil, err
	}

	rows, err := dbobj.Query(sys_rpc_001, jclaim.Domain_id)
	defer rows.Close()
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query domain info failed.", err)
		return nil, err
	}

	//	var oneLine ProjectMgr
	var rst []DomainInfo
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "query domain info failed.", err)
		return nil, err
	}

	var ret []DomainInfo
	for _, val := range getDomainTops(rst) {
		var tmp []DomainInfo
		dtree(rst, val.Domain_id, 2, &tmp)
		val.Domain_dept = "1"
		ret = append(ret, val)
		ret = append(ret, tmp...)
	}
	return ret, nil
}

func getDomainTops(node []DomainInfo) []DomainInfo {
	var ret []DomainInfo
	for _, val := range node {
		flag := true
		for _, iv := range node {
			if val.Domain_up_id == iv.Domain_id {
				flag = false
			}
		}
		if flag {
			ret = append(ret, val)
		}
	}
	return ret
}

func dtree(node []DomainInfo, id string, d int, result *[]DomainInfo) {
	var oneline DomainInfo
	for _, val := range node {
		if val.Domain_up_id == id {
			oneline.Domain_id = val.Domain_id
			oneline.Domain_name = val.Domain_name
			oneline.Domain_up_id = val.Domain_up_id
			oneline.Domain_status_id = val.Domain_status_id
			oneline.Create_date = val.Create_date
			oneline.Create_user_id = val.Create_user_id
			oneline.Domain_dept = strconv.Itoa(d)
			oneline.Modify_date = val.Modify_date
			oneline.Modify_user_id = val.Modify_user_id
			*result = append(*result, oneline)
			dtree(node, val.Domain_id, d+1, result)
		}
	}
}
