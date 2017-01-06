package start

import (
	"html/template"
	"net/http"

	"github.com/hzwy23/hcloud/plat/route"
)

type IndexSystem struct {
	route.RouteControl
}

func (this *IndexSystem) Get(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {

		huang, _ := template.ParseFiles("./views/login.tpl")

		huang.Execute(w, nil)

	} else {
		huang, _ := template.ParseFiles("./views/error.tpl")
		huang.Execute(w, nil)
	}
}
