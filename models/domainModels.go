package models

import (
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils/logs"
)

type ProjectMgr struct {
	Project_id            string `json:"domain_id"`
	Project_name          string `json:"domain_desc"`
	Project_status        string `json:"domain_status"`
	Maintance_date        string `json:"maintance_date"`
	User_id               string `json:"create_user_id"`
	Domain_maintance_date string `json:"domain_modify_date"`
	Domain_maintance_user string `json:"domain_modify_user"`
	Domain_dept           string `json:"domain_dept"`
	Domain_up_id          string `json:"domain_up_id"`
}

func (ProjectMgr) Get(offset, limit string) ([]ProjectMgr, error) {
	rows, err := dbobj.Query(sys_rdbms_034, offset, limit)
	defer rows.Close()
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		return nil, err
	}

	//	var oneLine ProjectMgr
	var rst []ProjectMgr
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		return nil, err
	}
	return rst, nil
}

func (ProjectMgr) Post(domain_id, domain_desc, domain_status, user_id string) error {
	return dbobj.Exec(sys_rdbms_036, domain_id, domain_desc, domain_status, user_id, user_id)
}

func (ProjectMgr) Delete(js []ProjectMgr) error {
	tx, err := dbobj.Begin()
	if err != nil {
		logs.Error(err)
		return err
	}
	for _, val := range js {
		_, err := tx.Exec(sys_rdbms_037, val.Project_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
func (ProjectMgr) Update(domainDesc, domainStatus, user_id, domainId string) error {
	return dbobj.Exec(sys_rdbms_038, domainDesc, domainStatus, user_id, domainId)
}

func (ProjectMgr) GetDomainInfoByUser(domain_id string) ([]ProjectMgr, error) {
	rows, err := dbobj.Query(sys_rdbms_035, domain_id)
	defer rows.Close()
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		return nil, err
	}

	//	var oneLine ProjectMgr
	var rst []ProjectMgr
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error("query data error.", dbobj.GetErrorMsg(err))
		return nil, err
	}
	return rst, nil
}

func (ProjectMgr) GetDomainInfoByUpId(domainid string) ([]ProjectMgr, error) {
	rows, err := dbobj.Query(sys_rdbms_067, domainid)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	var rst []ProjectMgr
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return rst, nil
}
