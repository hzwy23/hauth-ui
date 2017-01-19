package rdbms

import (
	"html/template"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type mSysRoleUsersRel struct {
}

func getRoleUsersRel(ctx *context.Context) {
	hz, _ := template.ParseFiles("./views/platform/resource/role_users_rel.tpl")
	hz.Execute(ctx.ResponseWriter, nil)
}

func init() {
	beego.Get("/v1/auth/userAndRoles/page", getRoleUsersRel)
}
