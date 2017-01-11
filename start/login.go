package start

import (
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/logs"

	"github.com/hzwy23/hauth/rdbms"

	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/token/hjwt"
	"github.com/hzwy23/hauth/utils"
)

func indexPage(ctx *context.Context) {
	defer hret.HttpPanic(func() {
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/", http.StatusMovedPermanently)
	})
	cok, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cok.Value)
	if err != nil {
		logs.Error(err)
		http.Redirect(ctx.ResponseWriter, ctx.Request, "/", http.StatusMovedPermanently)
		return
	}

	row := dbobj.QueryRow(platform_resource_login_1, jclaim.User_id)
	var url = "./views/theme/default/index.tpl"
	err = row.Scan(&url)
	if err != nil {
		url = "./views/theme/default/index.tpl"
		logs.Debug("get default theme.")
	}
	h, err := ioutil.ReadFile(url)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 402, "get default index html failed..")
		return
	}
	ctx.ResponseWriter.Write(h)
}

func loginSystem(ctx *context.Context) {
	ctx.Request.ParseForm()

	userId := ctx.Request.FormValue("username")

	userPasswd := ctx.Request.FormValue("password")

	psd, err := utils.Encrypt(userPasswd)
	if err != nil {
		logs.Error("decrypt passwd failed.", psd)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 400, "encrypt user passwd failed.")
		return
	}

	domainId := ""
	err = dbobj.QueryRow(platform_resource_login_2, userId).Scan(&domainId)
	if err != nil {
		logs.Error(userId, " 用户没有指定的域", err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 401, "user not specify domain.")
		return
	}

	orgid := ""
	err = dbobj.QueryRow(platform_resource_login_3, userId).Scan(&orgid)
	if err != nil {
		logs.Error(userId, " 用户没有指定的域", err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 401, "can't get org id of user")
		return
	}

	if ok, code, cnt, rmsg := rdbms.BasicAuth(userId, psd); ok {
		token := hjwt.GenToken(userId, domainId, orgid)
		cookie := http.Cookie{Name: "Authorization", Value: token, Path: "/", MaxAge: 3600}
		http.SetCookie(ctx.ResponseWriter, &cookie)
		hret.WriteHttpOkMsgs(ctx.ResponseWriter, "login successfully.")
	} else {
		emsg := hret.NewHttpErrMsg(code, rmsg, cnt)
		hret.WriteHttpErrMsg(ctx.ResponseWriter, emsg)
	}
}

func init() {

	beego.Get("/indexPage", indexPage)

	beego.Post("/login", loginSystem)

}
