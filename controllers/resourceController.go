package controllers

import (
	"io/ioutil"

	"github.com/astaxie/beego/context"
)

type ResourceController struct {
}

var ResourceCtl = &ResourceController{}

func (ResourceController) Page(ctx *context.Context) {
	hz, _ := ioutil.ReadFile("./views/hauth/res_info_page.tpl")
	ctx.ResponseWriter.Write(hz)
}
