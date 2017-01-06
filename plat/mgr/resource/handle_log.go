package resource

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"

	"github.com/hzwy23/dbobj"
)

type HandleLog struct {
	Uuid           string
	User_id        string
	Hander_date    string
	Request_method string
	Res_url        string
	Ip_addr        string
	Return_id      string
	Error_msg      string
	Cnt            string
	route.RouteControl
}

type HandleLogPage struct {
	route.RouteControl
}

func (this *HandleLog) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	deftDb := dbobj.GetDefaultName()
	r.ParseForm()
	searchId := r.FormValue("SearchTrue")
	userId := r.FormValue("UserId")
	startDate := r.FormValue("StartDate")
	endDate := r.FormValue("EndDate")
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	cond := " 1=1 "
	sql := sqlText.PLATFORM_RESOURCE_LOGS1
	switch deftDb {
	case "mysql":
		if searchId != "" {
			if userId != "" {
				cond = " user_id = '" + userId + "'"
			}
			if startDate != "" {
				cond += " and str_to_date(hander_date,'%Y-%m-%d %H:%i:%s') >= str_to_date('" + startDate + "','%Y-%m-%d %H:%i:%s')"
			}
			if endDate != "" {
				cond += " and str_to_date(hander_date,'%Y-%m-%d %H:%i:%s') < str_to_date('" + endDate + "','%Y-%m-%d %H:%i:%s')"
			}
		}
	case "oracle":
		if searchId != "" {
			if userId != "" {
				cond = " user_id = '" + userId + "'"
			}
			if startDate != "" {
				cond += " and to_timestamp(hander_date,'YYYY-MM-DD HH24:MI:SS') >= to_timestamp('" + startDate + "','YYYY-MM-DD HH24:MI:SS')"
			}
			if endDate != "" {
				cond += " and to_timestamp(hander_date,'YYYY-MM-DD HH24:MI:SS') < to_timestamp('" + endDate + "','YYYY-MM-DD HH24:MI:SS')"
			}
		}
	}
	sql = strings.Replace(sql, "HZWY23", cond, -1)

	rows, err := dbobj.Query(sql, offset, limit+offset)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		return
	}
	//	var one HandleLog
	var rst []HandleLog
	dbobj.Scan(rows, &rst)
	//	for rows.Next() {
	//		err := rows.Scan(&one.Uuid,
	//			&one.User_id,
	//			&one.Hander_date,
	//			&one.Request_method,
	//			&one.Res_url,
	//			&one.Ip_addr,
	//			&one.Return_id,
	//			&one.Error_msg,
	//			&one.cnt)
	//		if err != nil {
	//			logs.Error(err)
	//			return
	//		}
	//		rst = append(rst, one)
	//	}
	total := dbobj.Count("select count(*) from sys_user_login_records")
	this.WritePage(w, total, rst)
}

//func (this *HandleLog) Post(w http.ResponseWriter, r *http.Request) {
//	if sys.Privilege.Access(w, r) == false {
//		w.WriteHeader(http.StatusForbidden)
//		return
//	}

//	r.ParseForm()
//	ijson := r.FormValue("JSON")
//	m := r.FormValue("_Method")
//	if m == "Delete" {
//		var ijs []HandleLog
//		err := json.Unmarshal([]byte(ijson), &ijs)
//		if err != nil {
//			logs.Error(err)
//			w.WriteHeader(http.StatusExpectationFailed)
//			w.Write([]byte("json数据解析失败"))
//			return
//		}
//		sql := ""
//		if "oracle" == dbobj.DefaultDB() {
//			sql = "delete from sys_user_login_records where uuid = :1"
//		} else if "mysql" == dbobj.DefaultDB() {
//			sql = "delete from sys_user_login_records where uuid = ?"
//		}
//		for _, val := range ijs {
//			err := dbobj.Default.Exec(sql, val.Uuid)
//			if err != nil {
//				logs.Error(err)
//				return
//			}
//		}
//		w.Write([]byte("success"))
//	} else {

//	}
//}

func (this *HandleLogPage) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hz, _ := template.ParseFiles("./views/platform/resource/handle_logs_page.tpl")
	hz.Execute(w, nil)
}
