package resource

import (
	"net/http"

	"github.com/hzwy23/hcloud/plat/route"
)

type UserDomainRel struct {
	User_id   string
	Domain_id string
	route.RouteControl
}

func (this *UserDomainRel) Post(w http.ResponseWriter, r *http.Request) {

}
