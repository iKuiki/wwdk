package wwdk

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/ikuiki/wwdk/conf"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
)

// LoginChannelItem 登录时channel内返回的东西
type LoginChannelItem struct {
	Err  error
	Code LoginStatus
	Msg  string
}

// LoginStatus 登录状态枚举
type LoginStatus int32

const (
	// LoginStatusErrorOccurred 发生异常
	LoginStatusErrorOccurred LoginStatus = -1
	// LoginStatusWaitForScan 等待扫码
	// 返回Msg: 待扫码url
	LoginStatusWaitForScan LoginStatus = 1
	// LoginStatusScanedWaitForLogin 用户已经扫码
	// 返回Msg: 用户头像的base64
	LoginStatusScanedWaitForLogin LoginStatus = 2
	// LoginStatusScanedFinish 用户已同意登陆
	LoginStatusScanedFinish LoginStatus = 3
	// LoginStatusGotCookie 已获取到Cookie
	LoginStatusGotCookie LoginStatus = 4
	// LoginStatusInitFinish 登陆初始化完成
	LoginStatusInitFinish LoginStatus = 5
	// LoginStatusGotContact 已获取到联系人
	LoginStatusGotContact LoginStatus = 6
	// LoginStatusGotBatchContact 已获取到群聊成员
	LoginStatusGotBatchContact LoginStatus = 7
)

func (wxwb *WechatWeb) getUUID(loginChannel chan<- LoginChannelItem) (uuid string) {
	params := url.Values{}
	params.Set("appid", conf.AppID)
	params.Set("fun", "new")
	params.Set("lang", conf.Lang)
	params.Set("_", tool.GetWxTimeStamp())
	req, _ := http.NewRequest("GET", "https://login.weixin.qq.com/jslogin?"+params.Encode(), nil)
	resp, err := wxwb.request(req)
	if err != nil {
		panic(errors.New("request error: " + err.Error()))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ret := tool.ExtractWxWindowRespond(string(body))
	if ret["window.QRLogin.code"] != "200" {
		panic(errors.New("window.QRLogin.code = " + ret["window.QRLogin.code"]))
	}
	uuid = ret["window.QRLogin.uuid"]
	loginChannel <- LoginChannelItem{
		Code: LoginStatusWaitForScan,
		Msg:  "https://login.weixin.qq.com/l/" + uuid,
	}
	return uuid
}

// waitForScan 等待用户扫描二维码登陆
// 当扫码超时、扫码失败时，应当从getUUID方法重新开始
func (wxwb *WechatWeb) waitForScan(uuid string, loginChannel chan<- LoginChannelItem) (redirectURL string) {
	var ret map[string]string
	tip := "1"
	for true {
		redirectURL = func() (redirectURL string) {
			params := url.Values{}
			params.Set("tip", tip)
			tip = "0" // 在第二次轮询的时候tip就为0了
			params.Set("uuid", uuid)
			params.Set("_", tool.GetWxTimeStamp())
			req, _ := http.NewRequest(`GET`, "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?"+params.Encode(), nil)
			resp, err := wxwb.request(req)
			if err != nil {
				wxwb.logger.Infof("waitForScan request error: %v\n", err)
				return "" // return empty for continue
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			ret = tool.ExtractWxWindowRespond(string(body))
			switch ret["window.code"] {
			case "200": // 确认登陆
				wxwb.logger.Info("Login success\n")
				loginChannel <- LoginChannelItem{
					Code: LoginStatusScanedFinish,
				}
				return ret["window.redirect_uri"]
			case "201": // 用户已扫码
				wxwb.logger.Info("Scan success, waiting for login\n")
				loginChannel <- LoginChannelItem{
					Code: LoginStatusScanedWaitForLogin,
					Msg:  ret["window.userAvatar"],
				}
				return "" // continue
			case "400": // 登陆失败(二维码失效)
				panic(errors.New("Login fail: qrcode has run out"))
			case "408": // 等待登陆
				time.Sleep(500 * time.Microsecond)
			default:
				panic(errors.New("Login fail: unknown response code: " + ret["window.code"]))
			}
			return
		}()
		if redirectURL != "" {
			break
		}
	}
	return redirectURL
}

func (wxwb *WechatWeb) getCookie(redirectURL string, loginChannel chan<- LoginChannelItem) {
	req, _ := http.NewRequest(`GET`, redirectURL+"&fun=new&version=v2", nil)
	resp, err := wxwb.request(req)
	if err != nil {
		panic(errors.New("getCookie request error: " + err.Error()))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errors.New("Read respond body error: " + err.Error()))
	}
	var bodyResp datastruct.GetCookieRespond
	err = xml.Unmarshal(body, &bodyResp)
	if err != nil {
		panic(errors.New("Unmarshal respond xml error: " + err.Error()))
	}
	loginChannel <- LoginChannelItem{
		Code: LoginStatusGotCookie,
	}
	wxwb.loginInfo.sKey = bodyResp.Skey
	wxwb.loginInfo.PassTicket = bodyResp.PassTicket
}

func (wxwb *WechatWeb) wxInit(loginChannel chan<- LoginChannelItem) {
	data, err := json.Marshal(datastruct.WxInitRequestBody{
		BaseRequest: wxwb.baseRequest(),
	})
	if err != nil {
		panic(errors.New("json.Marshal error: " + err.Error()))
	}
	params := url.Values{}
	params.Set("pass_ticket", wxwb.loginInfo.PassTicket)
	// params.Set("skey", wxwb.loginInfo.sKey)
	params.Set("r", tool.GetWxTimeStamp())
	// resp, err := wxwb.apiRuntime.client.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?"+params.Encode(),
	// 	"application/json;charset=UTF-8",
	// 	bytes.NewReader(data))

	req, err := http.NewRequest("POST",
		"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?"+params.Encode(),
		bytes.NewReader(data))
	if err != nil {
		panic(errors.New("create request error: " + err.Error()))
	}
	resp, err := wxwb.request(req)
	if err != nil {
		panic(errors.New("do request error: " + err.Error()))
	}
	defer resp.Body.Close()

	respStruct := datastruct.WxInitRespond{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		panic(errors.New("respond json Unmarshal to struct fail: " + err.Error()))
	}
	if respStruct.BaseResponse.Ret != 0 {
		panic(errors.Errorf("respond ret(%d) error: %s", respStruct.BaseResponse.Ret, string(body)))
	}
	loginChannel <- LoginChannelItem{
		Code: LoginStatusInitFinish,
	}
	for _, contact := range respStruct.ContactList {
		wxwb.userInfo.contactList[contact.UserName] = contact
	}
	wxwb.userInfo.user = respStruct.User
	wxwb.loginInfo.syncKey = respStruct.SyncKey
	wxwb.loginInfo.sKey = respStruct.SKey
}

// 获取联系人
// 注：坑！此处获取到的居然不是完整的联系人，必须和init中获取到的合并后才是完整的联系人列表
func (wxwb *WechatWeb) getContactList() (err error) {
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	resp, err := wxwb.apiRuntime.client.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?" + params.Encode())
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	respStruct := datastruct.GetContactRespond{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return errors.New("respond json Unmarshal to struct fail: " + err.Error())
	}
	if respStruct.BaseResponse.Ret != 0 {
		return errors.Errorf("respond ret error(%d): %s", respStruct.BaseResponse.Ret, string(body))
	}
	for _, contact := range respStruct.MemberList {
		wxwb.userInfo.contactList[contact.UserName] = contact
	}
	return nil
}

// 获取群聊的成员
func (wxwb *WechatWeb) getBatchContact() (err error) {
	dataStruct := datastruct.GetBatchContactRequest{
		BaseRequest: wxwb.baseRequest(),
	}
	for _, contact := range wxwb.userInfo.contactList {
		if contact.IsChatroom() {
			dataStruct.List = append(dataStruct.List, datastruct.GetBatchContactRequestListItem{
				UserName: contact.UserName,
			})
		}
	}
	dataStruct.Count = int64(len(dataStruct.List))
	if dataStruct.Count == 0 {
		return nil
	}
	data, err := json.Marshal(dataStruct)
	if err != nil {
		return errors.New("json.Marshal error: " + err.Error())
	}
	params := url.Values{}
	params.Set("type", "ex")
	params.Set("r", tool.GetWxTimeStamp())
	resp, err := wxwb.apiRuntime.client.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?"+params.Encode(),
		"application/json;charset=UTF-8",
		bytes.NewReader(data))
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	respStruct := datastruct.GetBatchContactResponse{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return errors.New("respond json Unmarshal to struct fail: " + err.Error())
	}
	if respStruct.BaseResponse.Ret != 0 {
		return errors.Errorf("respond ret error(%d): %s", respStruct.BaseResponse.Ret, string(body))
	}
	for _, contact := range respStruct.ContactList {
		if c, ok := wxwb.userInfo.contactList[contact.UserName]; ok {
			c.MemberCount = contact.MemberCount
			c.MemberList = contact.MemberList
			c.EncryChatRoomID = contact.EncryChatRoomID
			wxwb.userInfo.contactList[c.UserName] = c
		}
	}
	return
}

// Login 登陆方法总成
// param loginChannel 登陆状态channel，从中可以读取到登录情况
func (wxwb *WechatWeb) Login(loginChannel chan<- LoginChannelItem) {
	// 尝试使用已存在的登录信息登录
	go func() {
		defer close(loginChannel)
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					// 发生了panic
					loginChannel <- LoginChannelItem{
						Code: LoginStatusErrorOccurred,
						Err:  err,
					}
				} else {
					wxwb.logger.Errorf("WechatWeb.Login panic: \n%+v\n", e)
				}
			}
		}()
		// 尝试使用已存在的登录信息登录
		logined := false
		readed, _ := wxwb.readLoginInfo()
		if readed {
			wxwb.logger.Info("loaded stored login info")
			if wxwb.getContactList() == nil {
				// 获取联系人成功，则为已登陆状态
				logined = true
				wxwb.logger.Infof("reuse loginInfo [%s] logined at %v\n", wxwb.userInfo.user.NickName, wxwb.runInfo.LoginAt.Format("2006-01-02 15:04:05"))
			}
		}
		if !logined {
			if readed {
				// 仅当成功读取了login信息并且登陆失败，才输出此log
				wxwb.logger.Info("stored login info not avaliable\n")
			}
			wxwb.resetLoginInfo()
			uuid := wxwb.getUUID(loginChannel)
			redirectURL := wxwb.waitForScan(uuid, loginChannel)
			// panic(redirectUrl)
			wxwb.getCookie(redirectURL, loginChannel)
			wxwb.wxInit(loginChannel)
			err := wxwb.getContactList()
			if err != nil {
				loginChannel <- LoginChannelItem{
					Code: LoginStatusErrorOccurred,
					Err:  err,
				}
				return
			}
			// 此处即认为登陆成功
			wxwb.runInfo.LoginAt = time.Now()
		}
		loginChannel <- LoginChannelItem{
			Code: LoginStatusGotContact,
		}
		// err = wxwb.StatusNotify(wxwb.userInfo.user.UserName, wxwb.userInfo.user.UserName, 3)
		// if err != nil {
		// 	return errors.New("StatusNotify error: " + err.Error())
		// }
		err := wxwb.getBatchContact()
		if err != nil {
			loginChannel <- LoginChannelItem{
				Code: LoginStatusErrorOccurred,
				Err:  err,
			}
			return
		}
		loginChannel <- LoginChannelItem{
			Code: LoginStatusGotBatchContact,
		}
		wxwb.logger.Infof("User %s has Login Success, total %d contacts\n", wxwb.userInfo.user.NickName, len(wxwb.userInfo.contactList))
		// 如有必要，记录login信息到storer
		wxwb.writeLoginInfo()
	}()
}

// Logout 退出登录
func (wxwb *WechatWeb) Logout() (err error) {
	params := url.Values{}
	params.Set("redirect", "0")
	params.Set("type", "1")
	params.Set("skey", wxwb.loginInfo.sKey)
	form := url.Values{}
	form.Set("sid", wxwb.loginInfo.cookie.Wxsid)
	form.Set("uin", wxwb.loginInfo.cookie.Wxuin)
	resp, err := wxwb.apiRuntime.client.PostForm("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxlogout?"+params.Encode(), form)
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	return nil
}
