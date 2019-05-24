package wwdk

// 此文件中存储的为WechatWeb读写登录凭据的方法

import (
	"encoding/json"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/pkg/errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// storeLoginInfo 用于储存的登录信息
type storeLoginInfo struct {
	Cookies     map[string][]*http.Cookie
	Cookie      wechatCookie
	SyncKey     *datastruct.SyncKey
	SKey        string
	PassTicket  string
	RunInfo     WechatRunInfo    // 运行统计信息
	DeviceID    string           // 由客户端生成，为e+15位随机数
	User        *datastruct.User // 用户信息
	ContactList map[string]datastruct.Contact
}

// 重置登录信息
func (wxwb *WechatWeb) resetLoginInfo() (err error) {
	defer func() {
		if r := recover(); r != nil {
			wxwb.logger.Infof("Recovered in resetLoginInfo: %v\n", r)
			wxwb.runInfo.PanicCount++
			err = errors.Errorf("panic recovered: %+v", r)
		}
	}()
	if wxwb.loginStorer != nil {
		wxwb.loginStorer.Truncate()
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return errors.WithStack(err)
	}
	wxwb.apiRuntime.client.Jar = jar
	wxwb.loginInfo = wechatLoginInfo{}
	return nil
}

// 往storer中写入信息
func (wxwb *WechatWeb) writeLoginInfo() (err error) {
	defer func() {
		if r := recover(); r != nil {
			wxwb.logger.Infof("Recovered in writeLoginInfo: %v\n", r)
			wxwb.runInfo.PanicCount++
			err = errors.Errorf("panic recovered: %+v", r)
		}
	}()
	cookieMap := make(map[string][]*http.Cookie)
	for _, host := range syncHosts {
		u, _ := url.Parse("https://" + host)
		cookieMap[host] = wxwb.apiRuntime.client.Jar.Cookies(u)
	}
	if wxwb.loginStorer != nil {
		storeInfo := storeLoginInfo{
			Cookies:     cookieMap,
			Cookie:      wxwb.loginInfo.cookie,
			SyncKey:     wxwb.loginInfo.syncKey,
			SKey:        wxwb.loginInfo.sKey,
			PassTicket:  wxwb.loginInfo.PassTicket,
			User:        wxwb.userInfo.user,
			ContactList: wxwb.userInfo.contactList,
			DeviceID:    wxwb.apiRuntime.deviceID,
			RunInfo:     wxwb.runInfo,
		}
		data, err := json.Marshal(storeInfo)
		if err != nil {
			return errors.WithStack(err)
		}
		err = wxwb.loginStorer.Write(data)
		return errors.WithStack(err)
	}
	return nil
}

// 从storer中读取信息
// 返回是否成功读取到信息
func (wxwb *WechatWeb) readLoginInfo() (readed bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			wxwb.logger.Infof("Recovered in readLoginInfo: %v\n", r)
			wxwb.runInfo.PanicCount++
			err = errors.Errorf("panic recovered: %+v", r)
		}
	}()
	if wxwb.loginStorer != nil {
		data, err := wxwb.loginStorer.Read()
		if err != nil {
			return false, errors.WithStack(err)
		}
		var storeInfo storeLoginInfo
		err = json.Unmarshal(data, &storeInfo)
		if err != nil {
			return false, errors.WithStack(err)
		}
		if storeInfo.DeviceID == "" {
			// 未读取到信息
			return false, nil
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
		{
			// 先暂存StartAt，对StartAt不做覆盖
			started := wxwb.runInfo.StartAt
			wxwb.runInfo = storeInfo.RunInfo
			// 还原startat
			wxwb.runInfo.StartAt = started
		}
		if storeInfo.User == nil {
			return false, nil
		}
		wxwb.userInfo.user = storeInfo.User
		for _, contact := range storeInfo.ContactList {
			wxwb.userInfo.contactList[contact.UserName] = contact
		}
		wxwb.apiRuntime.deviceID = storeInfo.DeviceID
		// 读取到了信息并进行还原
		return true, nil
	}
	return false, nil
}
