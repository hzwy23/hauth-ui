package resource

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/route"
)

type ProjectPage struct {
	route.RouteControl
}

func (this *ProjectPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := ioutil.ReadFile("./views/platform/resource/projectPage.tpl")
	w.Write(hz)
}

func (this *ProjectPage) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := template.ParseFiles("./views/platform/resource/projectPage.tpl")
	hz.Execute(w, nil)
}
