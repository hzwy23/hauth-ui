package rpc

import (
	"net/http"
	"strconv"
)

type rpcOrgInfo struct {
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

func GetParentAndSubOrgs(r *http.Request) ([]rpcDomainInfo, error) {

	return nil, nil
}

func GetSubOrgs(r *http.Request) ([]rpcOrgInfo, error) {
	return nil, nil
}

func RpcFindOrg(d []rpcOrgInfo, id string) bool {
	for _, val := range d {
		if val.Org_unit_id == id {
			return true
		}
	}
	return false
}

func dorgtree(node []rpcOrgInfo, id string, d int, result *[]rpcOrgInfo) {
	var oneline rpcOrgInfo
	for _, val := range node {
		if val.Up_org_id == id {
			oneline.Code_number = val.Code_number
			oneline.Create_date = val.Create_date
			oneline.Create_user = val.Create_user
			oneline.Domain_desc = val.Domain_desc
			oneline.Domain_id = val.Domain_id
			oneline.End_date = val.End_date
			oneline.Org_dept = strconv.Itoa(d)
			oneline.Maintance_date = val.Maintance_date
			oneline.Maintance_user = val.Maintance_user
			*result = append(*result, oneline)
			dorgtree(node, val.Org_unit_id, d+1, result)
		}
	}
}
