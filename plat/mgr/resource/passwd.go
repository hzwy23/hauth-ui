package resource

import (
	"net/http"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"

	"github.com/hzwy23/hcloud/utils"
)

type Passwd struct {
	route.RouteControl
}

func (this *Passwd) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()
	oriPasswd := r.FormValue("ora_passwd")
	newPasswd := r.FormValue("new_passwd")
	surePasswd := r.FormValue("sure_passwd")
	logs.Debug(r.Form)
	if newPasswd != surePasswd {
		logs.Error("new passwd is not match")
		return
	}
	pd := session.Get(w, r, "userPasswd")
	oriEn, err := utils.Encrypt(oriPasswd)
	if err != nil {
		logs.Error("Encrypt failed.")
	}
	if oriEn != pd {
		logs.Error("current user passwd error.")
		return
	}
	newPd, err := utils.Encrypt(newPasswd)
	if err != nil {
		logs.Error("Encrypt failed. New passwd is invalied.")
		return
	}
	ur := session.Get(w, r, "userId")

	sql := sqlText.PLATFORM_RESOURCE_PASSWD

	err = dbobj.Exec(sql, newPd, ur)
	if err != nil {
		logs.Error(dbobj.GetErrorMsg(err))
		return
	}
	logs.LogToDB(r, ur, true, "修改密码成功")
	//w.Write([]byte("测试错误，抛出错误信息"))
	//	w.WriteHeader(http.StatusForbidden)
	//	w.Write([]byte("hello world"))
}
