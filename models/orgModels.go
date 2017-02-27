package models

import (
	"github.com/hzwy23/dbobj"
	"github.com/hzwy23/hauth/utils/logs"
)

type OrgModel struct {
}

type SysOrgInfo struct {
	Org_unit_id     string `json:"Org_unit_id"`
	Org_unit_desc   string `json:"Org_unit_desc"`
	Up_org_id       string `json:"Up_org_id"`
	Org_status_id   string `json:"Org_status_id"`
	Org_status_desc string `json:"Org_status_desc"`
	Domain_id       string `json:"Domain_id"`
	Domain_desc     string `json:"Domain_desc"`
	Start_date      string `json:"Start_date"`
	End_date        string `json:"End_date"`
	Create_date     string `json:"Create_date"`
	Maintance_date  string `json:"Maintance_date"`
	Create_user     string `json:"Create_user"`
	Maintance_user  string `json:"Maintance_user"`
	Code_number     string `json:"Code_number"`
	Org_dept        string `json:"Org_dept,omitempty"`
}

func (OrgModel) Get(domain_id, org_id, offset, limit string) ([]SysOrgInfo, error) {
	var rst []SysOrgInfo
	rows, err := dbobj.Query(sys_rdbms_041, domain_id, org_id, offset, limit)
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
