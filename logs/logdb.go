package logs

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hzwy23/dbobj"
)

func LogToDB(r *http.Request, userid string, status bool, msg string) {

	curdate := time.Now().Format("2006-01-02 03:04:05")

	url := r.RequestURI

	method := r.Method

	ip := strings.Split(r.RemoteAddr, ":")[0]

	flag := "成功"

	if status == false {
		flag = "失败"
	}

	sql := "" //sqlText.DBLOG_SQL

	err := dbobj.Exec(sql, userid, curdate, method, url, ip, flag, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
