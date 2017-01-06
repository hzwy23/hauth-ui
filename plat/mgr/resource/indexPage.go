package resource

import (
	"io/ioutil"
	"net/http"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"
)

type IndexPage struct {
	Res_id   string
	Res_name string
	route.RouteControl
}

func (this *IndexPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("no access"))
		logs.Error("没有权限")
		return
	}
	userId := session.Get(w, r, "userId")
	sql := sqlText.PLATFORM_RESOURCE_INDEX
	row := dbobj.QueryRow(sql, userId)
	var url = "./views/theme/default/back_index.tpl"
	err := row.Scan(&url)
	if err != nil {
		url = "./views/theme/default/back_index.tpl"
		logs.Debug("获取默认主题")
	}

	h, err := ioutil.ReadFile(url)
	if err != nil {
		logs.Error(err)
		return
	}
	w.Write(h)
}
