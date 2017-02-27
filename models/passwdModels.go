package models

import (
	"github.com/hzwy23/dbobj"
)

type PasswdModels struct {
}

func (PasswdModels) UpdateMyPasswd(newPd, User_id, oriEn string) error {
	return dbobj.Exec(sys_rdbms_014, newPd, User_id, oriEn)
}

func (PasswdModels) UpdateUserPasswd(newPd, userid string) error {
	return dbobj.Exec(sys_rdbms_015, newPd, userid)
}
