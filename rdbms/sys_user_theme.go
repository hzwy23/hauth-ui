package rdbms

import (
	"github.com/hzwy23/dbobj"
)

func GetSysUserTheme(id string) (themeId string, err error) {
	err = dbobj.QueryRow(sys_rdbms_001, id).Scan(&themeId)
	return
}

func DeleteSysUserTheme(id string) error {
	return dbobj.Exec(sys_rdbms_002, id)
}

func AddSysUserTheme(id string, themeid string) error {
	return dbobj.Exec(sys_rdbms_003, id, themeid)
}

func UpdateSysUserTheme(id string, themeid string) error {
	return dbobj.Exec(sys_rdbms_004, themeid, id)
}
