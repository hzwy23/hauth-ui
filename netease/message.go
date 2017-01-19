package netease

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"time"

	"github.com/hzwy23/hauth/logs"
)

type msgCodeInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Obj  string `json:"obj"`
}

type verifyCode struct {
	Code int `json:"code"`
}

var (
	msgurl    = "https://api.netease.im/sms/sendcode.action"
	verifyurl = "https://api.netease.im/sms/verifycode.action"
	appKey    = "0b5eb2a6a180bbe402aa732acbe80e71"
	appSecret = "0de6a3d20d2e"
	nonce     = "yph2b"
)

func genSHA1(sec, nonce, curtme string) string {
	var sum = sec + nonce + curtme
	h := sha1.New()
	h.Write([]byte(sum))
	bs := h.Sum(nil)
	sha := fmt.Sprintf("%x", bs)
	return sha
}

func checkCode(message []byte) (*msgCodeInfo, error) {
	var msg msgCodeInfo
	err := json.Unmarshal(message, &msg)
	if err != nil {
		logs.Error(err)
		return &msg, err
	}
	if msg.Code != 200 {
		return &msg, err
	}
	return &msg, nil
}

func SendCode(phone string) (*msgCodeInfo, error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := http.Client{Transport: tr}
	client.Jar, _ = cookiejar.New(nil)

	var mobile = url.Values{"mobile": {phone}}

	req, err := http.NewRequest("POST", msgurl, strings.NewReader(mobile.Encode()))
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	var curTime = strconv.Itoa(int(time.Now().Unix()))
	var checkSum = genSHA1(appSecret, nonce, curTime)

	req.Header.Add("AppKey", appKey)
	req.Header.Add("Nonce", nonce)
	req.Header.Add("CurTime", curTime)
	req.Header.Add("CheckSum", checkSum)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) //发送
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	msg, err := checkCode(data)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return msg, nil
}

func VerifyCode(phone string, code string) (int, bool) {

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	client := http.Client{Transport: tr}
	client.Jar, _ = cookiejar.New(nil)

	var mobile = url.Values{"mobile": {phone}, "code": {code}}

	req, err := http.NewRequest("POST", verifyurl, strings.NewReader(mobile.Encode()))
	if err != nil {
		logs.Error(err)
		return 408, false
	}

	var curTime = strconv.Itoa(int(time.Now().Unix()))
	var checkSum = genSHA1(appSecret, nonce, curTime)

	req.Header.Add("AppKey", appKey)
	req.Header.Add("Nonce", nonce)
	req.Header.Add("CurTime", curTime)
	req.Header.Add("CheckSum", checkSum)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) //发送
	if err != nil {
		logs.Error(err)
		return 405, false
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return 406, false
	}
	var ret verifyCode
	err = json.Unmarshal(data, &ret)
	if err != nil {
		logs.Error(err)
		return 407, false
	}
	if ret.Code == 200 {
		return ret.Code, true
	} else {
		return ret.Code, false
	}
}
