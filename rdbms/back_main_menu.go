package rdbms

import (
	"io/ioutil"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/utils/hjwt"
)

func getBackMainMenu(ctx *context.Context) {

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	row := dbobj.QueryRow(sys_rdbms_013, jclaim.User_id)
	var url = "./views/theme/default/back_index.tpl"
	err = row.Scan(&url)
	if err != nil {
		url = "./views/theme/default/back_index.tpl"
		logs.Debug("获取默认主题")
	}

	h, err := ioutil.ReadFile(url)
	if err != nil {
		logs.Error(err)
		return
	}
	ctx.ResponseWriter.Write(h)
}

func init() {
	beego.Get("/v1/auth/main/menu/back", getBackMainMenu)
}
