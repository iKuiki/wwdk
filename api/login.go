package api

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/ikuiki/wwdk/conf"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// JsLogin 获取uuid
// 初次打开时调用，获取uuid
// @return uuid 获取到的uuid
func (api *WechatwebAPI) JsLogin() (uuid string, body []byte, err error) {
	params := url.Values{}
	params.Set("appid", conf.AppID)
	params.Set("fun", "new")
	params.Set("lang", conf.Lang)
	params.Set("_", tool.GetWxTimeStamp())
	req, _ := http.NewRequest("GET", "https://login."+api.apiDomain+"/jslogin?"+params.Encode(), nil)
	resp, err := api.request(req)
	if err != nil {
		return "", body, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	ret := tool.ExtractWxWindowRespond(string(body))
	if ret["window.QRLogin.code"] != "200" {
		return "", body, errors.New("window.QRLogin.code = " + ret["window.QRLogin.code"])
	}
	uuid = ret["window.QRLogin.uuid"]
	return uuid, body, nil
}

// Login 等待用户扫码登陆
// 获取uuid后就可以生成二维码等用户扫码了
// @param tip tip参数，一般第一轮为1，第二轮开始就为0
// @param uuid getUuid接口获取到的uuid
// @return code 返回的windows.code，状态码
// @return userAvatar 当用户扫码后返回用户头像
// @return redirectURL 用户确认登陆后返回重定向地址
func (api *WechatwebAPI) Login(tip, uuid string) (code, userAvatar, redirectURL string, body []byte, err error) {
	params := url.Values{}
	params.Set("tip", tip)
	params.Set("uuid", uuid)
	params.Set("_", tool.GetWxTimeStamp())
	req, _ := http.NewRequest(`GET`, "https://login."+api.apiDomain+"/cgi-bin/mmwebwx-bin/login?"+params.Encode(), nil)
	resp, err := api.request(req)
	if err != nil {
		err = errors.Errorf("waitForScan request error: %v", err)
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	ret := tool.ExtractWxWindowRespond(string(body))
	code, userAvatar, redirectURL = ret["window.code"], ret["window.userAvatar"], ret["window.redirect_uri"]
	return
}

// WebwxNewLoginPage 获取登陆凭据
// 用户扫码确认登陆后，获取登陆凭据
// @param redirectURL 当用户扫码确认登陆后获取到的redirectURL
func (api *WechatwebAPI) WebwxNewLoginPage(redirectURL string) (body []byte, err error) {
	req, _ := http.NewRequest(`GET`, redirectURL+"&fun=new", nil) // 统一不加version=v2了
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("getCookie request error: " + err.Error())
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Read respond body error: " + err.Error())
		return
	}
	var bodyResp datastruct.GetCookieRespond
	err = xml.Unmarshal(body, &bodyResp)
	if err != nil {
		err = errors.New("Unmarshal respond xml error: " + err.Error())
		return
	}
	// 此处获取到skey与passTicket
	api.loginInfo.SKey = bodyResp.Skey
	api.loginInfo.PassTicket = bodyResp.PassTicket
	return
}

// WebwxInit 初始化微信
// 当获取到登陆凭证后，即可调用此接口初始化微信获取基本信息
// @return user 当前登陆的用户的信息
// @return contactList 部分联系人列表⚠️此列表不全，要和后面获取联系人的列表合并，切记切记
func (api *WechatwebAPI) WebwxInit() (user *datastruct.User, contactList []datastruct.Contact, body []byte, err error) {
	reqBody, err := json.Marshal(datastruct.WxInitRequestBody{
		BaseRequest: api.baseRequest(),
	})
	if err != nil {
		err = errors.New("json.Marshal error: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	// params.Set("skey", api.loginInfo.sKey)
	params.Set("r", tool.GetWxTimeStamp())
	// resp, err := api.apiRuntime.client.Post("https://"+api.apiRuntime.apiDomain+"/cgi-bin/mmwebwx-bin/webwxinit?"+params.Encode(),
	// 	"application/json;charset=UTF-8",
	// 	bytes.NewReader(reqBody))

	req, err := http.NewRequest("POST",
		"https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxinit?"+params.Encode(),
		bytes.NewReader(reqBody))
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("do request error: " + err.Error())
		return
	}
	defer resp.Body.Close()

	respStruct := datastruct.WxInitRespond{}
	body, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		err = errors.New("respond json Unmarshal to struct fail: " + err.Error())
		return
	}
	if respStruct.BaseResponse.Ret != 0 {
		err = errors.Errorf("respond ret(%d) error: %s", respStruct.BaseResponse.Ret, respStruct.BaseResponse.ErrMsg)
		return
	}
	user = respStruct.User
	contactList = respStruct.ContactList
	api.loginInfo.SyncKey = respStruct.SyncKey
	api.loginInfo.SKey = respStruct.SKey
	return
}
