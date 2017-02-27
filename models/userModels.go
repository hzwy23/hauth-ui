package models

import (
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils/logs"

)

type UserModel struct{

}

type userInfo struct {
	User_id             string `json:"user_id"`
	User_name           string `json:"user_name"`
	User_status_desc    string `json:"status_desc"`
	User_create_date    string `json:"create_date"`
	User_owner          string `json:"create_user"`
	User_email          string `json:"user_email"`
	User_phone          string `json:"user_phone"`
	Org_unit_id         string `json:"org_unit_id"`
	Org_unit_desc       string `json:"org_unit_desc"`
	Domain_id           string `json:"domain_id"`
	Domain_name         string `json:"domain_name"`
	User_maintance_date string `json:"modify_date"`
	User_maintance_user string `json:"modify_user"`
}

func (UserModel)GetDefault(domain_id ,offset, limit string)([]userInfo,error){

	row, err := dbobj.Query(sys_rdbms_017, domain_id ,offset, limit)
	defer row.Close()
	if err != nil {
		logs.Error(err)
		return nil,err
	}

	var rst []userInfo
	err = dbobj.Scan(row, &rst)
	return rst,err
}