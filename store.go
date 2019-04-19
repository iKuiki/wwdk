package wwdk

// 此文件中存储的为WechatWeb读写登录凭据的方法

import (
	"encoding/json"
	"github.com/ikuiki/wwdk/datastruct"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// storeLoginInfo 用于储存的登录信息
type storeLoginInfo struct {
	Cookies    map[string][]*http.Cookie
	Cookie     wechatCookie
	SyncKey    *datastruct.SyncKey
	SKey       string
	PassTicket string
	RunInfo    WechatRunInfo // 运行统计信息
	DeviceID   string        // 由客户端生成，为e+15位随机数
	User       *datastruct.User
}

// 重置登录信息
func (wxwb *WechatWeb) resetLoginInfo() error {
	func() {
		defer func() {
			if r := recover(); r != nil {
				wxwb.logger.Infof("Recovered in resetLoginInfo: %v\n", r)
				wxwb.runInfo.PanicCount++
			}
		}()
	}()
	if wxwb.loginStorer != nil {
		wxwb.loginStorer.Truncate()
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	wxwb.apiRuntime.client.Jar = jar
	wxwb.loginInfo = wechatLoginInfo{}
	return nil
}

// 往storer中写入信息
func (wxwb *WechatWeb) writeLoginInfo() error {
	func() {
		defer func() {
			if r := recover(); r != nil {
				wxwb.logger.Infof("Recovered in writeLoginInfo: %v\n", r)
				wxwb.runInfo.PanicCount++
			}
		}()
	}()
	cookieMap := make(map[string][]*http.Cookie)
	for _, host := range syncHosts {
		u, _ := url.Parse("https://" + host)
		cookieMap[host] = wxwb.apiRuntime.client.Jar.Cookies(u)
	}
	if wxwb.loginStorer != nil {
		storeInfo := storeLoginInfo{
			Cookies:    cookieMap,
			Cookie:     wxwb.loginInfo.cookie,
			SyncKey:    wxwb.loginInfo.syncKey,
			SKey:       wxwb.loginInfo.sKey,
			PassTicket: wxwb.loginInfo.PassTicket,
			User:       wxwb.userInfo.user,
			DeviceID:   wxwb.apiRuntime.deviceID,
			RunInfo:    wxwb.runInfo,
		}
		data, err := json.Marshal(storeInfo)
		if err != nil {
			return err
		}
		err = wxwb.loginStorer.Writer(data)
		return err
	}
	return nil
}

// 从storer中读取信息
func (wxwb *WechatWeb) readLoginInfo() error {
	func() {
		defer func() {
			if r := recover(); r != nil {
				wxwb.logger.Infof("Recovered in readLoginInfo: %v\n", r)
				wxwb.runInfo.PanicCount++
			}
		}()
	}()
	if wxwb.loginStorer != nil {
		data, err := wxwb.loginStorer.Read()
		if err != nil {
			return err
		}
		var storeInfo storeLoginInfo
		err = json.Unmarshal(data, &storeInfo)
		if err != nil {
			return err
		}
		for _, host := range syncHosts {
			u, _ := url.Parse("https://" + host)
			wxwb.apiRuntime.client.Jar.SetCookies(u, storeInfo.Cookies[host])
		}
		wxwb.loginInfo = wechatLoginInfo{
			cookie:     storeInfo.Cookie,
			syncKey:    storeInfo.SyncKey,
			sKey:       storeInfo.SKey,
			PassTicket: storeInfo.PassTicket,
		}
		wxwb.runInfo = storeInfo.RunInfo
		wxwb.userInfo.user = storeInfo.User
		wxwb.apiRuntime.deviceID = storeInfo.DeviceID
	}
	return nil
}
