package models

import (
	"encoding/json"

	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils/logs"
)

type UserRolesModel struct {
	User_id string
}

type userRoleData struct{
	User_id string  `json:"user_id"`
	Role_id string  `json:"role_id"`
	Code_number string `json:"code_number"`
	Role_name string `json:"role_name"`
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

func (UserRolesModel)GetRolesByUser(user_id string)([]userRoleData,error){

	rows,err:=dbobj.Query(sys_rdbms_094,user_id)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	var rst []userRoleData
	err = dbobj.Scan(rows,&rst)
	return rst,err
}

func (UserRolesModel)GetOtherRoles(user_id string)([]userRoleData,error){
	rows,err:=dbobj.Query(sys_rdbms_095,user_id)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	var rst []userRoleData
	err = dbobj.Scan(rows,&rst)
	return rst,err
}

func (UserRolesModel)Auth(user_id ,ijs string)error{
	var rst []userRoleData
	err:=json.Unmarshal([]byte(ijs),&rst)
	if err!=nil{
		logs.Error(err)
		return err
	}
	tx,err:=dbobj.Begin()
	if err!=nil{
		logs.Error(err)
		return err
	}
	for _,val:=range rst {
		_,err:=tx.Exec(sys_rdbms_096,val.Role_id,val.User_id,user_id)
		if err!=nil{
			logs.Error(err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (UserRolesModel)Revoke(user_id string,role_id string)error{
	return dbobj.Exec(sys_rdbms_097,user_id,role_id)
}