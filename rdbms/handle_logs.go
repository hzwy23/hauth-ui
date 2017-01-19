package rdbms

import (
	"html/template"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/hret"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/hauth/token/hjwt"
	"github.com/hzwy23/hauth/utils"
)

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

func GetHandleLogs(ctx *context.Context) {
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) order by handle_time desc limit ?,?`
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			)`
	hret.WriteBootstrapTableJson(ctx.ResponseWriter, dbobj.Count(cntsql, jclaim.Domain_id), rst)
}

func serachLogs(ctx *context.Context) {
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) and user_id = ? and handle_time >= str_to_date(?,'%Y-%m-%d')
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) and user_id = ? and handle_time >= str_to_date(?,'%Y-%m-%d')
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) and user_id = ? and handle_time >= str_to_date(?,'%Y-%m-%d')
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) and handle_time >= str_to_date(?,'%Y-%m-%d')
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) and handle_time >= str_to_date(?,'%Y-%m-%d')
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) and handle_time < str_to_date(?,'%Y-%m-%d')
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) and user_id = ?
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
			where exists (
					SELECT domain_id from sys_domain_info s
					where FIND_IN_SET(s.domain_id,getChildDomainList(?))
					and t.domain_id = s.domain_id  
			) order by user_id,handle_time desc`

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

func init() {
	beego.Get("/v1/auth/HandleLogsPage", func(ctx *context.Context) {
		hz, _ := template.ParseFiles("./views/platform/resource/handle_logs_page.tpl")
		hz.Execute(ctx.ResponseWriter, nil)
	})
	beego.Get("/v1/auth/handle/logs", GetHandleLogs)
	beego.Get("/v1/auth/handle/logs/search", serachLogs)
}
