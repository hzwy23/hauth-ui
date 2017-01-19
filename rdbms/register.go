package rdbms

import (
	"html/template"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/netease"
	"github.com/hzwy23/hauth/utils"
)

var retcode = map[int]string{
	416: "每一个手机号码每天最多能获取10条验证信息",
}

func getVerifyCode(ctx *context.Context) {
	ctx.Request.ParseForm()
	phone_number := ctx.Request.FormValue("inputPhoneNumber")

	var cnt = ""
	dbobj.QueryRow(sys_rdbms_068, phone_number).Scan(&cnt)
	if cnt == phone_number {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 308, "用户已经存在，请登录系统，如果忘记密码，请重置密码")
		return
	}

	msg, err := netease.SendCode(phone_number)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 498, "get message code failed.", err)
		return
	}
	if msg.Code != 200 {
		logs.Error("请求短信验证码失败", msg)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, msg.Code, retcode[msg.Code], msg)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "code is 200")
}

func httpRegisterPage(ctx *context.Context) {
	hz, err := template.ParseFiles("./views/platform/messageVerify.tpl")
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 404, "cant not find page", err)
		return
	}
	hz.Execute(ctx.ResponseWriter, nil)
}

func verifyCode(ctx *context.Context) {
	ctx.Request.ParseForm()
	phone_number := ctx.Request.FormValue("inputPhoneNumber")
	code := ctx.Request.FormValue("inputVerifyCode")
	pd := ctx.Request.FormValue("inputPassword")
	cpd := ctx.Request.FormValue("inputConfirmPassword")
	if pd != cpd {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 405, "两次输入密码不一致")
		return
	}
	msg, ret := netease.VerifyCode(phone_number, code)
	if ret {
		addDomain(phone_number)
		addOrg(phone_number)
		addUser(phone_number, pd)
		addrole(phone_number)
		hret.WriteHttpOkMsgs(ctx.ResponseWriter, msg)
	} else {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, msg, "校验码不正确，请重新输入。错误代码是：")
	}
}

func noVerifyRegister(ctx *context.Context) {
	ctx.Request.ParseForm()
	phone_number := ctx.Request.FormValue("inputPhoneNumber")
	pd := ctx.Request.FormValue("inputPassword")
	cpd := ctx.Request.FormValue("inputConfirmPassword")
	if pd != cpd {
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 405, "两次输入密码不一致")
		return
	}

	msg, err := addDomain(phone_number)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 404, "用户已经被注册", err)
		return
	}
	msg, err = addOrg(phone_number)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 404, "register org failed.", err)
		return
	}
	msg, err = addUser(phone_number, pd)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 404, "register user failed.", err)
		return
	}
	msg, err = addrole(phone_number)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 404, "register role failed.", err)
		return
	}
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, msg)
}

func addDomain(phone string) (string, error) {
	domainId := phone
	domainDesc := phone
	domainUpId := "devops_product"
	domainStatus := 0
	err := dbobj.Exec(sys_rdbms_036, domainId, domainDesc, domainUpId, domainStatus, "network", "network")
	if err != nil {
		logs.Error(err)
		return "add domain failed.", err
	}
	return "add domain successfully.", nil
}

func addOrg(phone string) (string, error) {
	org_unit_id := "888888"
	org_unit_desc := "顶层机构"
	up_org_id := "sys19890228"
	domain_id := phone
	start_date := "2016-01-01"
	end_date := "2116-01-01"
	create_user := "network"
	maintance_user := "network"
	org_status_id := 0
	id := phone + "_join_" + "888888"

	err := dbobj.Exec(sys_rdbms_043, org_unit_id, org_unit_desc, up_org_id, org_status_id,
		domain_id, start_date, end_date, create_user, maintance_user, id)
	if err != nil {
		logs.Error(err)
		return "add org failed.", err
	}
	return "add org successfully.", nil
}

func addUser(userid string, password string) (string, error) {

	stheme := `insert into sys_user_theme(user_id,theme_id) values(?,?)`

	userId := userid
	userDesc := userid

	userPasswd, err := utils.Encrypt(password)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	userStatus := 0
	userEmail := "net@163.com"
	userPhone := userid
	userOrgUnitId := userid + "_join_" + "888888"

	tx, err := dbobj.Begin()
	// insert user details
	//
	_, err = tx.Exec(sys_rdbms_018, userId, userDesc, "network", userEmail, userPhone, userOrgUnitId, "network")
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		return "", err
	}

	// insert user passwd
	//
	_, err = tx.Exec(sys_rdbms_019, userId, userPasswd, userStatus)
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		return "", err
	}

	// insert theme info
	//
	_, err = tx.Exec(stheme, userId, "1001")
	if err != nil {
		tx.Rollback()
		logs.Error(err.Error())
		return "", err
	}

	tx.Commit()
	return "注册成功", nil
}

func addrole(phone string) (string, error) {
	maintanceDate := time.Now().Format("2006-01-02")
	err := dbobj.Exec(sys_rdbms_024, "devops_product_join_networkadmin", phone, maintanceDate, "network")
	if err != nil {
		logs.Error(err)
		return "add role failed.", err
	}
	return "add role successfully.", nil
}

func init() {
	beego.Get("/plat/registerPage", httpRegisterPage)
	beego.Post("/plat/register/sendCode", getVerifyCode)
	beego.Post("/plat/register/verifyCode", verifyCode)
	beego.Post("/plat/register", noVerifyRegister)
}
