package resource

import (
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"
)

type GoEntry struct {
	route.RouteControl
}

func (this *GoEntry) Post(w http.ResponseWriter, r *http.Request) {

}

func (this *GoEntry) Get(w http.ResponseWriter, r *http.Request) {

	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		logs.Info("权限不足")
		return
	}

	r.ParseForm()
	id := r.FormValue("Id")
	var userId = session.Get(w, r, "userId")
	sql := sqlText.PLATFORM_RESOURCE_ENTRY1

	row := dbobj.QueryRow(sql, userId, id)

	var url string
	err := row.Scan(&url)
	if err != nil {
		logs.Error("cant not fetch menu_url", err)
		url = "./views/theme/default/sysconfig.tpl"
	}

	hz, _ := ioutil.ReadFile(url)
	w.Write(hz)
}
