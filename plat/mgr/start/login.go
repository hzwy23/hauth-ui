package start

import (
	"io/ioutil"
	"net/http"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"
	"github.com/hzwy23/hcloud/utils/hjwt"

	"github.com/hzwy23/hcloud/utils"
)

type LoginSystem struct {
	route.RouteControl
}

func (this *LoginSystem) Post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userId := r.FormValue("username")

	userPasswd := r.FormValue("password")

	psd, err := utils.Encrypt(userPasswd)
	if err != nil {
		logs.Error("decrypt passwd failed.", psd)
	}
	domainId := ""

	err = dbobj.QueryRow("SELECT distinct domain_id FROM bigdata.sys_user_domain_rel where user_id = ?", userId).Scan(&domainId)
	if err != nil {
		logs.Error(userId, " 用户没有指定的域", err)
	}

	if auth.UserValid(userId, psd) {

		err = session.Set(w, r, map[string]string{"userId": userId, "userPasswd": psd, "domainId": domainId})
		if err != nil {
			logs.Error("set session failed.", err)
		}
		sql := sqlText.PLATFORM_RESOURCE_LOGIN1
		row := dbobj.QueryRow(sql, userId)
		var url = "./views/theme/default/index.tpl"
		err = row.Scan(&url)
		if err != nil {
			url = "./views/theme/default/index.tpl"
			logs.Debug("获取默认主题")
		}
		h, err := ioutil.ReadFile(url)
		if err != nil {
			logs.Error(err)
			return
		}
		token := hjwt.GenToken()
		cookie := http.Cookie{Name: "Authorization", Value: token, Path: "/", MaxAge: 3600}
		http.SetCookie(w, &cookie)
		w.Write(h)
		logs.LogToDB(r, userId, true, "登陆成功")

	} else {
		logs.LogToDB(r, userId, false, "用户名或密码错误")

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func (this *LoginSystem) Get(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	if auth.Access(w, r) == false {

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	} else {
		userId := session.Get(w, r, "userId")
		sql := sqlText.PLATFORM_RESOURCE_LOGIN2
		row := dbobj.QueryRow(sql, userId)
		var url = "./views/theme/default/index.tpl"
		err := row.Scan(&url)
		if err != nil {
			url = "./views/theme/default/index.tpl"
			logs.Debug("获取默认主题")
		}
		h, err := ioutil.ReadFile(url)
		if err != nil {
			logs.Error(err)
			return
		}
		w.Write(h)
		logs.LogToDB(r, userId, true, "登陆成功")
	}
}
