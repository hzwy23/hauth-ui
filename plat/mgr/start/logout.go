package start

import (
	"net/http"

	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"
)

type LogoutSystem struct {
	route.RouteControl
}

func (this *LogoutSystem) Post(w http.ResponseWriter, r *http.Request) {
	logs.LogToDB(r, session.Get(w, r, "userId"), true, "退出登陆")
	err := session.Delete(w, r, "userId")
	if err != nil {
		logs.Error("delete session info failed.", err)
	}

	err = session.Delete(w, r, "userPasswd")
	if err != nil {
		logs.Error("delete session info failed.", err)
	}

	session.SessionDestroy(w, r)
}
