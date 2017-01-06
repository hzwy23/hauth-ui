package resource

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"

	"github.com/hzwy23/hcloud/logs"
)

type ResInfo struct {
	Res_id        string
	Res_name      string
	Res_attr      string
	Res_attr_desc string
	Res_url       string
	Res_up_id     string
	Res_type      string
	Res_type_desc string
	Res_dept      string
	route.RouteControl
}

type ResInfoPage struct {
	route.RouteControl
}

func (this *ResInfoPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := ioutil.ReadFile("./views/platform/resource/res_info_page.tpl")
	w.Write(hz)
}

func (this *ResInfo) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	sql := sqlText.PLATFORM_RESOURCE_RESINFO1
	rows, err := dbobj.Query(sql, offset, limit+offset)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
	}
	var one ResInfo
	var rst []ResInfo
	for rows.Next() {
		err := rows.Scan(&one.Res_id,
			&one.Res_name,
			&one.Res_attr,
			&one.Res_attr_desc,
			&one.Res_url,
			&one.Res_up_id,
			&one.Res_type,
			&one.Res_type_desc)
		if err != nil {
			logs.Error(err)
		}
		rst = append(rst, one)
	}
	var ret []ResInfo
	this.resTree(rst, "-1", 1, &ret)
	this.WritePage(w, dbobj.Count("select count(*) from sys_resource_info"), ret)
}

func (this *ResInfo) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()

	resId := r.FormValue("resId")
	resDesc := r.FormValue("resDesc")
	resUrl := r.FormValue("resUrl")
	resUpId := r.FormValue("resUpId")
	resAttr := r.FormValue("resAttr")
	resType := r.FormValue("resType")

	sql := sqlText.PLATFORM_RESOURCE_RESINFO3
	err := dbobj.Exec(sql, resId, resDesc, resAttr, resUrl, resUpId, resType)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("新增菜单失败, 菜单编码是:" + resId))
	}

}

func (this *ResInfo) Put(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()
	resId := r.FormValue("resId")
	resDesc := r.FormValue("resDesc")
	resUrl := r.FormValue("resUrl")
	resUpId := r.FormValue("resUpId")
	resAttr := r.FormValue("resAttr")
	resConf := r.FormValue("resConf")
	resType := r.FormValue("resType")
	logs.Debug(resId, resDesc, resUrl, resUpId, resAttr, resConf, resType)
	sql := sqlText.PLATFORM_RESOURCE_RESINFO4
	err := dbobj.Exec(sql, resDesc, resAttr, resUrl, resUpId, resConf, resType, resId)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("新增菜单失败, 菜单编码是:" + resId))
	}
}

func (this *ResInfo) Delete(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.ParseForm()
	ijs := r.FormValue("JSON")
	var js []ResInfo
	err := json.Unmarshal([]byte(ijs), &js)
	if err != nil {
		logs.Error(err)
	}

	sql := sqlText.PLATFORM_RESOURCE_RESINFO2
	for _, val := range js {
		err := dbobj.Exec(sql, val.Res_id)
		if err != nil {
			logs.Error(err)
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("删除菜单失败." + val.Res_id))
			return
		}
	}
}

func (this *ResInfo) resTree(node []ResInfo, id string, d int, result *[]ResInfo) {
	var oneline ResInfo
	for _, val := range node {
		if val.Res_up_id == id {
			oneline.Res_id = val.Res_id
			oneline.Res_name = val.Res_name
			oneline.Res_attr = val.Res_attr
			oneline.Res_attr_desc = val.Res_attr_desc
			oneline.Res_url = val.Res_url
			oneline.Res_up_id = val.Res_up_id
			oneline.Res_dept = strconv.Itoa(d)
			oneline.Res_type = val.Res_type
			oneline.Res_type_desc = val.Res_type_desc
			*result = append(*result, oneline)
			this.resTree(node, val.Res_id, d+1, result)
		}
	}
}
