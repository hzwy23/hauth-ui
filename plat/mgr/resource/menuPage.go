package resource

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/route"
)

type MenuPage struct {
	route.RouteControl
}

// Query menu info
func (this *MenuPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := template.ParseFiles("./views/platform/resource/menuPage.tpl")
	hz.Execute(w, nil)
}

// Insert menu into menu table
func (this *MenuPage) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := template.ParseFiles("./views/platform/resource/menuPage.tpl")
	hz.Execute(w, nil)
}

// Delete menu info from menu table
func (this *MenuPage) Delete(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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
