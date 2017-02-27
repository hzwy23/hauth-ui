package models

import (
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils/logs"
)

type OrgModel struct {
}

type SysOrgInfo struct {
	Org_unit_id     string `json:"org_id"`
	Org_unit_desc   string `json:"org_desc"`
	Up_org_id       string `json:"up_org_id"`
	Org_status_id   string `json:"status_id"`
	Org_status_desc string `json:"status_desc"`
	Domain_id       string `json:"domain_id"`
	Domain_desc     string `json:"domain_desc"`
	Start_date      string `json:"start_date"`
	End_date        string `json:"end_date"`
	Create_date     string `json:"create_date"`
	Maintance_date  string `json:"modify_date"`
	Create_user     string `json:"create_user"`
	Maintance_user  string `json:"modify_user"`
	Code_number     string `json:"code_number"`
	Org_dept        string `json:"org_dept,omitempty"`
}

//获取域下边所有机构号
func (OrgModel) Get(domain_id string) ([]SysOrgInfo, error) {
	var rst []SysOrgInfo
	rows, err := dbobj.Query(sys_rdbms_041, domain_id)
	if err != nil {
		return nil, err
	}

	err = dbobj.Scan(rows, &rst)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (OrgModel) Delete(mjs []SysOrgInfo) error {
	tx, _ := dbobj.Begin()
	for _, val := range mjs {
		_, err := tx.Exec(sys_rdbms_044, val.Org_unit_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (OrgModel) Update(org_unit_desc, up_org_id, org_status_id, start_date, end_date, maintance_user, org_unit_id string) error {
	return dbobj.Exec(sys_rdbms_069, org_unit_desc, up_org_id, org_status_id,
		start_date, end_date, maintance_user, org_unit_id)
}

func (OrgModel) Post(org_unit_id, org_unit_desc, up_org_id, org_status_id, domain_id, start_date, end_date, create_user, maintance_user, id string) error {
	return dbobj.Exec(sys_rdbms_043, org_unit_id, org_unit_desc, up_org_id, org_status_id,
		domain_id, start_date, end_date, create_user, maintance_user, id)
}

func (OrgModel) GetOrgByDomainId(org_id string, domain_id string, did string) ([]SysOrgInfo, error) {
	var rst []SysOrgInfo
	if did != domain_id {
		rows, err := dbobj.Query(sys_rdbms_061, domain_id, did)
		if err != nil {
			return nil, err
		}

		err = dbobj.Scan(rows, &rst)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err := dbobj.Query(sys_rdbms_060, org_id, domain_id, did)
		if err != nil {
			return nil, err
		}

		err = dbobj.Scan(rows, &rst)
		if err != nil {
			return nil, err
		}

	}
	return rst, nil
}
