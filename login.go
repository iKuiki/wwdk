package wxweb

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ikuiki/wwdk/conf"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"github.com/mdp/qrterminal"
)

func (wxwb *WechatWeb) getUUID() (uuid string, err error) {
	params := url.Values{}
	params.Set("appid", conf.AppID)
	params.Set("fun", "new")
	params.Set("lang", conf.Lang)
	params.Set("_", tool.GetWxTimeStamp())
	req, _ := http.NewRequest("GET", "https://login.weixin.qq.com/jslogin?"+params.Encode(), nil)
	resp, err := wxwb.request(req)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ret := tool.AnalysisWxWindowRespond(string(body))
	if ret["window.QRLogin.code"] != "200" {
		return "", errors.New("window.QRLogin.code = " + ret["window.QRLogin.code"])
	}
	return ret["window.QRLogin.uuid"], nil
}

// getQrCode 通过uuid生成二维码并输出到控制台
func (wxwb *WechatWeb) getQrCode(uuid string) (err error) {
	if os.Getenv("DEBUG_PRINT_QRURL") == "true" {
		log.Println("qrcode url: https://login.weixin.qq.com/l/" + uuid)
	}
	qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)
	return nil
}

// waitForScan 等待用户扫描二维码登陆
// 当扫码超时、扫码失败时，应当从getUUID方法重新开始
func (wxwb *WechatWeb) waitForScan(uuid string) (redirectURL string, err error) {
	var ret map[string]string
	scaned := false
	scaned2TipMap := map[bool]string{
		false: "1",
		true:  "0",
	}
	for true {
		redirectURL, err = func() (redirectURL string, err error) {
			params := url.Values{}
			params.Set("tip", scaned2TipMap[scaned])
			params.Set("uuid", uuid)
			params.Set("_", tool.GetWxTimeStamp())
			req, _ := http.NewRequest(`GET`, "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?"+params.Encode(), nil)
			resp, err := wxwb.request(req)
			if err != nil {
				log.Println("waitForScan request error: " + err.Error())
				return "", nil // return nil error for continue
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			ret = tool.AnalysisWxWindowRespond(string(body))
			switch ret["window.code"] {
			case "200": // 确认登陆
				log.Println("Login success")
				return ret["window.redirect_uri"], nil
			case "201": // 用户已扫码
				scaned = true
				log.Println("Scan success, waiting for login")
				return "", nil // continue
			case "400": // 登陆失败(二维码失效)
				return "", errors.New("Login fail: qrcode has run out")
			case "408": // 等待登陆
				time.Sleep(500 * time.Microsecond)
			default:
				return "", errors.New("Login fail: unknown response code: " + ret["window.code"])
			}
			return "", nil
		}()
		if err != nil || redirectURL != "" {
			break
		}
	}
	return redirectURL, err
}

func (wxwb *WechatWeb) getCookie(redirectURL string) (err error) {

	req, _ := http.NewRequest(`GET`, redirectURL+"&fun=new", nil)
	resp, err := wxwb.request(req)
	if err != nil {
		return errors.New("getCookie request error: " + err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Read respond body error: " + err.Error())
	}
	var bodyResp datastruct.GetCookieRespond
	err = xml.Unmarshal(body, &bodyResp)
	if err != nil {
		return errors.New("Unmarshal respond xml error: " + err.Error())
	}
	// wxwb.refreshCookie(resp.Cookies()) // 在统一的Request方法已经调用了
	wxwb.loginInfo.PassTicket = bodyResp.PassTicket
	wxwb.loginInfo.sKey = bodyResp.Skey
	return nil
}

func (wxwb *WechatWeb) wxInit() (err error) {
	data, err := json.Marshal(datastruct.WxInitRequestBody{
		BaseRequest: wxwb.baseRequest(),
	})
	if err != nil {
		return errors.New("json.Marshal error: " + err.Error())
	}
	params := url.Values{}
	params.Set("pass_ticket", wxwb.loginInfo.PassTicket)
	params.Set("skey", wxwb.loginInfo.sKey)
	params.Set("r", tool.GetWxTimeStamp())
	// resp, err := wxwb.client.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?"+params.Encode(),
	// 	"application/json;charset=UTF-8",
	// 	bytes.NewReader(data))

	req, err := http.NewRequest("POST",
		"https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?"+params.Encode(),
		bytes.NewReader(data))
	if err != nil {
		return errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.request(req)
	if err != nil {
		return errors.New("do request error: " + err.Error())
	}
	defer resp.Body.Close()

	respStruct := datastruct.WxInitRespond{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return errors.New("respond json Unmarshal to struct fail: " + err.Error())
	}
	if respStruct.BaseResponse.Ret != 0 {
		return fmt.Errorf("respond ret error: %d", respStruct.BaseResponse.Ret)
	}
	for _, contact := range respStruct.ContactList {
		wxwb.contactList[contact.UserName] = contact
	}
	wxwb.user = respStruct.User
	wxwb.loginInfo.syncKey = respStruct.SyncKey
	wxwb.loginInfo.sKey = respStruct.SKey
	return nil
}

func (wxwb *WechatWeb) getContactList() (err error) {
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	resp, err := wxwb.client.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?" + params.Encode())
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
		return fmt.Errorf("respond ret error: %d", respStruct.BaseResponse.Ret)
	}
	for _, contact := range respStruct.MemberList {
		wxwb.contactList[contact.UserName] = contact
	}
	return nil
}

func (wxwb *WechatWeb) getBatchContact() (err error) {
	dataStruct := datastruct.GetBatchContactRequest{
		BaseRequest: wxwb.baseRequest(),
	}
	for _, contact := range wxwb.contactList {
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
	resp, err := wxwb.client.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?"+params.Encode(),
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
		return fmt.Errorf("respond ret error: %d", respStruct.BaseResponse.Ret)
	}
	for _, contact := range respStruct.ContactList {
		if c, ok := wxwb.contactList[contact.UserName]; ok {
			c.MemberCount = contact.MemberCount
			c.MemberList = contact.MemberList
			c.EncryChatRoomID = contact.EncryChatRoomID
			wxwb.contactList[c.UserName] = c
		}
	}
	return
}

// Login 登陆方法总成
func (wxwb *WechatWeb) Login() (err error) {
	wxwb.readLoginInfo()
	err = wxwb.getContactList()
	if err != nil {
		log.Println("stored login info not avaliable")
		wxwb.resetLoginInfo()
		uuid, err := wxwb.getUUID()
		if err != nil {
			return errors.New("Get UUID fail: " + err.Error())
		}
		err = wxwb.getQrCode(uuid)
		if err != nil {
			return errors.New("Get QrCode fail: " + err.Error())
		}
		redirectURL, err := wxwb.waitForScan(uuid)
		if err != nil {
			return errors.New("waitForScan error: " + err.Error())
		}
		// panic(redirectUrl)
		err = wxwb.getCookie(redirectURL)
		if err != nil {
			return errors.New("getCookie error: " + err.Error())
		}
		err = wxwb.wxInit()
		if err != nil {
			return errors.New("wxInit error: " + err.Error())
		}
		// 此处即认为登陆成功
		wxwb.runInfo.LoginAt = time.Now()
	} else {
		log.Printf("reuse loginInfo [%s] logined at %v\n", wxwb.user.NickName, wxwb.runInfo.LoginAt.Format("2006-01-02 15:04:05"))
	}
	// err = wxwb.StatusNotify(wxwb.user.UserName, wxwb.user.UserName, 3)
	// if err != nil {
	// 	return errors.New("StatusNotify error: " + err.Error())
	// }
	err = wxwb.getContactList()
	if err != nil {
		return errors.New("getContactList error: " + err.Error())
	}
	err = wxwb.getBatchContact()
	if err != nil {
		return errors.New("getBatchContact error: " + err.Error())
	}
	log.Printf("User %s has Login Success, total %d contacts\n", wxwb.user.NickName, len(wxwb.contactList))
	// 如有必要，记录login信息到storer
	wxwb.writeLoginInfo()
	return nil
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
	resp, err := wxwb.client.PostForm("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxlogout?"+params.Encode(), form)
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	return nil
}
