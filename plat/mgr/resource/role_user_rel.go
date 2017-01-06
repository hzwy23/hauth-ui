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
)

type RoleUserRel struct {
	Uuid           string
	User_id        string
	Role_id        string
	Role_name      string
	Maintance_date string
	Maintance_user string
	cnt            int
	route.RouteControl
}
type RoleUserRelPage struct {
	route.RouteControl
}

func (this *RoleUserRelPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := template.ParseFiles("./views/platform/resource/role_user_rel_page.tpl")
	hz.Execute(w, nil)
}

func (this *RoleUserRel) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	userId := r.FormValue("UserId")

	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	sql := sqlText.PLATFORM_RESOURCE_ROLEUSER1

	rows, err := dbobj.Query(sql, userId, offset, limit+offset)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		return
	}
	var one RoleUserRel
	var rst []RoleUserRel
	for rows.Next() {
		err := rows.Scan(&one.Uuid,
			&one.User_id,
			&one.Role_id,
			&one.Role_name,
			&one.Maintance_date,
			&one.Maintance_user)
		if err != nil {
			logs.Error(err)
			return
		}
		rst = append(rst, one)
	}
	this.WritePage(w, dbobj.Count("select count (*) from sys_role_user_relation T INNER JOIN SYS_ROLE_INFO I ON T.ROLE_ID = I.ROLE_ID where t.user_id = '"+userId+"'"), rst)
}

type SysRoleUserRelation struct {
	Uuid   string
	roleid string
	userid string
}

//hujian add 2016.7.19  用户角色
func (this *RoleUserRel) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	ijs := []byte(r.FormValue("JSON"))
	var rst []RoleUserRel
	err := json.Unmarshal(ijs, &rst)
	if err != nil {
		logs.Error(err.Error())
	}

	maintanceDate := time.Now().Format("2006-01-02")
	maintanceUser := session.Get(w, r, "userId")

	sql := sqlText.PLATFORM_RESOURCE_ROLEUSER3

	for _, val := range rst {
		err := dbobj.Exec(sql, val.Role_id, val.User_id, maintanceDate, maintanceUser)
		if err != nil {
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("插入角色域失败"))
		}
	}
}

func (this *RoleUserRel) Delete(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()
	mjson := []byte(r.FormValue("JSON"))
	var allrole []RoleUserRel
	err := json.Unmarshal(mjson, &allrole)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("json解析失败，请重新选择需要删除的角色信息"))
		return
	}
	sql := sqlText.PLATFORM_RESOURCE_ROLEUSER2
	for _, val := range allrole {
		logs.Debug("角色编码 ", val.Uuid)
		err := dbobj.Exec(sql, val.Uuid)
		if err != nil {
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("删除角色 " + val.Uuid + " 失败"))
			return
		}
		logs.Debug("删除角色", val.Uuid)
	}
}
