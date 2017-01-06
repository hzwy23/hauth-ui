package resource

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"
)

type ResourceRoleRel struct {
	Uuid      string
	Role_id   string
	Res_id    string
	Res_name  string
	Res_attr  string
	Res_up_id string
	Res_dept  string
	cnt       int
	route.RouteControl
}
type ResourceRoleRelPage struct {
	route.RouteControl
}

func (this *ResourceRoleRelPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := ioutil.ReadFile("./views/platform/resource/res_role_rel_page.tpl")
	w.Write(hz)
}

func (this *ResourceRoleRel) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	roleId := r.FormValue("RoleId")
	dataArgs := r.FormValue("RoleFlag")
	sql := sqlText.PLATFORM_RESOURCE_RESROLE1

	//参数等于0时,查询未获取的权限信息
	if dataArgs == "0" {
		sql = sqlText.PLATFORM_RESOURCE_RESROLE4
	}

	rows, err := dbobj.Query(sql, roleId)
	defer rows.Close()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	var one ResourceRoleRel
	var rst []ResourceRoleRel
	for rows.Next() {
		err := rows.Scan(&one.Uuid, &one.Res_id, &one.Res_name, &one.Res_attr, &one.Res_up_id)
		if err != nil {
			logs.Error(err.Error())
			return
		}
		one.Role_id = roleId
		rst = append(rst, one)
	}
	var ret []ResourceRoleRel

	//参数等于0时,查询未获取的权限信息
	if dataArgs == "0" {

		//get access resources
		sql = sqlText.PLATFORM_RESOURCE_RESROLE5
		rows, err = dbobj.Query(sql, roleId)
		if err != nil {
			logs.Error(err.Error())
			return
		}
		var urst []ResourceRoleRel
		for rows.Next() {
			err := rows.Scan(&one.Uuid, &one.Res_id, &one.Res_name, &one.Res_attr, &one.Res_up_id)
			if err != nil {
				logs.Error(err.Error())
				return
			}
			one.Role_id = roleId
			urst = append(urst, one)
		}

		//if getted resource id in ungetted resource up_id,
		//so show this resource.
		rst = this.addResinfo(urst, rst)
	}

	this.tree(rst, "-1", 1, &ret)

	ojs, err := json.Marshal(ret)
	if err != nil {
		logs.Error(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("获取资源信息失败"))
		return
	}
	w.Write(ojs)
}

//hujian add 2016.7.19  角色资源
func (this *ResourceRoleRel) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	ijs := []byte(r.FormValue("JSON"))
	var rst []ResourceRoleRel
	err := json.Unmarshal(ijs, &rst)
	if err != nil {
		logs.Error(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("解析参数失败"))
		return
	}

	sql := sqlText.PLATFORM_RESOURCE_RESROLE3

	for _, val := range rst {
		err := dbobj.Exec(sql, val.Role_id, val.Res_id)
		if err != nil {
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("插入角色失败"))
		}
	}
	logs.LogToDB(r, session.Get(w, r, "userId"), true, "更新角色与资源信息成功")
}

func (this *ResourceRoleRel) Delete(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	mjson := []byte(r.FormValue("JSON"))
	var all []RoleDomainRel
	err := json.Unmarshal(mjson, &all)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("json解析失败，请重新选择需要删除的角色信息"))
		return
	}

	tx, err := dbobj.Begin()
	if err != nil {
		logs.Error("dbobj begin failed.")
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("数据库开始事物处理失败，请联系管理员"))
		return
	}
	for _, val := range all {

		_, err := tx.Exec(sqlText.PLATFORM_RESOURCE_RESROLE2, val.Uuid)
		if err != nil {
			tx.Rollback()
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("删除角色 " + val.Uuid + " 失败"))
			return
		}
		logs.Debug("delete role info, uuid is : ", val.Uuid)
	}
	tx.Commit()
	logs.LogToDB(r, session.Get(w, r, "userId"), true, "删除角色与资源信息成功")
}

func (this *ResourceRoleRel) tree(node []ResourceRoleRel, id string, d int, result *[]ResourceRoleRel) {
	var oneline ResourceRoleRel
	for _, val := range node {
		if val.Res_up_id == id {
			oneline.Uuid = val.Uuid
			oneline.Role_id = val.Role_id
			oneline.Res_id = val.Res_id
			oneline.Res_name = val.Res_name
			oneline.Res_attr = val.Res_attr
			oneline.Res_up_id = val.Res_up_id
			oneline.Res_dept = strconv.Itoa(d)
			*result = append(*result, oneline)
			this.tree(node, val.Res_id, d+1, result)
		}
	}

}

func (this *ResourceRoleRel) addResinfo(urst, rst []ResourceRoleRel) []ResourceRoleRel {
	for _, val := range urst {
		for _, r := range rst {
			if r.Res_up_id == val.Res_id {
				var flag = false
				for _, e := range rst {
					if e.Res_id == val.Res_id {
						flag = true
					}
				}
				if flag == false {
					rst = append(rst, val)
				}
			}
		}
	}
	return rst
}
