package mgr

import (
	"github.com/hzwy23/hcloud/plat/mgr/resource"

	"github.com/hzwy23/hcloud/plat/mgr/start"

	"github.com/hzwy23/hcloud/plat/route"
)

func init() {
	route.AddRoute("/", &start.IndexSystem{})
	route.AddRoute("/login", &start.LoginSystem{})
	route.AddRoute("/logout", &start.LogoutSystem{})
	route.AddRoute("/plat/help", &resource.PlatMgrHelp{})
	route.AddRoute("/platform/user/domain", &resource.UserDomainRel{})
	//系统管理－

	route.AddRoute("/platform/menu", &resource.Menu{})
	route.AddRoute("/platform/menuMgr/Page", &resource.MenuPage{})
	route.AddRoute("/platform/menuMgr", &resource.MenuMgr{})

	route.AddRoute("/platform/passwd", &resource.Passwd{})
	//系统管理－入口页面
	route.AddRoute("/platform/select", &resource.GoEntry{})
	//导航页
	route.AddRoute("/platform/IndexPage", &resource.IndexPage{})
	//域名管理
	route.AddRoute("/platform/DomainMgr", &resource.ProjectMgr{})
	route.AddRoute("/platform/DomainMgr/page", &resource.ProjectPage{})

	route.AddRoute("/platform/HomePage", &resource.HomePage{})
	//用户管理
	route.AddRoute("/platform/UserInfo", &resource.UserInfo{})
	route.AddRoute("/platform/UserInfoPage", &resource.UserInfoPage{})
	//菜单管理
	route.AddRoute("/platform/ResInfoPage", &resource.ResInfoPage{})
	route.AddRoute("/platform/ResInfo", &resource.ResInfo{})
	//角色管理
	route.AddRoute("/platform/RoleInfoPage", &resource.RoleInfoPage{})
	route.AddRoute("/platform/RoleInfo", &resource.RoleInfo{})
	//日志管理
	route.AddRoute("/platform/HandleLogs", &resource.HandleLog{})
	route.AddRoute("/platform/HandleLogsPage", &resource.HandleLogPage{})
	//角色资源管理
	route.AddRoute("/platform/ResRoleRelPage", &resource.ResourceRoleRelPage{})
	route.AddRoute("/platform/ResRoleRel", &resource.ResourceRoleRel{})
	//角色域
	route.AddRoute("/platform/RoleDomainRelPage", &resource.RoleDomainRelPage{})
	route.AddRoute("/platform/RoleDomainRel", &resource.RoleDomainRel{})
	//用户角色
	route.AddRoute("/platform/RoleUserRelPage", &resource.RoleUserRelPage{})
	route.AddRoute("/platform/RoleUserRel", &resource.RoleUserRel{})

	route.AddRoute("/platform/DefaultMenu", &resource.DefaultMenu{})
}
