package rdbms

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/auth/sysuserinfo/:offset/:limit/:order/:mode", &SysUserInfo{}, "get:GetAllUsers")
	beego.Router("/v1/auth/sysuserinfo/:offset/:limit", &SysUserInfo{}, "get:GetAllUsersNoOrder")
	beego.Router("/v1/auth/sysuserinfo/?:user_id", &SysUserInfo{})
}
