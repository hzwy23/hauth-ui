package rpc

import (
	"github.com/hzwy23/dbobj"
)

func init() {
	defdb := dbobj.GetDefaultName()
	if "mysql" == defdb {
		sys_rpc_001 = `select t.domain_id as project_id, t.domain_name as project_name, t.domain_up_id, s.domain_status_name  as status_name, t.domain_create_date  as maintance_date, t.domain_owner as user_id,t.domain_maintance_date,t.domain_maintance_user
						from SYS_domain_info t inner join sys_domain_status_attr s  on t.domain_status_id = s.domain_status_id 
						where exists (
							SELECT domain_id from sys_domain_info s
							where FIND_IN_SET(s.domain_id,getChildDomainList(?))
							and t.domain_id = s.domain_id 
							and t.domain_status_id = '0'
						)`
		sys_rpc_002 = `select t.domain_id as project_id, t.domain_name as project_name, t.domain_up_id, s.domain_status_name  as status_name, t.domain_create_date  as maintance_date, t.domain_owner as user_id,t.domain_maintance_date,t.domain_maintance_user
						from SYS_domain_info t inner join sys_domain_status_attr s  on t.domain_status_id = s.domain_status_id 
						where exists (
							SELECT domain_id from sys_domain_info s
							where FIND_IN_SET(s.domain_id,getParentAndChindDomainList(?))
							and t.domain_id = s.domain_id
							and t.domain_status_id = '0'
						)`
		sys_rpc_003 = `select 1 from sys_role_user_relation t inner join sys_role_resource_relat r on t.role_id = r.role_id inner join sys_user_theme e on t.user_id = e.user_id inner join sys_theme_value v on r.res_id = v.res_id and e.theme_id = v.theme_id where t.user_id = ? and v.res_url = ?`
		sys_rpc_004 = `select 1 from sys_role_user_relation t inner join sys_role_resource_relat r on t.role_id = r.role_id where t.user_id = ? and r.res_id = ?`

	}
}
