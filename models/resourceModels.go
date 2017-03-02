package models

import (
	"github.com/hzwy23/dbobj"
	"github.com/astaxie/beego/logs"
	"fmt"
)

type ResourceModel struct{
}

type resData struct{
	Res_id string `json:"res_id"`
	Res_name string `json:"res_name"`
	Res_attr string `json:"res_attr"`
	Res_attr_desc string `json:"res_attr_desc"`
	Res_up_id string `json:"res_up_id"`
	Res_type string `json:"res_type"`
	Res_type_desc string `json:"res_type_desc"`
}

type themeData struct{
	Theme_id      string `json:"theme_id"`
	Theme_desc    string `json:"theme_desc"`
	Res_id        string `json:"res_id"`
	Res_url       string `json:"res_url"`
	Res_type      string `json:"res_type"`
	Res_bg_color  string `json:"res_bg_color"`
	Res_class     string `json:"res_class"`
	Group_id      string `json:"group_id"`
	Res_img       string `json:"res_img"`
	Sort_id       string `json:"sort_id"`
}

func (ResourceModel)Get()([]resData,error){
	rows,err:=dbobj.Query(sys_rdbms_071)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	var rst []resData
	err=dbobj.Scan(rows,&rst)
	return rst,err
}

func (ResourceModel)Query(res_id string)([]resData,error){

	rows,err:=dbobj.Query(sys_rdbms_089,res_id)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	var rst []resData
	err = dbobj.Scan(rows,&rst)
	return rst,err
}

func (ResourceModel)QueryTheme(res_id string,theme_id string)([]themeData,error){

	rows,err:=dbobj.Query(sys_rdbms_070,theme_id,res_id)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	var rst []themeData
	err = dbobj.Scan(rows,&rst)
	return rst,err
}

func (ResourceModel)Post(res_id,res_name,res_attr,res_up_id,res_type,theme_id,res_url,res_bg_color,res_class,group_id,res_img,sort_id string)error{
	tx,err:=dbobj.Begin()
	if err!=nil{
		logs.Error(err)
		return err
	}
	_,err=tx.Exec(sys_rdbms_072,res_id,res_name,res_attr,res_up_id,res_type)
	if err!=nil{
		logs.Error(err)
		tx.Rollback()
		return err
	}
	_,err=tx.Exec(sys_rdbms_073,theme_id,res_id,res_url,res_type,res_bg_color,res_class,group_id,res_img,sort_id)
	if err!=nil{
		logs.Error(err)
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(sys_rdbms_074, "vertex_root_join_sysadmin", res_id)
	if err != nil {
		logs.Error(err)
		tx.Rollback()
	}
	return 	tx.Commit()
}

func (this ResourceModel)search(rst []resData,all []resData){
	for _,val:=range rst{
		for i,v:=range all{
			if val.Res_id == v.Res_up_id{
				rst = append(rst,v)
				if len(all)==i+1{
					all = all[:i]
				}else{
					all = append(all[:i],all[i+1:]...)
				}

				this.search(rst,all)
			}
		}
	}
}

func (this ResourceModel)Delete(res_id string)error{
	var rst []resData
	rst = append(rst,resData{Res_id:res_id,})
	all,err:=this.Get()
	if err!=nil{
		logs.Error(err)
		return err
	}

	this.search(rst,all)
	fmt.Println(all)
	fmt.Println(rst)
	return nil

	tx, err := dbobj.Begin()
	if err != nil {
		logs.Error(err)
		return err
	}
	for _, val := range rst {

		_, err = tx.Exec(sys_rdbms_075, val.Res_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(sys_rdbms_076, val.Res_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(sys_rdbms_077, val.Res_id)
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}