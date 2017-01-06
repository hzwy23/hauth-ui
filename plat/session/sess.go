package session

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/hzwy23/hcloud/logs"
)

func SessionDestroy(w http.ResponseWriter, r *http.Request) {
	beego.GlobalSessions.SessionDestroy(w, r)
}

func Get(w http.ResponseWriter, r *http.Request, key string) string {

	seion, _ := beego.GlobalSessions.SessionStart(w, r)
	w.Header().Set("Content-Type", "text/html")

	if seion.Get(key) != nil {
		logs.Debug("Get info success in session :", key)
		return seion.Get(key).(string)
	} else {
		logs.Debug("Get info failed in session :", key)
		return ""
	}
}

func Set(w http.ResponseWriter, r *http.Request, m map[string]string) error {

	seion, _ := beego.GlobalSessions.SessionStart(w, r)
	w.Header().Set("Content-Type", "text/html")
	for key, val := range m {
		err := seion.Set(key, val)
		if err != nil {
			logs.Error("set user info failed in this session ->", err)
			return err
		}
	}
	seion.SessionRelease(w)
	return nil
}

func Delete(w http.ResponseWriter, r *http.Request, key string) error {

	sess, _ := beego.GlobalSessions.SessionStart(w, r)

	w.Header().Set("Content-Type", "text/html")

	err := sess.Delete(key)
	if err != nil {
		logs.Error("Delete user in session failed:", err)
		return err
	}
	return nil
}
