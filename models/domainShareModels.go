package models

import (
	"github.com/hzwy23/dbobj"
	"github.com/astaxie/beego/logs"
)

type DomainShareModel struct{
}

type dsModel struct{
	Uuid                string `json:"uuid"`
	Target_domain_id    string `json:"target_domain_id"`
	Domain_name         string `json:"domain_name"`
	Authorization_level string `json:"auth_level"`
	Create_user         string `json:"create_user"`
	Create_date         string `json:"create_date"`
	Modify_user         string `json:"modify_user"`
	Modify_date         string `json:"modify_date"`
}

func (DomainShareModel)Get(domain_id string)([]dsModel,error){

	rows,err:=dbobj.Query(sys_rdbms_083,domain_id)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}

	var rst []dsModel

	err = dbobj.Scan(rows,&rst)

	return rst,err
}