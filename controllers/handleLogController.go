package controllers

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/astaxie/beego/context"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils"
	"github.com/hzwy23/hauth/utils/hret"
	"github.com/hzwy23/hauth/utils/logs"
	"github.com/hzwy23/hauth/utils/token/hjwt"
	"github.com/tealeg/xlsx"
)

type testdate struct {
	Id   string
	Name string
	Age  string
}

func genXlsxData() *bytes.Buffer {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	//	var cell *xlsx.Cell
	var err error

	var td = testdate{
		Id:   "hi",
		Name: "hello",
		Age:  "f10"}

	var arr = []string{"1", "2", "3", "4"}

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	row.WriteSlice(&arr, -1)

	row2 := sheet.AddRow()
	row2.WriteStruct(&td, -1)
	//cell = row.AddCell()
	//cell.Value = "I am a cell!"
	var tmp = make([]byte, 1)
	var buf = bytes.NewBuffer(tmp)
	err = file.Write(buf)

	if err != nil {
		fmt.Printf(err.Error())
		return buf
	}
	return buf
}

type HandleLogsController struct {
}

type handleLogs struct {
	Uuid        string `json:"uuid"`
	User_id     string `json:"user_id"`
	Handle_time string `json:"handle_time"`
	Client_ip   string `json:"client_ip"`
	Status_code string `json:"status_code"`
	Method      string `json:"method"`
	Url         string `json:"url"`
	Data        string `json:"data"`
}

var HandleLogsCtl = new(HandleLogsController)

func (HandleLogsController) GetHandleLogPage(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/hauth/handle_logs_page.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func (HandleLogsController) Download(ctx *context.Context) {
	//ctx.ResponseWriter.Header().Set("Content-Type", "application/vnd.ms-excel")
	var buf = genXlsxData().Bytes()
	fd, _ := os.Create("testxlsx.xlsx")
	fd.Write(buf)
	fd.Sync()
	fd.Close()
	ctx.ResponseWriter.Write(buf)
}

func (HandleLogsController) GetHandleLogs(ctx *context.Context) {
	ctx.Request.ParseForm()
	offset := ctx.Request.FormValue("offset")
	limit := ctx.Request.FormValue("limit")
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? order by handle_time desc limit ?,?`
	var rst []handleLogs
	rows, err := dbobj.Query(sql, jclaim.Domain_id, offset, limit)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
		return
	}
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
		return
	}
	cntsql := `select count(*) from sys_handle_logs t
			where t.domain_id = ?`
	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(cntsql, jclaim.Domain_id), rst)
}

func (HandleLogsController) SerachLogs(ctx *context.Context) {
	ctx.Request.ParseForm()
	userid := ctx.Request.FormValue("UserId")
	start := ctx.Request.FormValue("StartDate")
	end := ctx.Request.FormValue("EndDate")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "No Auth")
		return
	}
	var rst []handleLogs
	if userid != "" && utils.ValidDate(start) && utils.ValidDate(end) {
		sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? and user_id = ? and handle_time >= str_to_date(?,'%Y-%m-%d')
			and handle_time < str_to_date(?,'%Y-%m-%d')
			order by handle_time desc`

		rows, err := dbobj.Query(sql, jclaim.Domain_id, userid, start, end)
		defer rows.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
		err = dbobj.Scan(rows, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
	} else if userid != "" && utils.ValidDate(start) {
		sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? and user_id = ? and handle_time >= str_to_date(?,'%Y-%m-%d')
			order by handle_time desc`

		rows, err := dbobj.Query(sql, jclaim.Domain_id, userid, start)
		defer rows.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
		err = dbobj.Scan(rows, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
	} else if userid != "" && utils.ValidDate(end) {
		sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? and user_id = ? and handle_time >= str_to_date(?,'%Y-%m-%d')
			and handle_time < str_to_date(?,'%Y-%m-%d')
			order by handle_time desc`

		rows, err := dbobj.Query(sql, jclaim.Domain_id, userid, start, end)
		defer rows.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
		err = dbobj.Scan(rows, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
	} else if utils.ValidDate(start) && utils.ValidDate(end) {
		sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? and handle_time >= str_to_date(?,'%Y-%m-%d')
			and handle_time < str_to_date(?,'%Y-%m-%d')
			order by handle_time desc`

		rows, err := dbobj.Query(sql, jclaim.Domain_id, start, end)
		defer rows.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
		err = dbobj.Scan(rows, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
	} else if utils.ValidDate(start) {
		sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? and handle_time >= str_to_date(?,'%Y-%m-%d')
			order by handle_time desc`

		rows, err := dbobj.Query(sql, jclaim.Domain_id, start, end)
		defer rows.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
		err = dbobj.Scan(rows, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
	} else if utils.ValidDate(end) {
		sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? and handle_time < str_to_date(?,'%Y-%m-%d')
			order by handle_time desc`

		rows, err := dbobj.Query(sql, jclaim.Domain_id, start, end)
		defer rows.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
		err = dbobj.Scan(rows, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
	} else if userid != "" {
		sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? and user_id = ?
			order by handle_time desc`

		rows, err := dbobj.Query(sql, jclaim.Domain_id, userid)
		defer rows.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
		err = dbobj.Scan(rows, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
	} else {
		sql := `select uuid,user_id,handle_time,client_ip,status_code,method,url,data from sys_handle_logs t
			where t.domain_id = ? order by user_id,handle_time desc`

		rows, err := dbobj.Query(sql, jclaim.Domain_id)
		defer rows.Close()
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
		err = dbobj.Scan(rows, &rst)
		if err != nil {
			logs.Error(err)
			hret.WriteHttpErrMsgs(ctx.ResponseWriter, 310, "query failed.")
			return
		}
	}

	hret.WriteJson(ctx.ResponseWriter, rst)
}
