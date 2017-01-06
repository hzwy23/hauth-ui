package sqlText

import (
	"os"
	"path"

	"github.com/hzwy23/hcloud/utils/config"

	"github.com/hzwy23/hcloud/logs"
)

var (
	PLATFORM_RESOURCE_USERINFO1        = `insert into sys_user_info (user_id,user_name,user_create_date,user_owner,user_email,user_phone,org_unit_id) values(:1,:2,sysdate,:3,:4,:5,:6)`
	PLATFORM_RESOURCE_USERINFO2        = `insert into sys_sec_user(user_id,user_passwd,status_id) values(:1,:2,:3)`
	PLATFORM_RESOURCE_USERCHECK        = `select user_id,user_passwd from SYS_SEC_USER where status_id='0' and user_id = :1`
	PLATFORM_RESOURCE_INDEX            = `select t2.res_url from sys_user_theme t1 inner join sys_theme_value t2 on t1.theme_id = t2.theme_id where t1.user_id = :1 and t2.res_id = 'backindex'`
	PLATFORM_RESOURCE_DEFAULTMENULIST1 = `select t.res_id,t.res_name,v.res_url,v.res_bg_color,v.res_class, v.res_img, v.group_id from (select * from sys_resource_info t start with t.res_up_id = :1 connect by prior t.res_id = t.res_up_id) t inner join SYS_THEME_VALUE v on t.res_id = v.res_id inner join sys_user_theme h on v.theme_id = h.theme_id inner join sys_role_user_relation r on r.user_id = h.user_id inner join sys_role_resource_relat e on r.role_id = e.role_id and e.res_id = t.res_id inner join sys_role_info f on f.role_id = e.role_id and f.role_status_id = '0' where v.res_type = :2 and h.user_id = :3 order by sort_id asc`
	PLATFORM_RESOURCE_ENTRY1           = `select t2.res_url from sys_user_theme t1 inner join SYS_THEME_VALUE t2 on t1.theme_id = t2.theme_id where t1.user_id = :1 and t2.res_id = :2 and t2.res_type = '0'`
	PLATFORM_RESOURCE_INDEXPAGE        = `select res_id,res_name from sys_resource_info where res_type = '0'`
	PLATFORM_RESOURCE_MENU1            = `select res_icon as menu_icon,res_id as menu_id,res_name as menu_name,res_url as menu_url,res_up_id as up_menu_id,res_attr as menu_attr ,res_img as menu_img,res_color as menu_color from sys_resource_info t where res_type = '1' start with  res_up_id = :1 connect by prior res_id = res_up_id`
	PLATFORM_RESOURCE_MENUMGR1         = `select res_icon as menu_icon ,res_id as menu_id,res_name as menu_name,res_url as menu_url,res_up_id as up_menu_id,res_attr_desc as menu_attr from sys_resource_info t inner join sys_resource_info_attr b on t.res_attr = b.res_attr  where res_type in ('0','1') start with  res_up_id = :1 connect by prior res_id = res_up_id`
	PLATFORM_RESOURCE_PASSWD           = `update SYS_SEC_USER set user_passwd = :1 where user_id = :2`
	PLATFORM_RESOURCE_PROJECT1         = `select project_id, project_name, domain_up_id, status_name, maintance_date, user_id from (select t.domain_id as project_id, t.domain_name as project_name, t.domain_up_id, s.domain_status_name  as status_name, t.domain_create_date  as maintance_date, t.domain_owner as user_id, row_number() over(order by t.domain_id) as rk from SYS_domain_info t inner join sys_domain_status_attr s  on t.domain_status_id = s.domain_status_id) where rk > :1  and rk <= :2`
	PLATFORM_RESOURCE_PROJECT2         = `delete from sys_domain_info where domain_id = :1`
	PLATFORM_RESOURCE_PROJECT3         = `insert into sys_domain_info(domain_id,domain_name,domain_up_id,domain_status_id,domain_create_date,domain_owner) values(:1,:2,:3,:4,:5,:6)`
	PLATFORM_RESOURCE_PROJECT4         = `update sys_domain_info set domain_name = :1,domain_up_id=:2,domain_status_id = :3 where domain_id = :4`
	PLATFORM_RESOURCE_RESROLE1         = `SELECT T.UUID, T.RES_ID, I.RES_NAME, I.RES_ATTR, I.RES_UP_ID FROM sys_role_resource_relat T INNER JOIN SYS_RESOURCE_INFO I ON T.RES_ID = I.RES_ID WHERE T.ROLE_ID = :1`
	PLATFORM_RESOURCE_RESROLE4         = `SELECT sys_guid(), I.RES_ID, I.RES_NAME, I.RES_ATTR, I.RES_UP_ID FROM SYS_RESOURCE_INFO I left JOIN sys_role_resource_relat T ON T.RES_ID = I.RES_ID and T.ROLE_ID = :1 where t.rowid is null`
	PLATFORM_RESOURCE_RESROLE5         = `SELECT sys_guid(), I.RES_ID, I.RES_NAME, I.RES_ATTR, I.RES_UP_ID FROM SYS_RESOURCE_INFO I left JOIN sys_role_resource_relat T ON T.RES_ID = I.RES_ID and T.ROLE_ID = :1 where t.rowid is not null`
	PLATFORM_RESOURCE_RESROLE2         = `delete from sys_role_resource_relat where uuid = :1`
	PLATFORM_RESOURCE_RESROLE3         = `insert into sys_role_resource_relat(uuid,role_id,res_id) values(sys_guid(),:1,:2)`
	PLATFORM_RESOURCE_RESINFO1         = `select RES_ID,RES_NAME,RES_ATTR,RES_attr_DESC,RES_URL,RES_UP_ID,RES_TYPE,RES_TYPE_DESC from (SELECT T.RES_ID,T.RES_NAME,T.RES_ATTR, A.RES_attr_DESC,T.RES_URL,T.RES_UP_ID,T.res_type,R.RES_TYPE_DESC,row_number() over(order by t.res_id) as rk FROM sys_resource_info T INNER JOIN sys_resource_info_attr A ON T.RES_ATTR = A.RES_ATTR INNER JOIN SYS_RESOURCE_TYPE_ATTR R ON T.RES_TYPE = R.RES_TYPE ) where rk >:1 and rk <= :2`
	PLATFORM_RESOURCE_RESINFO2         = `delete from sys_resource_info where res_id = :1`
	PLATFORM_RESOURCE_RESINFO3         = `insert into sys_resource_info(res_id,res_name,res_attr,res_url,res_up_id,res_type) values(:1,:2,:3,:4,:5,:6)`
	PLATFORM_RESOURCE_RESINFO4         = `update sys_resource_info set res_name=:1,res_attr=:2,res_url=:3,res_up_id =:4,res_icon=:5,res_type=:6 where res_id = :7`
	PLATFORM_RESOURCE_ROLEDOMAIN1      = `SELECT T.UUID,T.ROLE_ID, t.domain_id,i.domain_name,i.domain_up_id FROM sys_role_domain_relation T INNER JOIN SYS_DOMAIN_INFO I ON T.DOMAIN_ID = I.DOMAIN_ID where t.role_id = :1`
	PLATFORM_RESOURCE_ROLEDOMAIN2      = `delete from sys_role_domain_relation where uuid = :1`
	PLATFORM_RESOURCE_ROLEDOMAIN3      = `insert into sys_role_domain_relation(uuid,role_id,domain_id) values(:1,:2,:3)`
	PLATFORM_RESOURCE_ROLEUSER1        = `select UUID, User_id, Role_id, Role_name,maintance_date,maintance_user from (SELECT T.UUID,T.User_id,t.Role_id,i.Role_name,t.maintance_date,t.maintance_user,row_number() over(order by t.uuid) as rk FROM sys_role_user_relation T INNER JOIN SYS_ROLE_INFO I ON T.ROLE_ID = I.ROLE_ID where t.user_id = :1) where rk > :2 and rk <= :3`
	PLATFORM_RESOURCE_ROLEUSER2        = `delete from  sys_role_user_relation where uuid = :1`
	PLATFORM_RESOURCE_ROLEUSER3        = `insert into sys_role_user_relation(role_id,user_id,maintance_date,maintance_user) values(:1,:2,to_date(:3,'YYYY-MM-DD'),:4)`
	PLATFORM_RESOURCE_ROLEINFO1        = `select role_id,role_name,role_owner,role_create_date,role_status_desc,role_status_id from (select t.role_id,t.role_name,t.role_owner,t.role_create_date,a.role_status_desc,a.role_status_id,row_number() over(order by t.role_id) as rk from sys_role_info t inner join sys_role_status_attr a on t.role_status_id = a.role_status_id) where rk > :1 and rk <= :2`
	PLATFORM_RESOURCE_ROLEINFO2        = `insert into sys_role_info(role_id,role_name,role_owner,role_create_date,role_status_id) values(:1,:2,:3,to_date(:4,'YYYY-MM-DD'),:5)`
	PLATFORM_RESOURCE_ROLEINFO3        = `delete from  sys_role_info where role_id = :1`
	PLATFORM_RESOURCE_ROLEINFO4        = `update sys_role_info t set t.role_name = :1 ,t.role_status_id = :2 where t.role_id = :3`
	PLATFORM_RESOURCE_USER1            = `select  user_id,user_name,status_desc,user_create_date,user_owner,user_email,user_phone,org_unit_desc from (select t.user_id,t.user_name,a.status_desc,t.user_create_date,t.user_owner,t.user_email,t.user_phone,i.org_unit_desc,row_number() over(order by t.user_id) as rk from sys_user_info t  inner join sys_sec_user u on t.user_id = u.user_id inner join sys_user_status_attr a on u.status_id = a.status_id inner join sys_org_info i on i.org_unit_id = t.org_unit_id) where rk > :1 and rk <= :2`
	PLATFORM_RESOURCE_USER2            = `delete from sys_user_info where user_id = :1`
	PLATFORM_RESOURCE_USER3            = `delete from sys_sec_user where user_id = :1`
	PLATFORM_RESOURCE_USER4            = `update SYS_USER_INFO t set t.user_name = :1 , t.user_phone = :2, t.user_email = :3 , t.org_unit_id = :4 where t.user_id = :5`
	PLATFORM_RESOURCE_USER6            = `SELECT count(*) FROM sys_user_info t INNER JOIN sys_sec_user u ON t.user_id = u.user_id  INNER JOIN sys_user_status_attr a ON u.status_id = a.status_id INNER JOIN sys_org_info i ON i.org_unit_id = t.org_unit_id`
	PLATFORM_RESOURCE_USER5            = `update sys_sec_user t set t.status_id = :1 where t.user_id = :2 `
	PLATFORM_RESOURCE_USER7            = `delete from sys_user_theme where user_id = :1`
	PLATFORM_RESOURCE_LOGS1            = `select uuid,user_id ,hander_date,request_method ,res_url,ip_addr,return_id,error_msg from ( select uuid ,user_id  ,hander_date ,request_method ,res_url ,nvl(ip_addr,' ') as ip_addr ,return_id,error_msg ,row_number() over(order by hander_date desc) as rk  from sys_user_login_records where HZWY23) where rk > :1 and rk <= :2`
	DBLOG_SQL                          = `insert into sys_user_login_records(uuid,user_id,hander_date,request_method,res_url,ip_addr,return_id,error_msg) values(sys_guid(),:1,:2,:3,:4,:5,:6,:7)`
	PLATFORM_RESOURCE_LOGIN1           = `select t2.res_url from sys_user_theme t1 inner join sys_theme_value t2 on t1.theme_id = t2.theme_id where t1.user_id = ? and t2.res_id = 'index'`
	PLATFORM_RESOURCE_LOGIN2           = `select t2.res_url from sys_user_theme t1 inner join sys_theme_value t2 on t1.theme_id = t2.theme_id  where t1.user_id = ? and t2.res_id = 'index'`
)

func init() {

	filedir := path.Join(os.Getenv("HBIGDATA_HOME"), "conf", "system.properties")

	logs.Debug("init system. read system config file. dir is :", filedir)
	red, err := config.GetConfig(filedir)
	if err != nil {
		logs.Error("cant not read ./conf/system.properties.please check this file.")
	}
	defdb, _ := red.Get("DB.type")

	if defdb == "oracle" {
	} else if defdb == "mysql" {
		DBLOG_SQL = `insert into sys_user_login_records(uuid,user_id,hander_date,request_method,res_url,ip_addr,return_id,error_msg) values(uuid(),?,?,?,?,?,?,?)`
		PLATFORM_RESOURCE_LOGIN1 = `select t2.res_url  from sys_user_theme t1 inner join sys_theme_value t2 on t1.theme_id = t2.theme_id where t1.user_id = ? and t2.res_id = 'index'`
		PLATFORM_RESOURCE_LOGIN2 = `select t2.res_url  from sys_user_theme t1 inner join sys_theme_value t2 on t1.theme_id = t2.theme_id  where t1.user_id = ? and t2.res_id = 'index'`
		PLATFORM_RESOURCE_USERINFO1 = `insert into sys_user_info (user_id,user_name,user_create_date,user_owner,user_email,user_phone,org_unit_id) values(?,?,now(),?,?,?,?)`
		PLATFORM_RESOURCE_USERINFO2 = `insert into sys_sec_user(user_id,user_passwd,status_id) values(?,?,?)`
		PLATFORM_RESOURCE_USERCHECK = `select user_id,user_passwd from SYS_SEC_USER where status_id='0' and user_id = ?`
		PLATFORM_RESOURCE_INDEX = `select t2.res_url from sys_user_theme t1 inner join sys_theme_value t2 on t1.theme_id = t2.theme_id where t1.user_id = ? and t2.res_id = 'backindex'`
		PLATFORM_RESOURCE_DEFAULTMENULIST1 = `select distinct t.res_id,t.res_name,v.res_url,v.res_bg_color,v.res_class, v.res_img, v.group_id from (select * from sys_resource_info t where find_in_set(res_up_id,getChildList(?))) t inner join SYS_THEME_VALUE v on t.res_id = v.res_id inner join sys_user_theme h on v.theme_id = h.theme_id inner join sys_role_user_relation r on r.user_id = h.user_id inner join sys_role_resource_relat e on r.role_id = e.role_id and e.res_id = t.res_id inner join sys_role_info f on f.role_id = e.role_id and f.role_status_id = '0' where v.res_type = ? and h.user_id = ? order by sort_id asc`
		PLATFORM_RESOURCE_ENTRY1 = `select distinct t2.res_url from sys_user_theme t1 inner join SYS_THEME_VALUE t2 on t1.theme_id = t2.theme_id where t1.user_id = ? and t2.res_id = ? and t2.res_type = '0'`
		PLATFORM_RESOURCE_INDEXPAGE = `select res_id,res_name from sys_resource_info where res_type = '0'`
		PLATFORM_RESOURCE_MENU1 = `select res_icon as menu_icon,res_id as menu_id,res_name as menu_name,res_url as menu_url,res_up_id as up_menu_id,res_attr as menu_attr ,res_img as menu_img,res_color as menu_color from sys_resource_info t where res_type = '1' start with  res_up_id = ? connect by prior res_id = res_up_id`
		PLATFORM_RESOURCE_MENUMGR1 = `select res_icon as menu_icon ,res_id as menu_id,res_name as menu_name,res_url as menu_url,res_up_id as up_menu_id,res_attr_desc as menu_attr from sys_resource_info t inner join sys_resource_info_attr b on t.res_attr = b.res_attr  where res_type in ('0','1') start with  res_up_id = ? connect by prior res_id = res_up_id`
		PLATFORM_RESOURCE_PASSWD = `update SYS_SEC_USER set user_passwd = ? where user_id = ?`
		PLATFORM_RESOURCE_PROJECT1 = `select t.domain_id as project_id, t.domain_name as project_name, t.domain_up_id, s.domain_status_name  as status_name, t.domain_create_date  as maintance_date, t.domain_owner as user_id from SYS_domain_info t inner join sys_domain_status_attr s  on t.domain_status_id = s.domain_status_id limit ?,?`
		PLATFORM_RESOURCE_PROJECT2 = `delete from sys_domain_info where domain_id = ?`
		PLATFORM_RESOURCE_PROJECT3 = `insert into sys_domain_info(domain_id,domain_name,domain_up_id,domain_status_id,domain_create_date,domain_owner) values(?,?,?,?,?,?)`
		PLATFORM_RESOURCE_PROJECT4 = `update sys_domain_info set domain_name = ?, domain_up_id=?, domain_status_id = ? where domain_id = ?`
		PLATFORM_RESOURCE_RESROLE1 = `SELECT T.UUID, T.RES_ID, I.RES_NAME, I.RES_ATTR, I.RES_UP_ID FROM sys_role_resource_relat T INNER JOIN SYS_RESOURCE_INFO I ON T.RES_ID = I.RES_ID WHERE T.ROLE_ID = ?`
		PLATFORM_RESOURCE_RESROLE4 = `SELECT uuid(), I.RES_ID, I.RES_NAME, I.RES_ATTR, I.RES_UP_ID FROM SYS_RESOURCE_INFO I left JOIN sys_role_resource_relat T ON T.RES_ID = I.RES_ID and T.ROLE_ID = ? where t.role_id is null`
		PLATFORM_RESOURCE_RESROLE5 = `SELECT uuid(), I.RES_ID, I.RES_NAME, I.RES_ATTR, I.RES_UP_ID FROM SYS_RESOURCE_INFO I left JOIN sys_role_resource_relat T ON T.RES_ID = I.RES_ID and T.ROLE_ID = ? where t.role_id is not null`
		PLATFORM_RESOURCE_RESROLE2 = `delete from  sys_role_resource_relat where uuid = ?`
		PLATFORM_RESOURCE_RESROLE3 = `insert into sys_role_resource_relat(uuid,role_id,res_id) values(uuid(),?,?)`
		PLATFORM_RESOURCE_RESINFO1 = `SELECT T.RES_ID,T.RES_NAME,T.RES_ATTR, A.RES_attr_DESC,T.RES_URL,T.RES_UP_ID,T.res_type,R.RES_TYPE_DESC FROM sys_resource_info T INNER JOIN sys_resource_info_attr A ON T.RES_ATTR = A.RES_ATTR INNER JOIN SYS_RESOURCE_TYPE_ATTR R ON T.RES_TYPE = R.RES_TYPE limit ?,?`
		PLATFORM_RESOURCE_RESINFO2 = `delete from sys_resource_info where res_id = ?`
		PLATFORM_RESOURCE_RESINFO3 = `insert into sys_resource_info(res_id,res_name,res_attr,res_url,res_up_id,res_type) values(?,?,?,?,?,?)`
		PLATFORM_RESOURCE_RESINFO4 = `update sys_resource_info set res_name=?,res_attr=?,res_url=?,res_up_id =?,res_icon=?,res_type=? where res_id = ?`
		PLATFORM_RESOURCE_ROLEDOMAIN1 = `SELECT T.UUID,T.ROLE_ID, t.domain_id,i.domain_name,i.domain_up_id FROM sys_role_domain_relation T INNER JOIN SYS_DOMAIN_INFO I ON T.DOMAIN_ID = I.DOMAIN_ID where t.role_id = ?`
		PLATFORM_RESOURCE_ROLEDOMAIN2 = `delete from  sys_role_domain_relation where uuid = ?`
		PLATFORM_RESOURCE_ROLEDOMAIN3 = `insert into sys_role_domain_relation(uuid,role_id,domain_id) values(?,?,?)`
		PLATFORM_RESOURCE_ROLEUSER1 = `SELECT T.UUID,T.User_id,t.Role_id,i.Role_name,t.maintance_date,t.maintance_user FROM sys_role_user_relation T INNER JOIN SYS_ROLE_INFO I ON T.ROLE_ID = I.ROLE_ID where t.user_id = ? limit ?,?`
		PLATFORM_RESOURCE_ROLEUSER2 = `delete from  sys_role_user_relation where uuid = ?`
		PLATFORM_RESOURCE_ROLEUSER3 = `insert into sys_role_user_relation(uuid,role_id,user_id,maintance_date,maintance_user) values(uuid(),?,?,str_to_date(?,'%Y-%m-%d'),?)`
		PLATFORM_RESOURCE_ROLEINFO1 = `select t.role_id,t.role_name,t.role_owner,t.role_create_date,a.role_status_desc,a.role_status_id from sys_role_info t inner join sys_role_status_attr a on t.role_status_id = a.role_status_id limit ?,?`
		PLATFORM_RESOURCE_ROLEINFO2 = `insert into sys_role_info(role_id,role_name,role_owner,role_create_date,role_status_id) values(?,?,?,str_to_date(?,'%Y-%m-%d'),?)`
		PLATFORM_RESOURCE_ROLEINFO3 = `delete from  sys_role_info where role_id = ?`
		PLATFORM_RESOURCE_ROLEINFO4 = `update sys_role_info t set t.role_name = ? ,t.role_status_id = ? where t.role_id = ?`
		PLATFORM_RESOURCE_USER1 = `
									select t.user_id,t.user_name,a.status_desc,t.user_create_date,
									t.user_owner,t.user_email,t.user_phone,i.org_unit_desc,
									dr.domain_id,di.domain_name 
									from sys_user_info t  
									inner join sys_sec_user u on t.user_id = u.user_id 
									inner join sys_user_status_attr a on u.status_id = a.status_id 
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where exists (
									    SELECT domain_id from sys_domain_info s
									    where FIND_IN_SET(s.domain_id,getChildDomainList(?)) 
									    and s.domain_status_id = '0'
									    and ( di.domain_id = s.domain_id or di.domain_id is null )
									)
									order by HZWSORTCOL HZWSORTAD limit ?,?`
		PLATFORM_RESOURCE_USER2 = `delete from sys_user_info where user_id = ?`
		PLATFORM_RESOURCE_USER3 = `delete from sys_sec_user where user_id = ?`
		PLATFORM_RESOURCE_USER4 = `update SYS_USER_INFO t set t.user_name = ?, t.user_phone = ?, t.user_email = ? , t.org_unit_id = ? where t.user_id = ?`
		PLATFORM_RESOURCE_USER5 = `update sys_sec_user t set t.status_id = ? where t.user_id = ?`
		PLATFORM_RESOURCE_USER7 = `delete from sys_user_theme where user_id = ?`
		PLATFORM_RESOURCE_LOGS1 = `select uuid ,user_id  ,hander_date ,request_method ,res_url ,ifnull(ip_addr,' ') as ip_addr ,return_id,error_msg from sys_user_login_records where HZWY23 order by hander_date desc limit ?,?`
	}
}
