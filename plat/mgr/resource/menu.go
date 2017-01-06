package resource

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"

	"github.com/hzwy23/hcloud/utils"
)

type Menu struct {
	route.RouteControl
}

// Query menu info
func (this *Menu) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()
	menuId := r.FormValue("Id")

	sql := sqlText.PLATFORM_RESOURCE_MENU1

	rows, err := dbobj.Query(sql, menuId)
	defer rows.Close()
	if err != nil {
		logs.Error("query failed.", sql)
	}

	var retSet []utils.TreeMenuStruct
	var one utils.TreeMenuStruct
	var unResultSet []utils.TreeMenuStruct

	for rows.Next() {
		err := rows.Scan(&one.Menu_icon, &one.Menu_id, &one.Menu_name, &one.Menu_route, &one.Menu_up_id, &one.Menu_leaf_flag, &one.Menu_img, &one.Menu_color)
		if err != nil {
			logs.Error("get row failed.", err)
		}
		unResultSet = append(unResultSet, one)
	}

	utils.GetJSONMenuTree(unResultSet, menuId, 1, &retSet)
	ojs, err := json.Marshal(retSet)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("请求菜单信息失败，请联系管理员"))
		return
	}
	w.Write(ojs)
}

// Insert menu into menu table
func (this *Menu) Post(w http.ResponseWriter, r *http.Request) {
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
	logs.LogToDB(r, session.Get(w, r, "userId"), true, "更新菜单信息成功")
}

// Delete menu info from menu table
func (this *Menu) Delete(w http.ResponseWriter, r *http.Request) {
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
	var abc = &test{"Delete", "jack", 22}

	//	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	//	err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")

	t, err := template.New("foo").Parse(`This student id is :{{.Id}},Name is :{{.Name}},Age is :{{.Age}}`)
	err = t.Execute(w, abc)
	if err != nil {
		logs.Error("hello world")
	}
}

func calcPaddingLeft(args ...interface{}) (rst int) {
	if len(args) == 1 {
		rst, _ = strconv.Atoi(args[0].(string))
		rst = rst*15 - 15
	} else {
		rst = 0
	}
	return
}
