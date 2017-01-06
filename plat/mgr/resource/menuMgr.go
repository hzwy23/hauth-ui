package resource

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"

	"github.com/hzwy23/hcloud/utils"
)

type MenuMgr struct {
	route.RouteControl
}

// Query menu info
func (this *MenuMgr) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()
	//	menuId := r.FormValue("Id")
	menuId := "-1"
	logs.Debug("MenuId is :", menuId)

	sql := sqlText.PLATFORM_RESOURCE_MENUMGR1

	rows, err := dbobj.Query(sql, menuId)
	defer rows.Close()
	if err != nil {
		logs.Error("query failed.", sql)
	}

	var retSet []utils.TreeMenuStruct
	var one utils.TreeMenuStruct
	var unResultSet []utils.TreeMenuStruct

	for rows.Next() {
		err := rows.Scan(&one.Menu_icon, &one.Menu_id, &one.Menu_name, &one.Menu_route, &one.Menu_up_id, &one.Menu_leaf_flag)
		if err != nil {
			logs.Error("get row failed.", err)
		}
		logs.Debug(one)
		unResultSet = append(unResultSet, one)
	}

	utils.GetJSONMenuTree(unResultSet, menuId, 1, &retSet)
	ojs, err := json.Marshal(retSet)
	if err != nil {
		logs.Error(err)
	}
	w.Write(ojs)
}

// Insert menu into menu table
func (this *MenuMgr) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()

	type test struct {
		Id   string
		Name string
		Age  int
	}
	var abc = &test{"post", "mike", 21}

	//	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	//	err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")

	t, err := template.New("foo").Parse(`This student id is :{{.Id}},Name is :{{.Name}},Age is :{{.Age}}`)
	err = t.Execute(w, abc)
	if err != nil {
		logs.Error("hello world")
	}
}

// Delete menu info from menu table
func (this *MenuMgr) Delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Method)
	type test struct {
		Id   string
		Name string
		Age  int
	}
	var abc = &test{"Delete", "jack", 22}

	//	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	//	err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")

	t, err := template.New("foo").Parse(`This student id is :{{.Id}},Name is :{{.Name}},Age is :{{.Age}}`)
	err = t.Execute(w, abc)
	if err != nil {
		logs.Error("hello world")
	}
}
