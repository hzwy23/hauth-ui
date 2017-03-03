package models

import (
	"database/sql"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils/logs"
)

const (
	error_querydb  string = "can not found user info in system."
	error_maxerror string = "user was forbided, you have continued type password error 6 times."
	error_password string = "user's password error.please check your password"
)

type mSysUserSec struct {
	User_id                 string        `json:"user_id"`
	User_passwd             string        `json:"user_passwd"`
	User_status             sql.NullInt64 `json:"user_status"`
	User_continue_error_cnt sql.NullInt64
}

func updateContinueErrorCnt(cnt int64, user_id string) {
	dbobj.Exec("update sys_sec_user set continue_error_cnt = ? where user_id = ?", cnt, user_id)
}

func forbidUsers(user_id string) {
	dbobj.Exec("update sys_sec_user set status_id = 1 where user_id = ?", user_id)
}


// check user's passwd is right.
func BasicAuth(user_id, user_passwd string) (bool, int, int64, string) {
	var sec mSysUserSec
	err := dbobj.QueryRow(sys_rdbms_010, user_id).Scan(&sec.User_id, &sec.User_passwd, &sec.User_status, &sec.User_continue_error_cnt)
	if err != nil {
		return false, 402, 0, error_querydb
	}

	if sec.User_status.Int64 != 0 {
		return false, 406, sec.User_status.Int64, error_maxerror
	}

	if sec.User_continue_error_cnt.Int64 > 6 {
		forbidUsers(user_id)
		return false, 403, sec.User_continue_error_cnt.Int64, error_maxerror
	}

	if sec.User_id == user_id && sec.User_passwd == user_passwd {
		updateContinueErrorCnt(0, user_id)
		return true, 200, 0, ""
	} else {
		updateContinueErrorCnt(sec.User_continue_error_cnt.Int64+1, user_id)
		return false, 405, sec.User_continue_error_cnt.Int64 + 1, error_password
	}
}


// check the user wheather handle the domain
// return value :
// -1   : have no right to handle the domain
// 1    : can read the domain info
// 2    : can read and wirte the domain info
func CheckDomainRights(user_id string,domain_id string)int{
	var cnt = -1
	err := dbobj.QueryRow(sys_rdbms_001,domain_id,user_id).Scan(&cnt)
	if err!=nil{
		logs.Error(err)
		return -1
	}
	return cnt
}
