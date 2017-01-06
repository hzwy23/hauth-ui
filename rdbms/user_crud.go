package rdbms

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/utils"
	"github.com/hzwy23/hauth/utils/hjwt"
)

type usersOfOrg struct {
	User_id     string
	User_name   string
	Domain_desc string
}

type usersOfDomain struct {
	User_id       string
	User_name     string
	Org_unit_desc string
}

type UserInfo struct {
	User_id             string
	User_name           string
	User_status_desc    string
	User_create_date    string
	User_owner          string
	User_email          string
	User_phone          string
	Org_unit_id         string
	Org_unit_desc       string
	Domain_id           string
	Domain_name         string
	User_maintance_date string
	User_maintance_user string
}

type userHandle struct {
	beego.Controller
}

// 用户在查询时，能够获取到自己域中，机构层级比自己低的用户信息
// 也可以获取到子域中，所有用户信息
func (this *userHandle) Get() {

	this.Ctx.Request.ParseForm()
	offset, _ := strconv.Atoi(this.Ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(this.Ctx.Request.FormValue("limit"))

	cookie, _ := this.Ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, 310, "No Auth")
		return
	}

	row, err := dbobj.Query(sys_rdbms_017, jclaim.User_id, jclaim.Domain_id, jclaim.Org_id, offset, limit)
	defer row.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, 311, "Query User info faild.", err)
		return
	}

	var rst []UserInfo
	err = dbobj.Scan(row, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, 312, "Query User info faild.", err)
		return
	}
	hret.WriteBootstrapTableJson(this.Ctx.ResponseWriter, dbobj.Count(sys_rdbms_016, jclaim.User_id, jclaim.Domain_id, jclaim.Org_id), rst)
}

func (this *userHandle) Post() {
	this.Ctx.Request.ParseForm()

	stheme := `insert into sys_user_theme(user_id,theme_id) values(?,?)`

	userId := this.Ctx.Request.FormValue("userId")
	userDesc := this.Ctx.Request.FormValue("userDesc")

	if !utils.ValidAlnumAndSymbol(userId) {
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "user name must be alpha or number")
		return
	}
	//

	if !utils.ValidHanWord(userDesc) {
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "user name must be words")
		return
	}
	//
	password := this.Ctx.Request.FormValue("userPasswd")
	if !utils.ValidAlphaNumber(password, 6, 12) {
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "user password must be 6-12 bits")
		return
	}

	userPasswd, err := utils.Encrypt(this.Ctx.Request.FormValue("userPasswd"))
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "email account is not valid.please check your email address.", err)
		return
	}
	userStatus := this.Ctx.Request.FormValue("userStatus")
	userEmail := this.Ctx.Request.FormValue("userEmail")
	userPhone := this.Ctx.Request.FormValue("userPhone")
	userOrgUnitId := this.Ctx.Request.FormValue("userOrgUnitId")
	userDomainId := this.Ctx.Request.FormValue("domainId")
	//
	if !utils.ValidEmail(userEmail) {
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "email account is not valid.please check your email address.", err)
		return
	}
	//
	if !utils.ValidMobile(userPhone) {
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "phone number is not valid.please check your phone number.", err)
		return
	}

	cookie, _ := this.Ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, 310, "No Auth")
		return
	}

	tx, err := dbobj.Begin()
	// insert user details
	//
	_, err = tx.Exec(sys_rdbms_018, userId, userDesc, jclaim.User_id, userEmail, userPhone, userOrgUnitId, jclaim.User_id)
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "add new user failed.", err)
		return
	}

	// insert user passwd
	//
	_, err = tx.Exec(sys_rdbms_019, userId, userPasswd, userStatus)
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "add new user failed. set password failed.", err)
		return
	}

	// insert theme info
	//
	_, err = tx.Exec(stheme, userId, "1001")
	if err != nil {
		tx.Rollback()
		logs.Error(err.Error())
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "add new user failed. set user theme failed.", err)
		return
	}

	_, err = tx.Exec(sys_rdbms_039, userId, userDomainId, jclaim.User_id)
	if err != nil {
		tx.Rollback()
		logs.Error(err.Error())
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusExpectationFailed, "add new user failed. set user theme failed.", err)
		return
	}
	//logs.LogToDB(r, owner, true, "新增用户成功")
	tx.Commit()
}

func (this *userHandle) Put() {
	this.Ctx.Request.ParseForm()

	userId := this.Ctx.Request.FormValue("userId")
	userDesc := this.Ctx.Request.FormValue("userDesc")
	userEmail := this.Ctx.Request.FormValue("userEmail")
	userPhone := this.Ctx.Request.FormValue("userPhone")
	cookie, _ := this.Ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, 310, "No Auth")
		return
	}

	err = dbobj.Exec(sys_rdbms_021, userDesc, userPhone, userEmail, jclaim.User_id, userId)
	if err != nil {
		logs.Error("更新用户信息失败", err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusForbidden, "update user info failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(this.Ctx.ResponseWriter, "update user info successfully.")

}

func getUserCrudPage(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/platform/resource/UserInfoPage.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func getUsersOfOrg(ctx *context.Context) {
	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))
	orgid := ctx.Request.FormValue("org_unit_id")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	row, err := dbobj.Query(sys_rdbms_058, jclaim.User_id, jclaim.Domain_id, orgid, offset, limit)
	defer row.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "Query User info faild.", err)
		return
	}

	var rst []usersOfOrg
	err = dbobj.Scan(row, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "Query User info faild.", err)
		return
	}
	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_059, jclaim.User_id, jclaim.Domain_id, orgid), rst)
}

func getUsersOfDomains(ctx *context.Context) {

	ctx.Request.ParseForm()
	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))
	domainid := ctx.Request.FormValue("domain_id")
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	var rst []usersOfDomain
	if jclaim.Domain_id == domainid {
		row, err := dbobj.Query(sys_rdbms_051, jclaim.User_id, domainid, jclaim.Org_id, offset, limit)
		defer row.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "Query User info faild.", err)
			return
		}
		err = dbobj.Scan(row, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "Query User info faild.", err)
			return
		}
		hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_049, jclaim.User_id, domainid, jclaim.Org_id), rst)
	} else {
		row, err := dbobj.Query(sys_rdbms_056, jclaim.User_id, jclaim.Domain_id, domainid, offset, limit)
		defer row.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "Query User info faild.", err)
			return
		}
		err = dbobj.Scan(row, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "Query User info faild.", err)
			return
		}
		hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_055, jclaim.User_id, jclaim.Domain_id, domainid), rst)

	}
}

func deleteUserInfo(ctx *context.Context) {

	ijs := []byte(ctx.Request.FormValue("JSON"))
	var js []UserInfo
	err := json.Unmarshal(ijs, &js)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "json info unmarsh1 failed.", err)
		return
	}

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}

	tx, _ := dbobj.Begin()
	for _, val := range js {
		//判断用户是否在线
		//如果在线,则不允许删除用户
		if val.User_id == "admin" {
			tx.Rollback()
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusForbidden, "admin is sysadmin, can not be deleted", err)
			return
		}

		// check user
		// can't delete yourself
		if jclaim.User_id == val.User_id {
			tx.Rollback()
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusForbidden, "you can't delete yourself.", err)
			return
		}

		_, err := tx.Exec(sys_rdbms_007, val.User_id)
		if err != nil {
			tx.Rollback()
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusForbidden, "delete user info failed.", err)
			return
		}
	}
	tx.Commit()
}

func forbiduser(ctx *context.Context) {
	ctx.Request.ParseForm()

	userId := ctx.Request.FormValue("userId")
	userStatus := ctx.Request.FormValue("userStatus")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	if jclaim.User_id == userId {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusForbidden, "forbidden change yourself status")
		return
	}

	err = dbobj.Exec(sys_rdbms_020, userStatus, userId)
	if err != nil {
		logs.Error("更新用户状体信息失败", err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "update user info failed..", err)
		return
	}
	logs.Info("update user status. user id is :", userId, " status is :", userStatus)
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "update user status successfully.")
}

func getUsersOfSpecialDomain(ctx *context.Context) {

	ctx.Request.ParseForm()
	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))
	domainid := ctx.Request.FormValue("domain_id")
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	var rst []UserInfo
	if jclaim.Domain_id == domainid {
		row, err := dbobj.Query(sys_rdbms_053, jclaim.User_id, domainid, jclaim.Org_id, offset, limit)
		defer row.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "Query User info faild.", err)
			return
		}
		err = dbobj.Scan(row, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "Query User info faild.", err)
			return
		}
		hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_052, jclaim.User_id, domainid, jclaim.Org_id), rst)
	} else {
		row, err := dbobj.Query(sys_rdbms_054, jclaim.User_id, jclaim.Domain_id, domainid, offset, limit)
		defer row.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "Query User info faild.", err)
			return
		}
		err = dbobj.Scan(row, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "Query User info faild.", err)
			return
		}
		hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_055, jclaim.User_id, jclaim.Domain_id, domainid), rst)
	}
}

func getUsersOfSpecialDomainAndOrg(ctx *context.Context) {

	ctx.Request.ParseForm()
	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))
	domainid := ctx.Request.FormValue("domain_id")
	orgid := ctx.Request.FormValue("org_unit_id")
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	var rst []UserInfo
	if jclaim.Domain_id == domainid {
		row, err := dbobj.Query(sys_rdbms_064, jclaim.User_id, domainid, jclaim.Org_id, orgid, offset, limit)
		defer row.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "Query User info faild.", err)
			return
		}
		err = dbobj.Scan(row, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "Query User info faild.", err)
			return
		}
		hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_062, jclaim.User_id, domainid, jclaim.Org_id, orgid), rst)
	} else {
		row, err := dbobj.Query(sys_rdbms_065, jclaim.User_id, jclaim.Domain_id, domainid, orgid, offset, limit)
		defer row.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 311, "Query User info faild.", err)
			return
		}
		err = dbobj.Scan(row, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 312, "Query User info faild.", err)
			return
		}
		hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(sys_rdbms_063, jclaim.User_id, jclaim.Domain_id, domainid, orgid), rst)
	}
}

func init() {
	beego.Get("/v1/auth/user/page", getUserCrudPage)
	beego.Router("/v1/auth/user", &userHandle{})
	beego.Post("/v1/auth/user/delete", deleteUserInfo)
	beego.Put("/v1/auth/user/forbid", forbiduser)
	beego.Get("/v1/auth/user/domain/get", getUsersOfDomains)
	beego.Get("/v1/auth/user/domain/details", getUsersOfSpecialDomain)
	beego.Get("/v1/auth/user/domain/org/list", getUsersOfSpecialDomainAndOrg)
	beego.Get("/v1/auth/user/org/get", getUsersOfOrg)
}
