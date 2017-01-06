package auth

import (
	"fmt"
	"net/http"
	"text/template"
)

func Sett(w http.ResponseWriter, r *http.Request) {
	val, err := Token.SetToken(w, r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("hello set,token is :", val)
	hz, _ := template.ParseFiles("./static/demo/login.tpl")
	hz.Execute(w, val)
}

func Gett(w http.ResponseWriter, r *http.Request) {
	val, err := Token.GetToken(w, r)
	if err != nil {
		fmt.Println("get token failed...")
	} else {
		fmt.Println("hello get token is :", val)
	}

	newVal, err := Token.SetToken(w, r)
	if err != nil {
		fmt.Println("new token failed...")
	}
	hz, _ := template.ParseFiles("./static/demo/login.tpl")
	hz.Execute(w, newVal)
}
