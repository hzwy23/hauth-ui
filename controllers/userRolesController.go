package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego/context"
	"github.com/hzwy23/hauth/models"
	"github.com/hzwy23/hauth/utils/hret"
	"github.com/hzwy23/hauth/utils/logs"
)

type userRolesController struct {
	models *models.UserRolesModel
}

var UserRolesController = &userRolesController{
	models: new(models.UserRolesModel),
}

func (this userRolesController) CleanUserRoles(ctx *context.Context) {
	ctx.Request.ParseForm()
	var rst []models.UserRolesModel
	err := json.Unmarshal([]byte(ctx.Request.FormValue("JSON")), &rst)
	if err != nil {
		logs.Error(err)
		hret.WriteHttpErrMsgs(ctx.ResponseWriter, http.StatusExpectationFailed, "unmarshal failed.", err)
		return
	}
	err = this.models.CleanRoles(rst)
	hret.WriteHttpOkMsgs(ctx.ResponseWriter, "clean roles of user successfully.")
}
