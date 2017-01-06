package resource

import (
	"html/template"
	"net/http"

	"github.com/hzwy23/hcloud/plat/route"
)

type PlatMgrHelp struct {
	route.RouteControl
}

func (this *PlatMgrHelp) Get(w http.ResponseWriter, r *http.Request) {
	hz, _ := template.ParseFiles("./views/platMarHelp.tpl")
	hz.Execute(w, nil)
}
