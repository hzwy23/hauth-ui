package models

import (
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils/logs"
	"encoding/json"
	"errors"
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
	User_status_id     string  `json:"status_cd"`
}

func (UserModel)GetDefault(domain_id string)([]userInfo,error){

	row, err := dbobj.Query(sys_rdbms_017, domain_id)
	defer row.Close()
	if err != nil {
		logs.Error(err)
		return nil,err
	}

	var rst []userInfo
	err = dbobj.Scan(row, &rst)
	return rst,err
}

func (UserModel)Post(userId,userPasswd,userDesc,userStatus,id,userEmail,userPhone,userOrgUnitId string)error{
	tx, err := dbobj.Begin()
	// insert user details
	//
	_, err = tx.Exec(sys_rdbms_018, userId, userDesc, id, userEmail, userPhone, userOrgUnitId, id)
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		return err
	}

	// insert user passwd
	//
	_, err = tx.Exec(sys_rdbms_019, userId, userPasswd, userStatus)
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		return err
	}

	// insert theme info
	//

	stheme := `insert into sys_user_theme(user_id,theme_id) values(?,?)`

	_, err = tx.Exec(stheme, userId, "1001")
	if err != nil {
		tx.Rollback()
		logs.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (UserModel)Delete(ijs []byte,user_id string)error{
	var js []userInfo
	err := json.Unmarshal(ijs, &js)
	if err != nil {
		logs.Error(err)
		return err
	}

	tx, _ := dbobj.Begin()
	for _, val := range js {
		//判断用户是否在线
		//如果在线,则不允许删除用户
		if val.User_id == "admin" {
			tx.Rollback()
			return errors.New("admin是系统内置管理员，无法被删除")
		}

		// check user
		// can't delete yourself
		if user_id == val.User_id {
			tx.Rollback()
			return errors.New("禁止将自己删除。")
		}

		_, err := tx.Exec(sys_rdbms_007, val.User_id)
		if err != nil {
			tx.Rollback()
			logs.Error(err)
			return err
		}
	}
	return tx.Commit()
}

func (this UserModel)Search(org_id string,status_id string,domain_id string)([]userInfo,error){
	var rst []userInfo
	var err error
	// 如果机构号为空
	// 直接查询指定域中所有的用户
	if org_id == "" {
		rst,err = this.GetDefault(domain_id)
		if err!=nil{
			logs.Error(err)
			return nil,err
		}
	}else{
		rows,err:=dbobj.Query(sys_rdbms_090,org_id)
		if err!=nil{
			logs.Error(err)
			return nil,err
		}

		err = dbobj.Scan(rows,&rst)
		if err!=nil{
			logs.Error(err)
			return nil,err
		}
	}


	if status_id == ""{
		return rst,nil
	} else if status_id == "0"{
		var ret []userInfo
		for _,val:=range rst{
			if val.User_status_id == "0"{
				ret = append(ret,val)
			}
		}
		return ret,nil
	} else {
		var ret []userInfo
		for _,val:=range rst{
			if val.User_status_id == "1"{
				ret = append(ret,val)
			}
		}
		return ret,nil
	}
}