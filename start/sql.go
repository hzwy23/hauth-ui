package start

import (
	"github.com/hzwy23/dbobj"
)

var (
	platform_resource_login_1 = ""
	platform_resource_login_2 = ""
	platform_resource_login_3 = ""
)

func init() {
	defdb := dbobj.GetDefaultName()
	if "mysql" == defdb {
		platform_resource_login_1 = `select t2.res_url from sys_user_theme t1 inner join sys_theme_value t2 on t1.theme_id = t2.theme_id where t1.user_id = ? and t2.res_id = 'index'`
		platform_resource_login_2 = `SELECT distinct domain_id FROM sys_user_info i inner join sys_org_info o on i.org_unit_id = o.id where user_id = ?`
		platform_resource_login_3 = `SELECT o.org_unit_id FROM sys_user_info i inner join sys_org_info o on i.org_unit_id = o.id where user_id = ?`
	}
}
