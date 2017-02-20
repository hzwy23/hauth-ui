package rpc

import (
	"net/http"
)

type rpcUserInfo struct {
	User_id          string `json:"user_id"`
	User_name        string `json:"user_name"`
	Org_unit_id      string `json:"org_unit_id"`
	Org_unit_desc    string `json:"org_unit_desc"`
	User_status_desc string `json:"user_status_desc"`
	Domain_name      string `json:"domain_name"`
	Domain_id        string `json:"domain_id"`
}

func GetSubUsers(r *http.Request) ([]rpcUserInfo, error) {
	return nil, nil
}
