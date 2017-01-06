package route

import (
	"net/http"

	"github.com/astaxie/beego"
)

type Controller struct {
	beego.Controller
}

func Router(rootpath string, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
	beego.BeeApp.Handlers.Add(rootpath, c, mappingMethods...)
	return beego.BeeApp
}

func Get(rootpath string, f beego.FilterFunc) *beego.App {
	beego.BeeApp.Handlers.Get(rootpath, f)
	return beego.BeeApp
}

func Post(rootpath string, f beego.FilterFunc) *beego.App {
	beego.BeeApp.Handlers.Post(rootpath, f)
	return beego.BeeApp
}

func Delete(rootpath string, f beego.FilterFunc) *beego.App {
	beego.BeeApp.Handlers.Delete(rootpath, f)
	return beego.BeeApp
}

func Put(rootpath string, f beego.FilterFunc) *beego.App {
	beego.BeeApp.Handlers.Put(rootpath, f)
	return beego.BeeApp
}

func Any(rootpath string, f beego.FilterFunc) *beego.App {
	beego.BeeApp.Handlers.Any(rootpath, f)
	return beego.BeeApp
}

func Handler(rootpath string, h http.Handler, options ...interface{}) *beego.App {
	beego.BeeApp.Handlers.Handler(rootpath, h, options...)
	return beego.BeeApp
}

func InsertFilter(pattern string, pos int, filter beego.FilterFunc, params ...bool) *beego.App {
	beego.BeeApp.Handlers.InsertFilter(pattern, pos, filter, params...)
	return beego.BeeApp
}

func Run(params ...string) {
	beego.Run(params...)
}
