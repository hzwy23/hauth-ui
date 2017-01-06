package resource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hzwy23/hcloud/logs"

	"github.com/hzwy23/hcloud/plat/auth"
	"github.com/hzwy23/hcloud/plat/mgr/sqlText"
	"github.com/hzwy23/hcloud/plat/route"
	"github.com/hzwy23/hcloud/plat/session"

	"github.com/hzwy23/dbobj"
)

type DefaultMenu struct {
	Res_id       string
	Res_name     string
	Res_bg_color string
	Res_class    string
	Res_url      string
	Res_img      string
	Group_id     string
	route.RouteControl
}

func (this *DefaultMenu) Get(w http.ResponseWriter, r *http.Request) {
	if auth.Access(w, r) == false {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("权限不足"))
		return
	}
	logs.Debug("请求菜单列表信息")
	typeId := r.FormValue("TypeId")
	Id := r.FormValue("Id")
	user := session.Get(w, r, "userId")
	sql := sqlText.PLATFORM_RESOURCE_DEFAULTMENULIST1
	var one DefaultMenu
	var rst []DefaultMenu

	rows, err := dbobj.Query(sql, Id, typeId, user)
	defer rows.Close()
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("请联系管理员,查询主题菜单失败"))
		return
	}
	for rows.Next() {
		err := rows.Scan(&one.Res_id,
			&one.Res_name,
			&one.Res_url,
			&one.Res_bg_color,
			&one.Res_class,
			&one.Res_img,
			&one.Group_id)
		if err != nil {
			fmt.Println(err)
			logs.Error(err)
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("请联系管理员,查询主题菜单失败"))
			return
		}
		rst = append(rst, one)
	}

	ojs, err := json.Marshal(rst)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("打包json失败,请联系管理员"))
		return
	}
	w.Write(ojs)
}
