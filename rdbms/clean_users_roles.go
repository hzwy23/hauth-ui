package rdbms

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
)

type mUsersList struct {
	User_id string
}

func cleanUsersRoles(ctx *context.Context) {
	ctx.Request.ParseForm()
	var rst []mUsersList
	err := json.Unmarshal([]byte(ctx.Request.FormValue("JSON")), &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "unmarshal failed.", err)
		return
	}
	tx, _ := dbobj.Begin()
	for _, val := range rst {
		_, err := tx.Exec(sys_rdbms_045, val.User_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, " delete roles of user failed.", err)
			return
		}
	}
	tx.Commit()
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "clean roles of user successfully.")
}

func init() {
	beego.Post("/v1/auth/role/users/clean", cleanUsersRoles)
}
