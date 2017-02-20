package rpc

type rpcRoleInfo struct {
	Role_id     string `json"role_id"`
	Role_name   string `json:"role_name"`
	Domain_id   string `json:"domain_id"`
	Domain_name string `json:"domain_name"`
	Code_number string `json:"code_number"`
	Role_owner  string `json:"role_owner"`
}
