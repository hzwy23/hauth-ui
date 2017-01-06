package resource

import (
	"encoding/json"
	"html/template"
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

type RoleInfo struct {
	Role_id          string
	Role_name        string
	Role_owner       string
	Role_create_date string
	Role_status_desc string
	Role_status      string
	cnt              int
	route.RouteControl
}

type RoleInfoPage struct {
	route.RouteControl
}

func (this *RoleInfoPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := template.ParseFiles("./views/platform/resource/role_info_page.tpl")
	hz.Execute(w, nil)
}

func (this *RoleInfo) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	sql := sqlText.PLATFORM_RESOURCE_ROLEINFO1
	rows, err := dbobj.Query(sql, offset, limit+offset)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
	}
	var one RoleInfo
	var rst []RoleInfo
	for rows.Next() {
		err := rows.Scan(&one.Role_id,
			&one.Role_name,
			&one.Role_owner,
			&one.Role_create_date,
			&one.Role_status_desc,
			&one.Role_status)
		if err != nil {
			logs.Error(err)
		}
		rst = append(rst, one)
	}
	this.WritePage(w, dbobj.Count("select count(*) from sys_role_info"), rst)
}

//hujian add 2016.7.19  角色管理
func (this *RoleInfo) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	//取数据
	roleid := r.FormValue("role_id")
	rolename := r.FormValue("role_name")
	roleowner := session.Get(w, r, "userId")
	rolecdate := time.Now().Format("2006-01-02")
	rolestatus := 0
	//校验
	if !utils.ValidWord(roleid, 1, 30) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("请输入角色编码"))
		return
	}
	//
	if !utils.ValidHanAndWord(rolename, 1, 30) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("角色名称必须是汉字,字母,或者下划线的组合,并且长度不能小于30"))
		return
	}

	sql := sqlText.PLATFORM_RESOURCE_ROLEINFO2

	err := dbobj.Exec(sql, roleid, rolename, roleowner, rolecdate, rolestatus)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("插入角色信息失败"))
	}
}

func (this *RoleInfo) Delete(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	mjson := []byte(r.FormValue("JSON"))
	var allrole []RoleInfo
	err := json.Unmarshal(mjson, &allrole)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("json解析失败，请重新选择需要删除的角色信息"))
		return
	}
	sql := sqlText.PLATFORM_RESOURCE_ROLEINFO3
	for _, val := range allrole {
		err := dbobj.Exec(sql, val.Role_id)
		if err != nil {
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("删除角色 " + val.Role_id + " 失败"))
			return
		}
		logs.Debug("删除角色", val.Role_id)
	}
}

func (this *RoleInfo) Put(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	Role_id := r.FormValue("Role_id")
	Role_name := r.FormValue("Role_name")
	Role_status := r.FormValue("Role_status")

	sql := sqlText.PLATFORM_RESOURCE_ROLEINFO4

	err := dbobj.Exec(sql, Role_name, Role_status, Role_id)
	if err != nil {
		logs.Error(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("更新角色信息失败"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
