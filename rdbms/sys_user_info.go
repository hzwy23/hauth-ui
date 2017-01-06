package rdbms

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
)

type mSysUserInfo struct {
	User_id          string `json:user_id`
	User_name        string `json:user_name`
	User_create_date string `json:user_create_date`
	User_owner       string `json:user_owner`
	User_email       string `json:user_email`
	User_phone       string `json:user_phone`
	Org_unit_id      string `json:org_unit_id`
	Org_unit_desc    string `json:org_unit_desc`
	Up_org_id        string `json:up_org_id`
	Org_status_id    string `json:org_status_id`
}

type SysUserInfo struct {
	beego.Controller
}

// if get user_id ,so return then details info of user_id
// no return all users info.
func (this *SysUserInfo) Get() {

	id := this.Ctx.Input.Param(":user_id")

	var ret mSysUserInfo

	if id == "" {
		id = this.Ctx.Request.FormValue("user_id")
	}

	if id == "" {
		logs.Error("user_id is empty. so didn't get user info.")
		hret.WriteHttpErrMsg(this.Ctx.ResponseWriter, hret.HttpErrMsg{
			Error_code:    http.StatusNotFound,
			Error_msg:     "user_id is empty. so didn't get user info.",
			Error_details: "v1.0",
		})
		return
	}
	row := dbobj.QueryRow(sys_rdbms_005, id)
	var (
		user_name        = ""
		user_create_date = ""
		user_owner       = ""
		user_email       = ""
		user_phone       = ""
		org_unit_id      = ""
		org_unit_desc    = ""
		up_org_id        = ""
		org_status_id    = ""
	)
	err := row.Scan(&ret.User_id,
		&user_name,
		&user_create_date,
		&user_owner,
		&user_email,
		&user_phone,
		&org_unit_id,
		&org_unit_desc,
		&up_org_id,
		&org_status_id)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(this.Ctx.ResponseWriter, http.StatusNotFound, "Get user info from database failed. error info is "+err.Error())
		return
	}
	ret.User_name = user_name
	ret.Org_status_id = org_status_id
	ret.Org_unit_desc = org_unit_desc
	ret.Org_unit_id = org_unit_id
	ret.Up_org_id = up_org_id
	ret.User_create_date = user_create_date
	ret.User_email = user_email
	ret.User_owner = user_owner
	ret.User_phone = user_phone
	hret.WriteHttpOkMsgs(this.Ctx.ResponseWriter, ret)
}

// update users info .
// accept argument slice .
func (this *SysUserInfo) Put() {
	this.Ctx.Request.ParseForm()
	info := this.Ctx.Request.FormValue("sys_user_info")
	var ret []mSysUserInfo
	err := json.Unmarshal([]byte(info), &ret)
	if err != nil {
		logs.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.Ctx.WriteString("Put all users info into database failed. error info is :" + err.Error())
		return
	}
	tx, err := dbobj.Begin()
	if err != nil {
		logs.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.Ctx.WriteString("Put all users info into database failed. error info is :" + err.Error())
		return
	}
	for _, val := range ret {
		_, err := tx.Exec(sys_rdbms_008, val.User_name, val.User_email, val.User_phone, val.Org_unit_id, val.User_id)
		if err != nil {
			logs.Error(err)
			this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
			this.Ctx.WriteString("Put all users info into database failed. error info is :" + err.Error())
			tx.Rollback()
			return
		}
	}
	tx.Commit()
}

// delete users info.
// accept argument slice.
func (this *SysUserInfo) Delete() {
	this.Ctx.Request.ParseForm()
	id := this.Ctx.Input.Param(":user_id")

	err := dbobj.Exec(sys_rdbms_007, id)
	if err != nil {
		logs.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.Ctx.WriteString("Post all users info into database failed. error info is :" + err.Error())
		return
	}
}

// insert users info.
// accept argument slice.
func (this *SysUserInfo) Post() {
	this.Ctx.Request.ParseForm()
	info := this.Ctx.Request.FormValue("sys_user_info")
	var ret []mSysUserInfo
	err := json.Unmarshal([]byte(info), &ret)
	if err != nil {
		logs.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.Ctx.WriteString("Post all users info into database failed. error info is :" + err.Error())
		return
	}
	tx, err := dbobj.Begin()
	if err != nil {
		logs.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.Ctx.WriteString("Post all users info into database failed. error info is :" + err.Error())
		return
	}
	for _, val := range ret {
		_, err := tx.Exec(sys_rdbms_008, val.User_id, val.User_name, val.User_owner, val.User_email, val.User_phone, val.Org_unit_id)
		if err != nil {
			logs.Error(err)
			this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
			this.Ctx.WriteString("Post all users info into database failed. error info is :" + err.Error())
			tx.Rollback()
			return
		}
	}
	tx.Commit()
}

func (this *SysUserInfo) GetAllUsers() {

	defer hret.HttpPanic()

	offset := this.Ctx.Input.Param(":offset")
	limit := this.Ctx.Input.Param(":limit")
	order := this.Ctx.Input.Param(":order")
	sort := this.Ctx.Input.Param(":mode")

	rows, err := dbobj.Query(sys_rdbms_006 + " order by " + order + " " + sort + " limit " + offset + "," + limit)
	defer rows.Close()

	if err != nil {
		logs.Error(err.Error())
		hret.WriteHttpErrMsg(this.Ctx.ResponseWriter, hret.HttpErrMsg{
			Error_code:    http.StatusNotFound,
			Error_msg:     "Get all users info from database failed. error info is :" + err.Error(),
			Error_details: "v1.0",
		})

		return
	}
	var ret []mSysUserInfo
	err = dbobj.Scan(rows, &ret)
	if err != nil {
		logs.Error(err)

		hret.WriteHttpErrMsg(this.Ctx.ResponseWriter, hret.HttpErrMsg{
			Error_code:    http.StatusNotFound,
			Error_msg:     "Get all users info from database failed. error info is" + err.Error(),
			Error_details: "v1.0",
		})

		return
	}
	hret.WriteHttpOkMsgs(this.Ctx.ResponseWriter, ret)
}

func (this *SysUserInfo) GetAllUsersNoOrder() {
	defer hret.HttpPanic()

	offset := this.Ctx.Input.Param(":offset")
	limit := this.Ctx.Input.Param(":limit")

	rows, err := dbobj.Query(sys_rdbms_006 + " limit " + offset + "," + limit)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsg(this.Ctx.ResponseWriter, hret.HttpErrMsg{
			Error_code:    http.StatusNotFound,
			Error_msg:     "Get all users info from database failed. error info is",
			Error_details: "v1.0",
		})
		return
	}
	var ret []mSysUserInfo
	err = dbobj.Scan(rows, &ret)
	if err != nil {
		logs.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.Ctx.WriteString("Get all users info from database failed. error info is :" + err.Error())
		return
	}
	hret.WriteHttpOkMsgs(this.Ctx.ResponseWriter, ret)
}
