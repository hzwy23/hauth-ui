package rdbms

import (
	"github.com/hzwy23/dbobj"
)

func init() {
	defdb := dbobj.GetDefaultName()
	if "mysql" == defdb {
		sys_rdbms_001 = `SELECT theme_id FROM sys_user_theme where user_id = ?`
		sys_rdbms_002 = `delete from sys_user_theme where user_id = ?`
		sys_rdbms_003 = `insert into sys_user_theme values(?,?)`
		sys_rdbms_004 = `update sys_user_theme set theme_id = ? where user_id = ?`
		sys_rdbms_005 = `SELECT i.user_id,i.user_name,i.user_create_date,i.user_owner,i.user_email,i.user_phone,i.org_unit_id,f.org_unit_desc,up_org_id,f.org_status_id FROM sys_user_info i left join sys_org_info f on i.org_unit_id = f.org_unit_id where user_id = ?`
		sys_rdbms_006 = `SELECT i.user_id,i.user_name,i.user_create_date,i.user_owner,i.user_email,i.user_phone,i.org_unit_id,f.org_unit_desc,up_org_id,f.org_status_id FROM sys_user_info i left join sys_org_info f on i.org_unit_id = f.org_unit_id `
		sys_rdbms_007 = `delete from sys_user_info where user_id = ?`
		sys_rdbms_008 = `insert into sys_user_info(user_id,user_name,user_create_date,user_owner,user_email,user_phone,org_unit_id) values(?,?,now(),?,?,?,?)`
		sys_rdbms_009 = `update sys_user_info set user_name = ?,user_email = ?,user_phone = ?,org_unit_id = ? where user_id = ?`

		sys_rdbms_010 = `SELECT user_id,user_passwd,status_id,continue_error_cnt FROM sys_sec_user where user_id = ?`

		sys_rdbms_011 = `select distinct t2.res_url from sys_user_theme t1 inner join SYS_THEME_VALUE t2 on t1.theme_id = t2.theme_id where t1.user_id = ? and t2.res_id = ? and t2.res_type = '0'`

		sys_rdbms_012 = `select distinct t.res_id,t.res_name,v.res_url,v.res_bg_color,v.res_class, v.res_img, v.group_id from (select * from sys_resource_info t where find_in_set(res_up_id,getChildList(?))) t inner join SYS_THEME_VALUE v on t.res_id = v.res_id inner join sys_user_theme h on v.theme_id = h.theme_id inner join sys_role_user_relation r on r.user_id = h.user_id inner join sys_role_resource_relat e on r.role_id = e.role_id and e.res_id = t.res_id inner join sys_role_info f on f.role_id = e.role_id and f.role_status_id = '0' where v.res_type = ? and h.user_id = ? order by sort_id asc`

		sys_rdbms_013 = `select t2.res_url from sys_user_theme t1 inner join sys_theme_value t2 on t1.theme_id = t2.theme_id where t1.user_id = ? and t2.res_id = 'backindex'`

		sys_rdbms_014 = `update SYS_SEC_USER set user_passwd = ? where user_id = ? and user_passwd = ?`

		sys_rdbms_015 = `update SYS_SEC_USER set user_passwd = ? where user_id = ?`

		sys_rdbms_016 = `select count(*) 
							from sys_user_info t  
							inner join sys_sec_user u on t.user_id = u.user_id 
							inner join sys_user_status_attr a on u.status_id = a.status_id 
							inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
							left join sys_user_domain_rel dr on t.user_id = dr.user_id 
							left join sys_domain_info di on dr.domain_id = di.domain_id 
							where 
							t.user_id <> ? and
							exists (
							    SELECT domain_id from sys_domain_info s
							    where FIND_IN_SET(s.domain_id,getChildDomainList(?)) 
							    and s.domain_status_id = '0'
							    and di.domain_id = s.domain_id
							) and
                            exists (
							    SELECT org_unit_id from sys_org_info s
							    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
							    and s.org_status_id = '0'
							    and s.org_unit_id = i.org_unit_id
							)`
		sys_rdbms_017 = `select t.user_id,t.user_name,a.status_desc,t.user_create_date,
							t.user_owner,t.user_email,t.user_phone,t.org_unit_id,i.org_unit_desc,
							dr.domain_id,di.domain_name,t.user_maintance_date,t.user_maintance_user 
							from sys_user_info t  
							inner join sys_sec_user u on t.user_id = u.user_id 
							inner join sys_user_status_attr a on u.status_id = a.status_id 
							inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
							left join sys_user_domain_rel dr on t.user_id = dr.user_id 
							left join sys_domain_info di on dr.domain_id = di.domain_id 
							where 
							t.user_id <> ? and
							exists (
							    SELECT domain_id from sys_domain_info s
							    where FIND_IN_SET(s.domain_id,getChildDomainList(?)) 
							    and s.domain_status_id = '0'
							    and di.domain_id = s.domain_id
							) and
                            exists (
							    SELECT org_unit_id from sys_org_info s
							    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
							    and s.org_status_id = '0'
							    and s.org_unit_id = i.org_unit_id
							)
                            limit ?,?`
		sys_rdbms_018 = `insert into sys_user_info (user_id,user_name,user_create_date,user_owner,user_email,user_phone,org_unit_id,user_maintance_date,user_maintance_user) values(?,?,now(),?,?,?,?,now(),?)`
		sys_rdbms_019 = `insert into sys_sec_user(user_id,user_passwd,status_id) values(?,?,?)`
		sys_rdbms_020 = `update sys_sec_user t set t.status_id = ?,continue_error_cnt = 0 where t.user_id = ?`
		sys_rdbms_021 = `update SYS_USER_INFO t set t.user_name = ?, t.user_phone = ?, t.user_email = ? ,t.user_maintance_date = now(), t.user_maintance_user = ? where t.user_id = ?`
		sys_rdbms_022 = `SELECT T.UUID,T.User_id,t.Role_id,i.Role_name,t.maintance_date,t.maintance_user FROM sys_role_user_relation T INNER JOIN SYS_ROLE_INFO I ON T.ROLE_ID = I.ROLE_ID where t.user_id = ? limit ?,?`
		sys_rdbms_023 = `select count(*) from sys_role_user_relation T INNER JOIN SYS_ROLE_INFO I ON T.ROLE_ID = I.ROLE_ID where t.user_id = ?`
		sys_rdbms_024 = `insert into sys_role_user_relation(uuid,role_id,user_id,maintance_date,maintance_user) values(uuid(),?,?,str_to_date(?,'%Y-%m-%d'),?)`
		sys_rdbms_025 = `delete from  sys_role_user_relation where uuid = ?`
		sys_rdbms_026 = `insert into sys_role_info(role_id,role_name,role_owner,role_create_date,role_status_id,domain_id,role_maintance_date,role_maintance_user) values(?,?,?,now(),?,?,now(),?)`
		sys_rdbms_027 = `delete from  sys_role_info where role_id = ?`
		sys_rdbms_028 = `select  t.role_id,t.role_name,t.role_owner,t.role_create_date,a.role_status_desc,a.role_status_id,t.domain_id,o.domain_name,t.role_maintance_date,t.role_maintance_user
								from sys_role_info t 
								inner join sys_role_status_attr a on t.role_status_id = a.role_status_id
								inner join sys_domain_info o on t.domain_id = o.domain_id
								inner join sys_role_user_relation n on t.role_id <> n.role_id and n.user_id = ?
								where exists (
									SELECT domain_id from sys_domain_info s
									where FIND_IN_SET(s.domain_id,getChildDomainList(?))
									and t.domain_id = s.domain_id  
								) limit ?,?`
		sys_rdbms_029 = `SELECT T.UUID, T.RES_ID, I.RES_NAME, I.RES_ATTR, I.RES_UP_ID FROM sys_role_resource_relat T INNER JOIN SYS_RESOURCE_INFO I ON T.RES_ID = I.RES_ID WHERE T.ROLE_ID = ?`
		sys_rdbms_030 = `select distinct e.res_id,e.res_id,i.res_name,i.res_attr,i.res_up_id
										from sys_role_user_relation r
										inner join sys_role_resource_relat e on r.role_id = e.role_id
										inner join sys_resource_info i on e.res_id = i.res_id
										where (r.user_id = ? or 'admin' = ?)
										and not exists (
										    select 1 from sys_role_resource_relat a
										    inner join sys_resource_info f on a.res_id = f.res_id
										    where a.role_id = ? and f.res_id = i.res_id
										)`
		sys_rdbms_031 = `SELECT uuid(), I.RES_ID, I.RES_NAME, I.RES_ATTR, I.RES_UP_ID FROM SYS_RESOURCE_INFO I inner JOIN sys_role_resource_relat T ON T.RES_ID = I.RES_ID and T.ROLE_ID = ?`
		sys_rdbms_032 = `insert into sys_role_resource_relat(uuid,role_id,res_id) values(uuid(),?,?)`
		sys_rdbms_033 = `delete from  sys_role_resource_relat where uuid = ?`
		sys_rdbms_034 = `select t.domain_id as project_id, t.domain_name as project_name, t.domain_up_id, s.domain_status_name  as status_name, t.domain_create_date  as maintance_date, t.domain_owner as user_id,t.domain_maintance_date,t.domain_maintance_user
								from SYS_domain_info t inner join sys_domain_status_attr s  on t.domain_status_id = s.domain_status_id 
								where exists (
									SELECT domain_id from sys_domain_info s
									where FIND_IN_SET(s.domain_id,getChildDomainList(?))
									and ( t.domain_id = s.domain_id or t.domain_id is null ) 
								) limit ?,?`
		sys_rdbms_035 = `select t.domain_id as project_id, t.domain_name as project_name, t.domain_up_id, s.domain_status_name  as status_name, t.domain_create_date  as maintance_date, t.domain_owner as user_id,t.domain_maintance_date,t.domain_maintance_user
								from SYS_domain_info t inner join sys_domain_status_attr s  on t.domain_status_id = s.domain_status_id 
								where exists (
									SELECT domain_id from sys_domain_info s
									where FIND_IN_SET(s.domain_id,getChildDomainList(?))
									and t.domain_id = s.domain_id )`
		sys_rdbms_036 = `insert into sys_domain_info(domain_id,domain_name,domain_up_id,domain_status_id,domain_create_date,domain_owner,domain_maintance_date,domain_maintance_user) values(?,?,?,?,now(),?,now(),?)`
		sys_rdbms_037 = `delete from sys_domain_info where domain_id = ?`
		sys_rdbms_038 = `update sys_domain_info set domain_name = ?, domain_up_id=?, domain_status_id = ?, domain_maintance_date = now(), domain_maintance_user = ? where domain_id = ?`
		sys_rdbms_039 = `insert into sys_user_domain_rel(uuid,user_id,domain_id,maintance_date,grant_user_id) values(uuid(),?,?,now(),?)`
		sys_rdbms_040 = `SELECT T.RES_ID,T.RES_NAME,T.RES_ATTR, A.RES_attr_DESC,T.RES_UP_ID,T.res_type,R.RES_TYPE_DESC FROM sys_resource_info T INNER JOIN sys_resource_info_attr A ON T.RES_ATTR = A.RES_ATTR INNER JOIN SYS_RESOURCE_TYPE_ATTR R ON T.RES_TYPE = R.RES_TYPE limit ?,?`

		sys_rdbms_041 = `SELECT org_unit_id,org_unit_desc,up_org_id,t.org_status_id,r.org_status_desc,t.domain_id,i.domain_name,start_date,end_date,create_date,maintance_date,create_user,maintance_user 
								FROM sys_org_info t
								inner join sys_domain_info i on t.domain_id = i.domain_id
								inner join sys_org_status_attr r on t.org_status_id = r.org_status_id
								where exists (
										SELECT 1 from sys_domain_info s
										where FIND_IN_SET(s.domain_id,getChildDomainList(?))
										and t.domain_id = s.domain_id
								) and 
                                exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
								)
								limit ?,?`
		sys_rdbms_042 = `SELECT count(*) FROM sys_org_info t
									where exists (
											SELECT domain_id from sys_domain_info s
											where FIND_IN_SET(s.domain_id,getChildDomainList(?))
											and ( t.domain_id = s.domain_id or t.domain_id is null ) 
										) and 
										exists (
											    SELECT 1 from sys_org_info s
											    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
											    and s.org_status_id = '0'
											    and s.org_unit_id = t.org_unit_id
										)`
		sys_rdbms_043 = `insert into sys_org_info(org_unit_id,org_unit_desc,up_org_id,org_status_id,domain_id,start_date,end_date,create_date,maintance_date,create_user,maintance_user) values(?,?,?,?,?,?,?,now(),now(),?,?)`
		sys_rdbms_044 = `delete from sys_org_info where org_unit_id = ?`
		sys_rdbms_045 = `delete from sys_role_user_relation where user_id = ?`
		sys_rdbms_046 = `select t.role_id,t.role_name
							from sys_role_info t 
							where ( t.role_owner = ? or 
							exists (
								select 1 from sys_role_user_relation r
							    where r.user_id = ? and t.role_id = r.role_id
							))`
		sys_rdbms_047 = ` select t.role_id,t.role_name
							from sys_role_info t 
							where ( t.role_owner = ? or 
							exists (
								select 1 from sys_role_user_relation r
							    where r.user_id = ? and t.role_id = r.role_id
							)) and not exists (
								select 1 from sys_role_user_relation n
                                where n.user_id = ? and t.role_id = n.role_id
                            )`
		sys_rdbms_048 = `insert into sys_role_user_relation(uuid,role_id,user_id,maintance_date,maintance_user) values(?,?,?,now(),?)`

		sys_rdbms_049 = `select count(*)
									from sys_user_info t  
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and di.domain_id = ?
									and exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
									)`
		sys_rdbms_050 = `update sys_role_info t set t.role_name = ? ,t.role_status_id = ? where t.role_id = ?`

		sys_rdbms_051 = `select t.user_id,t.user_name,i.org_unit_desc
									from sys_user_info t  
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and di.domain_id = ?
									and exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
									) limit ?,?`
		sys_rdbms_052 = `select count(*)
									from sys_user_info t
                                    inner join sys_sec_user u on t.user_id = u.user_id
                                    inner join sys_user_status_attr ra on ra.status_id = u.status_id
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and di.domain_id = ?
									and exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
									)`

		sys_rdbms_053 = `select t.user_id,t.user_name,ra.status_desc,t.user_create_date,t.User_owner,t.User_email,t.User_phone,t.Org_unit_id,i.org_unit_desc,dr.Domain_id,di.domain_name,t.User_maintance_date,t.User_maintance_user
									from sys_user_info t
                                    inner join sys_sec_user u on t.user_id = u.user_id
                                    inner join sys_user_status_attr ra on ra.status_id = u.status_id
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and di.domain_id = ?
									and exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
									) limit ?,?`
		sys_rdbms_054 = `select t.user_id,t.user_name,ra.status_desc,t.user_create_date,t.User_owner,t.User_email,t.User_phone,t.Org_unit_id,i.org_unit_desc,dr.Domain_id,di.domain_name,t.User_maintance_date,t.User_maintance_user
									from sys_user_info t  
                                    inner join sys_sec_user u on t.user_id = u.user_id
                                    inner join sys_user_status_attr ra on ra.status_id = u.status_id
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ?
                                    and exists (
										SELECT 1 from sys_domain_info s
										where FIND_IN_SET(s.domain_id,getChildDomainList(?))
										and di.domain_id = s.domain_id
										and s.domain_id = ?
									) limit ?,?`
		sys_rdbms_055 = `select count(*)
									from sys_user_info t  
                                    inner join sys_sec_user u on t.user_id = u.user_id
                                    inner join sys_user_status_attr ra on ra.status_id = u.status_id
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ?
                                    and exists (
										SELECT 1 from sys_domain_info s
										where FIND_IN_SET(s.domain_id,getChildDomainList(?))
										and di.domain_id = s.domain_id
										and s.domain_id = ?
									)`
		sys_rdbms_056 = `select t.user_id,t.user_name,i.org_unit_desc
									from sys_user_info t  
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ?
                                    and exists (
										SELECT 1 from sys_domain_info s
										where FIND_IN_SET(s.domain_id,getChildDomainList(?))
										and di.domain_id = s.domain_id
										and s.domain_id = ?
									) limit ?,?`

		sys_rdbms_057 = `
									select t.user_id,t.user_name,i.org_unit_desc
									from sys_user_info t  
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and di.domain_id = ? 
									and exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
									) limit ?,?`
		sys_rdbms_058 = `select t.user_id,t.user_name,di.domain_name
									from sys_user_info t  
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and 
									exists (
										SELECT domain_id from sys_domain_info s
										where FIND_IN_SET(s.domain_id,getChildDomainList(?))
										and di.domain_id = s.domain_id
									) and t.org_unit_id = ?
									limit ?,?`
		sys_rdbms_059 = `select count(*)
									from sys_user_info t  
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and 
									exists (
										SELECT domain_id from sys_domain_info s
										where FIND_IN_SET(s.domain_id,getChildDomainList(?))
										and di.domain_id = s.domain_id
									) and t.org_unit_id = ?`

		sys_rdbms_060 = `SELECT org_unit_id,org_unit_desc,up_org_id,t.org_status_id,r.org_status_desc,t.domain_id,i.domain_name,start_date,end_date,create_date,maintance_date,create_user,maintance_user 
								FROM sys_org_info t
								inner join sys_domain_info i on t.domain_id = i.domain_id
								inner join sys_org_status_attr r on t.org_status_id = r.org_status_id
								where exists (
								    SELECT 1 from sys_org_info s
								    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
								    and s.org_status_id = '0'
								    and s.org_unit_id = t.org_unit_id
								) and  exists (
									SELECT 1 from sys_domain_info s
									where FIND_IN_SET(s.domain_id,getChildDomainList(?))
									and t.domain_id = s.domain_id
								    and s.domain_id = ?
								) `
		sys_rdbms_061 = `SELECT org_unit_id,org_unit_desc,up_org_id,t.org_status_id,r.org_status_desc,t.domain_id,i.domain_name,start_date,end_date,create_date,maintance_date,create_user,maintance_user 
								FROM sys_org_info t
								inner join sys_domain_info i on t.domain_id = i.domain_id
								inner join sys_org_status_attr r on t.org_status_id = r.org_status_id
								where exists (
									SELECT 1 from sys_domain_info s
									where FIND_IN_SET(s.domain_id,getChildDomainList(?))
									and t.domain_id = s.domain_id
								    and s.domain_id = ?
								)`
		sys_rdbms_062 = `select count(*)
									from sys_user_info t
                                    inner join sys_sec_user u on t.user_id = u.user_id
                                    inner join sys_user_status_attr ra on ra.status_id = u.status_id
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and di.domain_id = ?
									and exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
									) and t.org_unit_id = ?`
		sys_rdbms_063 = `select count(*)
									from sys_user_info t  
                                    inner join sys_sec_user u on t.user_id = u.user_id
                                    inner join sys_user_status_attr ra on ra.status_id = u.status_id
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ?
                                    and exists (
										SELECT 1 from sys_domain_info s
										where FIND_IN_SET(s.domain_id,getChildDomainList(?))
										and di.domain_id = s.domain_id
										and s.domain_id = ?
									) and t.org_unit_id = ?`
		sys_rdbms_064 = `select t.user_id,t.user_name,ra.status_desc,t.user_create_date,t.User_owner,t.User_email,t.User_phone,t.Org_unit_id,i.org_unit_desc,dr.Domain_id,di.domain_name,t.User_maintance_date,t.User_maintance_user
									from sys_user_info t
                                    inner join sys_sec_user u on t.user_id = u.user_id
                                    inner join sys_user_status_attr ra on ra.status_id = u.status_id
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and di.domain_id = ?
									and exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
									) and t.org_unit_id = ? limit ?,?`
		sys_rdbms_065 = `select t.user_id,t.user_name,ra.status_desc,t.user_create_date,t.User_owner,t.User_email,t.User_phone,t.Org_unit_id,i.org_unit_desc,dr.Domain_id,di.domain_name,t.User_maintance_date,t.User_maintance_user
									from sys_user_info t  
                                    inner join sys_sec_user u on t.user_id = u.user_id
                                    inner join sys_user_status_attr ra on ra.status_id = u.status_id
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ?
                                    and exists (
										SELECT 1 from sys_domain_info s
										where FIND_IN_SET(s.domain_id,getChildDomainList(?))
										and di.domain_id = s.domain_id
										and s.domain_id = ?
									) and t.org_unit_id = ? limit ?,?`
		sys_rdbms_066 = `select count(*)
									from sys_user_info t  
									inner join sys_org_info i on i.org_unit_id = t.org_unit_id 
									left join sys_user_domain_rel dr on t.user_id = dr.user_id 
									left join sys_domain_info di on dr.domain_id = di.domain_id 
									where t.user_id <> ? and di.domain_id = ?
									and exists (
									    SELECT 1 from sys_org_info s
									    where FIND_IN_SET(s.org_unit_id,getChildOrgList(?)) 
									    and s.org_status_id = '0'
									    and s.org_unit_id = t.org_unit_id
									)`
		sys_rdbms_067 = `SELECT domain_id,domain_name,domain_up_id,domain_status_id from sys_domain_info s where FIND_IN_SET(s.domain_id,getChildDomainList(?))`

	}
}
