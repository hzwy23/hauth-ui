package resource

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"

	"github.com/hzwy23/hcloud/utils"
)

type UserInfo struct {
	User_id          string
	User_name        string
	User_status_desc string
	User_create_date string
	User_owner       string
	User_email       string
	User_phone       string
	Org_unit_desc    string
	Domain_id        string
	Domain_name      string
	cnt              int
	route.RouteControl
}

type UserInfoPage struct {
	route.RouteControl
}

func (this *UserInfoPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := template.ParseFiles("./views/platform/resource/UserInfoPage.tpl")
	hz.Execute(w, nil)
}

func (this *UserInfo) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		logs.LogToDB(r, session.Get(w, r, "userId"), false, "没有权限访问")
		return
	}

	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	domainId := session.Get(w, r, "domainId")

	sql := sqlText.PLATFORM_RESOURCE_USER1
	//sort column
	// if sort is null,then use user_id as sort id
	if val := r.FormValue("sort"); val != "" {
		sql = strings.Replace(sql, "HZWSORTCOL", val, -1)
	} else {
		sql = strings.Replace(sql, "HZWSORTCOL", "user_id", -1)
	}

	if val := r.FormValue("order"); val != "" {
		sql = strings.Replace(sql, "HZWSORTAD", val, -1)
	} else {
		sql = strings.Replace(sql, "HZWSORTAD", "asc", -1)
	}

	row, err := dbobj.Query(sql, domainId, offset, limit+offset)
	defer row.Close()
	if err != nil {
		logs.Error(err)
		logs.LogToDB(r, session.Get(w, r, "userId"), false, err.Error())
		return
	}

	var rst []UserInfo
	err = dbobj.Scan(row, &rst)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("查询用户信息失败"))
		return
	}
	this.WritePage(w, dbobj.Count(sqlText.PLATFORM_RESOURCE_USER6), rst)
}

func (this *UserInfo) Post(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		logs.LogToDB(r, session.Get(w, r, "userId"), false, "没有权限访问")
		return
	}

	r.ParseForm()

	sql := sqlText.PLATFORM_RESOURCE_USERINFO1
	ssql := sqlText.PLATFORM_RESOURCE_USERINFO2
	stheme := `insert into  sys_user_theme(user_id,theme_id) values(?,?)`

	userId := r.FormValue("userId")
	userDesc := r.FormValue("userDesc")

	if !utils.ValidAlphaNumber(userId) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("账户必须是字母，数字组成"))
		return
	}
	//

	if !utils.ValidHanWord(userDesc) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("用户名必须是汉字与字母组成"))
		return
	}
	//
	password := r.FormValue("userPasswd")
	if !utils.ValidAlphaNumber(password, 6, 12) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("密码长度不正确，密码长度6－－12位"))
		return
	}

	userPasswd, err := utils.Encrypt(r.FormValue("userPasswd"))
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("密码加密失败"))
	}
	userStatus := r.FormValue("userStatus")
	userEmail := r.FormValue("userEmail")
	userPhone := r.FormValue("userPhone")
	userOrgUnitId := r.FormValue("userOrgUnitId")
	//
	if !utils.ValidEmail(userEmail) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("请输正确的邮箱账号"))
		return
	}
	//
	if !utils.ValidMobile(userPhone) {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("请输正确的电话号码"))
		return
	}
	owner := session.Get(w, r, "userId")

	tx, err := dbobj.Begin()
	// insert user details
	//
	_, err = tx.Exec(sql, userId, userDesc, owner, userEmail, userPhone, userOrgUnitId)
	if err != nil {
		tx.Rollback()
		logs.Error(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("新增用户失败"))
		return
	}

	// insert user passwd
	//
	_, err = tx.Exec(ssql, userId, userPasswd, userStatus)
	if err != nil {
		tx.Rollback()
		logs.Error("添加用户失败,写入用户,密码,状态信息失败", err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("添加用户失败,写入用户,密码,状态信息失败"))
		return
	}

	// insert theme info
	//
	_, err = tx.Exec(stheme, userId, "1001")
	if err != nil {
		tx.Rollback()
		logs.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("添加用户主题信息失败"))
		return
	}
	logs.LogToDB(r, owner, true, "新增用户成功")
	tx.Commit()
}

func (this *UserInfo) Delete(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		logs.LogToDB(r, session.Get(w, r, "userId"), false, "没有权限访问")
		return
	}

	sql := sqlText.PLATFORM_RESOURCE_USER2
	ssql := sqlText.PLATFORM_RESOURCE_USER3
	stheme := sqlText.PLATFORM_RESOURCE_USER7
	ijs := []byte(r.FormValue("JSON"))
	var js []UserInfo
	err := json.Unmarshal(ijs, &js)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("待删除的用户账号信息错误" + string(ijs)))
	}

	for _, val := range js {

		//判断用户是否在线
		//如果在线,则不允许删除用户
		if val.User_id == "admin" {
			w.WriteHeader(http.StatusForbidden)
			logs.LogToDB(r, session.Get(w, r, "userId"), false, "无法删除管理员用户")
			w.Write([]byte(val.User_id + " 是系统管理员用户,无法删除"))
			return
		}
		tx, _ := dbobj.Begin()
		// delete user details
		//
		_, err := tx.Exec(ssql, val.User_id)
		if err != nil {
			tx.Rollback()
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			logs.LogToDB(r, session.Get(w, r, "userId"), false, "删除用户密码信息失败")
			w.Write([]byte("删除用户与密码信息失败: " + val.User_id))
		}

		// delete user passwd
		//
		_, err = tx.Exec(sql, val.User_id)
		if err != nil {
			tx.Rollback()
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			logs.LogToDB(r, session.Get(w, r, "userId"), false, "删除用户失败")
			w.Write([]byte("删除用户失败: " + val.User_id))
		}

		// delete user theme info
		//
		_, err = tx.Exec(stheme, val.User_id)
		if err != nil {
			tx.Rollback()
			logs.Error(err.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("删除用户主题信息失败"))
			logs.LogToDB(r, session.Get(w, r, "userId"), false, "删除用户主题失败")
		}
		logs.Info("删除用户 " + val.User_id + " 信息成功")
		logs.LogToDB(r, session.Get(w, r, "userId"), true, "删除用户成功")
		tx.Commit()
	}
}

func (this *UserInfo) Put(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		logs.LogToDB(r, session.Get(w, r, "userId"), false, "没有权限访问")
		return
	}
	r.ParseForm()
	offline := r.FormValue("OfflineUser")
	if offline == "true" {
		curId := session.Get(w, r, "userId")
		userId := r.FormValue("UserId")
		if curId == userId {
			w.WriteHeader(http.StatusExpectationFailed)
			logs.LogToDB(r, session.Get(w, r, "userId"), false, "不能将自己强制下限")
			w.Write([]byte("不能自己将自己强制下线"))
		} else if userId == "admin" {
			w.WriteHeader(http.StatusExpectationFailed)
			logs.LogToDB(r, session.Get(w, r, "userId"), false, "无法将超级管理员强制下限")
			w.Write([]byte("无法将超级管理员强制下线"))
		} else {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("用户不在线"))
		}
	} else {
		this.updateUserInfo(w, r)
	}
}

func (this *UserInfo) updateUserInfo(w http.ResponseWriter, r *http.Request) {
	sql := sqlText.PLATFORM_RESOURCE_USER4
	userId := r.FormValue("userId")
	userDesc := r.FormValue("userDesc")
	userStatus := r.FormValue("userStatus")
	userEmail := r.FormValue("userEmail")
	userPhone := r.FormValue("userPhone")
	userOrgUnitId := r.FormValue("userOrgUnitId")
	err := dbobj.Exec(sql, userDesc, userPhone, userEmail, userOrgUnitId, userId)
	if err != nil {
		logs.Error("更新用户信息失败", err)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("更新用户信息失败"))
		return
	}
	if session.Get(w, r, "userId") == userId && userStatus != "0" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("不允许自己禁用自己"))
		return
	}
	sql = sqlText.PLATFORM_RESOURCE_USER5
	err = dbobj.Exec(sql, userStatus, userId)
	if err != nil {
		logs.Error("更新用户状体信息失败", err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("更新用户状态信息失败"))
		return
	}
}

/*
func init() {
	var err error
	ce, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		logs.Error("init cache failed.")
		return
	} else {
		fmt.Println("init cache success.")
	}
}
*/
