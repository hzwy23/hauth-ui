package rdbms

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
)

type ResourceRoleRel struct {
	Uuid      string
	Role_id   string
	Res_id    string
	Res_name  string
	Res_attr  string
	Res_up_id string
	Res_dept  string
}

func getRoleResourcePage(ctx *context.Context) {
	hz, _ := ioutil.ReadFile("./views/platform/resource/res_role_rel_page.tpl")
	ctx.ResponseWriter.Write(hz)
}

func getRoleResourceRel(ctx *context.Context) {
	ctx.Request.ParseForm()
	roleId := ctx.Request.FormValue("RoleId")
	dataArgs := ctx.Request.FormValue("RoleFlag")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	var rst []ResourceRoleRel
	var ret []ResourceRoleRel
	//参数等于0时,查询未获取的权限信息
	if dataArgs == "0" {
		rows, err := dbobj.Query(sys_rdbms_030, jclaim.User_id, jclaim.User_id, roleId)
		defer rows.Close()
		if err != nil {
			logs.Error(err.Error())
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get resource info failed.", err)
			return
		}
		var one ResourceRoleRel
		for rows.Next() {
			err := rows.Scan(&one.Uuid, &one.Res_id, &one.Res_name, &one.Res_attr, &one.Res_up_id)
			if err != nil {
				logs.Error(err.Error())
				hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get resource info failed.", err)
				return
			}
			one.Role_id = roleId
			rst = append(rst, one)
		}

		//get access resources
		rows, err = dbobj.Query(sys_rdbms_031, roleId)
		defer rows.Close()
		if err != nil {
			logs.Error(err.Error())
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get resource info failed.", err)
			return
		}
		var urst []ResourceRoleRel
		for rows.Next() {
			err := rows.Scan(&one.Uuid, &one.Res_id, &one.Res_name, &one.Res_attr, &one.Res_up_id)
			if err != nil {
				logs.Error(err.Error())
				hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get resource info failed.", err)
				return
			}
			one.Role_id = roleId
			urst = append(urst, one)
		}
		rst = addResinfo(urst, rst)
	} else {
		rows, err := dbobj.Query(sys_rdbms_029, roleId)
		defer rows.Close()
		if err != nil {
			logs.Error(err.Error())
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get resource info failed.", err)
			return
		}
		var one ResourceRoleRel
		for rows.Next() {
			err := rows.Scan(&one.Uuid, &one.Res_id, &one.Res_name, &one.Res_attr, &one.Res_up_id)
			if err != nil {
				logs.Error(err.Error())
				hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get resource info failed.", err)
				return
			}
			one.Role_id = roleId
			rst = append(rst, one)
		}
	}

	tree(rst, "-1", 1, &ret)
	ojs, err := json.Marshal(ret)
	if err != nil {
		logs.Error(err.Error())
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "get resource info failed.", err)
		return
	}
	ctx.ResponseWriter.Write(ojs)
}

/*
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
*/

func postRoleResourceRel(ctx *context.Context) {

	ctx.Request.ParseForm()
	ijs := []byte(ctx.Request.FormValue("JSON"))
	var rst []ResourceRoleRel
	err := json.Unmarshal(ijs, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "package json data failed.", err)
		return
	}

	tx, _ := dbobj.Begin()
	for _, val := range rst {
		_, err := tx.Exec(sys_rdbms_032, val.Role_id, val.Res_id)
		if err != nil {
			tx.Rollback()
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "add role and resource relation failed.", err)
			return
		}
	}
	tx.Commit()
}

func deleteRoleResourceRel(ctx *context.Context) {
	ctx.Request.ParseForm()

	mjson := []byte(ctx.Request.FormValue("JSON"))
	var all []ResourceRoleRel
	err := json.Unmarshal(mjson, &all)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "package json data failed.", err)
		return
	}

	tx, err := dbobj.Begin()
	if err != nil {
		logs.Error("dbobj begin failed.")
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "dbobj begin failed.", err)
		return
	}

	for _, val := range all {
		_, err := tx.Exec(sys_rdbms_033, val.Uuid)
		if err != nil {
			tx.Rollback()
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "delete role and resource relation failed.", err)
			return
		}
		logs.Info("delete role and resource relation info, uuid is : ", val.Uuid)
	}
	tx.Commit()
}

func tree(node []ResourceRoleRel, id string, d int, result *[]ResourceRoleRel) {
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
			tree(node, val.Res_id, d+1, result)
		}
	}
}

func addResinfo(urst, rst []ResourceRoleRel) []ResourceRoleRel {
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

func init() {
	beego.Get("/v1/auth/role/resource/page", getRoleResourcePage)
	beego.Get("/v1/auth/role/resource/get", getRoleResourceRel)
	beego.Post("/v1/auth/role/resource/post", postRoleResourceRel)
	beego.Post("/v1/auth/role/resource/delete", deleteRoleResourceRel)
}
