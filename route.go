package main

import (
	"github.com/astaxie/beego"
	"github.com/hzwy23/hauth/controllers"
)

func init() {

	beego.Get("/HomePage", controllers.HomePage)

	beego.Post("/login", controllers.LoginSystem)

	beego.Any("/logout", controllers.LogoutSystem)

	beego.Get("/", controllers.IndexPage)

	beego.Get("/v1/auth/index/entry", controllers.SubSystemEntry)

	beego.Get("/v1/auth/main/menu", controllers.HomePageMenus)

	beego.Post("/v1/auth/passwd/update", controllers.PasswdController.PostModifyPasswd)

	beego.Post("/v1/auth/passwd/modify", controllers.PasswdController.AdminModifyPasswd)

	beego.Post("/v1/auth/role/users/clean", controllers.UserRolesController.CleanUserRoles)

	beego.Get("/v1/auth/HandleLogsPage", controllers.HandleLogsCtl.GetHandleLogPage)

	beego.Get("/v1/auth/handle/logs", controllers.HandleLogsCtl.GetHandleLogs)

	beego.Get("/v1/auth/handle/logs/search", controllers.HandleLogsCtl.SerachLogs)

	beego.Get("/v1/auth/domain/page", controllers.DomainCtl.GetDomainInfoPage)
	beego.Get("/v1/auth/domain/share/page",controllers.DomainShareCtl.Page)
	beego.Get("/v1/auth/domain/share/get",controllers.DomainShareCtl.Get)

	beego.Get("/v1/auth/domain/get", controllers.DomainCtl.GetDomainInfo)
	beego.Post("/v1/auth/domain/post", controllers.DomainCtl.PostDomainInfo)
	beego.Post("/v1/auth/domain/delete", controllers.DomainCtl.DeleteDomainInfo)
	beego.Put("/v1/auth/domain/update", controllers.DomainCtl.UpdateDomainInfo)
	beego.Get("/v1/auth/domain/owner",controllers.DomainCtl.GetDomainOwner)
	beego.Get("/v1/auth/domain/row/details",controllers.DomainCtl.GetDetails)

	beego.Get("/v1/auth/batch/page", controllers.AuthroityCtl.GetBatchPage)
	beego.Post("/v1/auth/batch/grant", controllers.AuthroityCtl.BatchGrants)

	beego.Get("/v1/auth/roles/getted", controllers.AuthroityCtl.GetGettedRoles)
	beego.Get("/v1/auth/roles/canGrant", controllers.AuthroityCtl.CanGrantRoles)

	beego.Get("/v1/auth/relation/domain/org", controllers.OrgCtl.GetSubOrgInfo)
	beego.Get("/v1/auth/resource/org/page", controllers.OrgCtl.GetOrgPage)
	beego.Get("/v1/auth/resource/org/get", controllers.OrgCtl.GetSysOrgInfo)
	beego.Post("/v1/auth/resource/org/insert", controllers.OrgCtl.InsertOrgInfo)
	beego.Post("/v1/auth/resource/org/delete", controllers.OrgCtl.DeleteOrgInfo)
	beego.Put("/v1/auth/resource/org/update", controllers.OrgCtl.UpdateOrgInfo)
	beego.Get("/v1/auth/resource/page", controllers.ResourceCtl.Page)
	beego.Get("/v1/auth/user/page", controllers.UserCtl.Page)
	beego.Get("/v1/auth/user/get/default",controllers.UserCtl.Get)
	beego.Post("/v1/auth/user/post",controllers.UserCtl.Post)
	beego.Post("/v1/auth/user/delete",controllers.UserCtl.Delete)

	beego.Get("/v1/auth/role/page", controllers.RoleCtl.Page)

	beego.Get("/v1/auth/role/get", controllers.RoleCtl.GetRoleInfo)
	beego.Post("/v1/auth/role/post", controllers.RoleCtl.PostRoleInfo)
	beego.Post("/v1/auth/role/delete", controllers.RoleCtl.DeleteRoleInfo)
	beego.Put("/v1/auth/role/update", controllers.RoleCtl.UpdateRoleInfo)
	beego.Get("/v1/auth/role/resource/details",controllers.RoleCtl.ResourcePage)

	beego.Get("/v1/auth/handle/logs/download", controllers.HandleLogsCtl.Download)
	beego.Get("/v1/auth/resource/org/download",controllers.OrgCtl.Download)
}
