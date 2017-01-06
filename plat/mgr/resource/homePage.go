package resource

import (
	"html/template"
	"net/http"

	"github.com/hzwy23/hcloud/plat/route"
)

type HomePage struct {
	route.RouteControl
}

func (this *HomePage) Get(w http.ResponseWriter, r *http.Request) {
	hz, _ := template.ParseFiles("./views/homepage.tpl")
	hz.Execute(w, nil)
}
