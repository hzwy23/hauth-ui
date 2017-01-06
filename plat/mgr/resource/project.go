package resource

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"

	"github.com/hzwy23/hcloud/logs"

	"github.com/hzwy23/hcloud/utils"
)

type ProjectMgr struct {
	Project_id     string
	Project_name   string
	Domain_up_id   string
	Project_status string
	Maintance_date string
	Domain_dept    string
	User_id        string
	cnt            int
	route.RouteControl
}

func (this *ProjectMgr) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	sql := sqlText.PLATFORM_RESOURCE_PROJECT1
	rows, err := dbobj.Query(sql, offset, limit+offset)
	defer rows.Close()
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		return
	}
	var oneLine ProjectMgr
	var rst []ProjectMgr
	for rows.Next() {
		err := rows.Scan(&oneLine.Project_id,
			&oneLine.Project_name,
			&oneLine.Domain_up_id,
			&oneLine.Project_status,
			&oneLine.Maintance_date,
			&oneLine.User_id)
		if err != nil {
			logs.Error("scan failed.", err)
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}
		rst = append(rst, oneLine)
	}
	var ret []ProjectMgr
	this.tree(rst, "-1", 1, &ret)
	this.WritePage(w, dbobj.Count("select count(*) from SYS_domain_info"), ret)
}

func (this *ProjectMgr) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()
	domainId := r.FormValue("domainId")
	domainDesc := r.FormValue("domainDesc")
	domainUpId := r.FormValue("domainUpId")
	domainStatus := r.FormValue("domainStatus")
	//校验
	if !utils.ValidAlphaNumber(domainId, 3, 10) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("域名编码格式错误,应为字母或数字组合，不为空"))
		return
	}

	if !utils.ValidAlphaNumber(domainUpId, 3, 10) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("上级域名编码为空，应为字母或数字组合，不为空"))
		return
	}
	//
	if !utils.ValidBool(domainStatus) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("域名状态❌"))
		return
	}

	logs.Debug(domainId, domainDesc, domainUpId, domainStatus)
	sql := sqlText.PLATFORM_RESOURCE_PROJECT3

	curtime := time.Now()
	userId := session.Get(w, r, "userId")
	err := dbobj.Exec(sql, domainId, domainDesc, domainUpId, domainStatus, curtime, userId)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("添加域失败" + domainId))
	}
	logs.LogToDB(r, userId, true, "更新域名信息成功")
}

func (this *ProjectMgr) Delete(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ijs := []byte(r.FormValue("JSON"))
	var js []ProjectMgr
	err := json.Unmarshal(ijs, &js)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("域编码格式错误,无法删除" + string(ijs)))
	}
	sql := sqlText.PLATFORM_RESOURCE_PROJECT2
	for _, val := range js {
		err := dbobj.Exec(sql, val.Project_id)
		if err != nil {
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("删除域失败" + val.Project_id))
		}
	}
	logs.LogToDB(r, session.Get(w, r, "userId"), true, "删除域名信息失败")
}

func (this *ProjectMgr) Put(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()
	sql := sqlText.PLATFORM_RESOURCE_PROJECT4
	domainId := r.FormValue("domainId")
	domainDesc := r.FormValue("domainDesc")
	domainUpId := r.FormValue("domainUpId")
	domainStatus := r.FormValue("domainStatus")

	err := dbobj.Exec(sql, domainDesc, domainUpId, domainStatus, domainId)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("更新域信息失败" + domainId))
		return
	}
	logs.LogToDB(r, session.Get(w, r, "userId"), true, "更新域名信息成功")
}

func (this *ProjectMgr) tree(node []ProjectMgr, id string, d int, result *[]ProjectMgr) {
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
			*result = append(*result, oneline)
			this.tree(node, val.Project_id, d+1, result)
		}
	}
}
