// Copyright 2009 hzwy23. All rights reserved.
// Package sys provides base common function
package auth

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hzwy23/hcloud/logs"
	"github.com/hzwy23/hcloud/plat/session"

	_ "github.com/astaxie/beego/session/redis"
)

// genToken function is general md5 info
// this function return two values
// first return is token value
// second return is error info,if it return nil, is right
func genToken() (string, error) {
	curitme := time.Now().Unix()
	md := md5.New()
	io.WriteString(md, strconv.FormatInt(curitme, 10))
	val := fmt.Sprintf("%x", md.Sum([]byte{'h', 'z', 'w'}))
	return val, nil
}

// SetToken function , insert token into session
// session use Privilege's handle
// It is right when return nil
func SetToken(w http.ResponseWriter, r *http.Request) error {
	token, err := genToken()
	if err != nil {
		logs.Error("cant not general token")
		return err
	}
	return session.Set(w, r, map[string]string{"token": token})
}

// GetToken function, get token info from session
// return token value
func GetToken(w http.ResponseWriter, r *http.Request) string {
	return session.Get(w, r, "token")
}
