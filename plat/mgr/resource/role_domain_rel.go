package resource

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"

	"github.com/hzwy23/dbobj"

	"github.com/hzwy23/hcloud/utils"
)

type RoleDomainRel struct {
	Uuid         string
	Role_id      string
	Domain_id    string
	Domain_name  string
	Domain_up_id string
	Domain_dept  string
	cnt          int
	route.RouteControl
}
type RoleDomainRelPage struct {
	route.RouteControl
}

func (this *RoleDomainRelPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := ioutil.ReadFile("./views/platform/resource/role_domain_rel_page.tpl")
	w.Write(hz)
}

func (this *RoleDomainRel) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()
	roleId := r.FormValue("RoleId")
	sql := sqlText.PLATFORM_RESOURCE_ROLEDOMAIN1

	rows, err := dbobj.Query(sql, roleId)

	defer rows.Close()
	if err != nil {
		logs.Error(err)
		return
	}
	var one RoleDomainRel
	var rst []RoleDomainRel
	for rows.Next() {
		err := rows.Scan(&one.Uuid,
			&one.Role_id,
			&one.Domain_id,
			&one.Domain_name,
			&one.Domain_up_id)
		if err != nil {
			logs.Error(err)
			return
		}
		rst = append(rst, one)
	}

	var ret []RoleDomainRel
	this.tree(rst, "-1", 1, &ret)
	ojs, err := json.Marshal(ret)
	if err != nil {
		logs.Error(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("打包数据失败"))
		return
	}
	w.Write(ojs)
}

type SysRoleDomian struct {
	uuid     string
	roleid   string
	domainid string
}

//hujian add 2016.7.19  角色域
func (this *RoleDomainRel) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	//取数据
	uuid := r.FormValue("uuid")
	roleid := r.FormValue("role_id")
	domainid := r.FormValue("udomain_id")
	//校验
	if !utils.ValidAlphaNumber(uuid) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("请确认角色编码是否正确"))
		return
	}
	//
	if !utils.ValidAlphaNumber(roleid) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("请确认角色名称是否正确"))
		return
	}
	sql := sqlText.PLATFORM_RESOURCE_ROLEDOMAIN3
	err := dbobj.Exec(sql, uuid, roleid, domainid)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("插入角色域失败"))
	}

}

func (this *RoleDomainRel) Delete(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	mjson := []byte(r.FormValue("JSON"))
	var allrole []RoleDomainRel
	err := json.Unmarshal(mjson, &allrole)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("json解析失败，请重新选择需要删除的角色信息"))
		return
	}
	sql := sqlText.PLATFORM_RESOURCE_ROLEDOMAIN2
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

func (this *RoleDomainRel) Put(w http.ResponseWriter, r *http.Request) {

}

func (this *RoleDomainRel) tree(node []RoleDomainRel, id string, d int, result *[]RoleDomainRel) {
	var oneline RoleDomainRel
	for _, val := range node {
		if val.Domain_up_id == id {
			oneline.Uuid = val.Uuid
			oneline.Role_id = val.Role_id
			oneline.Domain_id = val.Domain_id
			oneline.Domain_up_id = val.Domain_up_id
			oneline.Domain_name = val.Domain_name
			oneline.Domain_dept = strconv.Itoa(d)
			*result = append(*result, oneline)
			this.tree(node, val.Domain_id, d+1, result)
		}
	}

}
