package rpc

import (
	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/token/hjwt"
	"github.com/hzwy23/hauth/logs"
	"github.com/hzwy23/dbobj"
)

func HaveRightsById(ctx *context.Context,id string) bool {
	cookie, _ := ctx.Request.Cookie("Authorization")
	jc, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		return false
	}
	status := 0
	err=dbobj.QueryRow(sys_rdbms_078,jc.User_id,id).Scan(&status)
	if err!=nil{
		logs.Error("no rights")
		return false
	}
	if status==1{
		return true
	}
	return false
}

func HaveRightsByUri(ctx *context.Context)bool{

	cookie, _ := ctx.Request.Cookie("Authorization")
	jc, err := hjwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		return false
	}

	url := ctx.Request.RequestURI
	status := 0
	err=dbobj.QueryRow(sys_rdbms_079,jc.User_id,url).Scan(&status)

	if err!=nil{
		logs.Error("no rights")
		return false
	}
	if status==1{
		return true
	}
	return false

}