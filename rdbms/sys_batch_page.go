package rdbms

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
)

type mbatchUsers struct {
	User_id string `json:"user_id"`
	Role_id string `json:"role_id"`
}

func getBatchPage(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/platform/resource/sys_batch_page.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func batchGrants(ctx *context.Context) {
	ctx.Request.ParseForm()
	users := ctx.Request.FormValue("Users")
	roles := ctx.Request.FormValue("Roles")
	var ret []mbatchUsers
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	var user []mbatchUsers
	var role []mbatchUsers
	err = json.Unmarshal([]byte(users), &user)
	err = json.Unmarshal([]byte(roles), &role)
	if err != nil {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "batch grant failed.", err)
		return
	}
	err = json.Unmarshal([]byte(roles), &role)
	if err != nil {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "batch grant failed.", err)
		return
	}
	var tmp mbatchUsers
	tx, _ := dbobj.Begin()
	for _, u := range user {
		for _, r := range role {
			_, err := tx.Exec(sys_rdbms_048, u.User_id+"-"+r.Role_id, r.Role_id, u.User_id, jclaim.User_id)
			if err != nil {
				tmp.Role_id = r.Role_id
				tmp.User_id = u.User_id
				ret = append(ret, tmp)
			}
		}
	}
	tx.Commit()
	if len(ret) > 0 {
		hret.WriteHttpOkMsg(ctx.ResponseWriter, hret.HttpOkMsg{Version: "v1.0", Reply_code: 210, Reply_msg: "batch grant complete. but there are some role can't grant to users", Data: ret})
	} else {
		hret.WriteHttpOkMsgs(ctx.ResponseWriter, "batch grant successfully.")
	}

}

func init() {
	beego.Get("/v1/auth/batch/page", getBatchPage)
	beego.Post("/v1/auth/batch/grant", batchGrants)
}
