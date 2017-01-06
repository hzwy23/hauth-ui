package route

import (
	"net/http"

	"reflect"

	"github.com/hzwy23/hcloud/logs"

	"encoding/json"

	"github.com/astaxie/beego"
)

var routeMap = make(map[string]interface{})

type BeegoControl struct {
	beego.Controller
}

type RouteInterface interface {
	Post(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
}

type RouteControl struct {
}

// Insert menu into menu table
func (this *RouteControl) Post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Please implement Post method"))
}

// Delete menu info from menu table
func (this *RouteControl) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Please implement Delete method"))
}

// Update menu info from menu table
func (this *RouteControl) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Please implement Get method"))
}

func (this *RouteControl) Put(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Please implement Put method"))
}

func (this *BeegoControl) Post() {
	if this.Ctx.Request.FormValue("_Method") == "Delete" {
		this.Delete()
		return
	}
	pat := this.Data["RouterPattern"]
	f := routeMap[pat.(string)]
	rf := reflect.ValueOf(f)
	function := rf.MethodByName("Post")
	param := make([]reflect.Value, 2)
	param[0] = reflect.ValueOf(this.Ctx.ResponseWriter)
	param[1] = reflect.ValueOf(this.Ctx.Request)
	function.Call(param)
}

func (this *BeegoControl) Get() {
	pat := this.Data["RouterPattern"]
	f := routeMap[pat.(string)]
	rf := reflect.ValueOf(f)
	function := rf.MethodByName("Get")
	param := make([]reflect.Value, 2)
	param[0] = reflect.ValueOf(this.Ctx.ResponseWriter)
	param[1] = reflect.ValueOf(this.Ctx.Request)
	function.Call(param)
}

func (this *BeegoControl) Delete() {
	pat := this.Data["RouterPattern"]
	f := routeMap[pat.(string)]
	rf := reflect.ValueOf(f)
	function := rf.MethodByName("Delete")
	param := make([]reflect.Value, 2)
	param[0] = reflect.ValueOf(this.Ctx.ResponseWriter)
	param[1] = reflect.ValueOf(this.Ctx.Request)
	function.Call(param)
}

func (this *BeegoControl) Put() {

	pat := this.Data["RouterPattern"]
	f := routeMap[pat.(string)]
	rf := reflect.ValueOf(f)
	function := rf.MethodByName("Put")
	param := make([]reflect.Value, 2)
	param[0] = reflect.ValueOf(this.Ctx.ResponseWriter)
	param[1] = reflect.ValueOf(this.Ctx.Request)
	function.Call(param)
}

func AddRoute(rootpath string, rc RouteInterface) *beego.App {
	bc := &BeegoControl{}
	routeMap[rootpath] = rc
	beego.BeeApp.Handlers.Add(rootpath, bc)
	return beego.BeeApp
}

func Router(rootpath string, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
	return beego.Router(rootpath, c, mappingMethods...)
}

func (this *RouteControl) WritePage(w http.ResponseWriter, total int, rows interface{}) {
	type page struct {
		Total int         `json:"total"`
		Rows  interface{} `json:"rows"`
	}
	var p page
	p.Total = total
	p.Rows = rows
	ijs, err := json.Marshal(p)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("打包json数据失败"))
		return
	}
	w.Write(ijs)
}

func (this *RouteControl) WriteJson(w http.ResponseWriter, v interface{}) error {
	ojs, err := json.Marshal(v)
	if err != nil {
		logs.Error(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(err.Error()))
		return err
	}
	w.Write(ojs)
	return nil
}
