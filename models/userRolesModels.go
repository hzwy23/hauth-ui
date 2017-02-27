package models

import (
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils/logs"
)

type UserRolesModel struct {
	User_id string
}

func (UserRolesModel) CleanRoles(rst []UserRolesModel) error {
	tx, _ := dbobj.Begin()
	for _, val := range rst {
		_, err := tx.Exec(sys_rdbms_045, val.User_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
