package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

var (
	// ErrEmptyLoginInfo 登陆信息为空
	ErrEmptyLoginInfo = errors.New("empty login info")
)

type wechatwebAPIMarshalData struct {
	UserAgent string
	APIDomain string
	DeviceID  string
	LoginInfo LoginInfo
	CookieMap map[string][]*http.Cookie
}

func (api *wechatwebAPI) Marshal() (data []byte, err error) {
	// 储存cookie
	cookieMap := make(map[string][]*http.Cookie)
	for _, host := range []string{
		api.apiDomain,
		"webpush." + api.apiDomain,
		"file." + api.apiDomain,
		// "login." + api.apiDomain,
		".qq.com",
	} {
		u, _ := url.Parse("https://" + host)
		cookieMap[host] = api.client.Jar.Cookies(u)
	}
	data, err = json.Marshal(wechatwebAPIMarshalData{
		UserAgent: api.userAgent,
		APIDomain: api.apiDomain,
		DeviceID:  api.deviceID,
		LoginInfo: api.loginInfo,
		CookieMap: cookieMap,
	})
	return
}

func (api *wechatwebAPI) Unmarshal(data []byte) (err error) {
	var dataStruct wechatwebAPIMarshalData
	err = json.Unmarshal(data, &dataStruct)
	if err == nil {
		if dataStruct.UserAgent == "" ||
			dataStruct.APIDomain == "" ||
			dataStruct.DeviceID == "" {
			// 判断为不合法的恢复
			err = ErrEmptyLoginInfo
			return
		}
		api.userAgent = dataStruct.UserAgent
		api.apiDomain = dataStruct.APIDomain
		api.deviceID = dataStruct.DeviceID
		api.loginInfo = dataStruct.LoginInfo
		for _, host := range []string{
			api.apiDomain,
			"webpush." + api.apiDomain,
			"file." + api.apiDomain,
			// "login." + api.apiDomain,
			".qq.com",
		} {
			u, _ := url.Parse("https://" + host)
			api.client.Jar.SetCookies(u, dataStruct.CookieMap[host])
		}
	}
	return
}
