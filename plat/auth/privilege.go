package auth

import (
	"net/http"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/plat/mgr/sqlText"
	"github.com/hzwy23/hauth/plat/session"
)

type privilege struct {
}

var hp = new(privilege)

func Access(w http.ResponseWriter, r *http.Request) bool {
	return hp.Access(w, r)
}

func UserValid(user, passwd string) bool {
	return hp.UserValid(user, passwd)
}

//
func (this *privilege) Access(w http.ResponseWriter, r *http.Request) bool {

	userId := session.Get(w, r, "userId")
	userPd := session.Get(w, r, "userPasswd")

	return this.UserValid(userId, userPd)
}

func (this *privilege) UserValid(user, passwd string) bool {

	sql := sqlText.PLATFORM_RESOURCE_USERCHECK
	row := dbobj.QueryRow(sql, user)

	var ur string
	var pd string

	err := row.Scan(&ur, &pd)
	if err != nil {
		logs.Error(dbobj.GetErrorMsg(err), sql, user)
		return false
	}
	if ur == user && pd == passwd && ur != "" && pd != "" {
		return true
	} else {
		logs.Warn("user checkout failed.user or passwd is incorrect")
		return false
	}
}
