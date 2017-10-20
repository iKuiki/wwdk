package wxweb

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/mdp/qrterminal"
	"github.com/yinhui87/wechat-web/conf"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"
)

func (wxwb *WechatWeb) getUUID() (uuid string, err error) {
	// req := httplib.Get("https://login.weixin.qq.com/jslogin")
	// req.Param("appid", conf.APP_ID)
	// req.Param("fun", "new")
	// req.Param("lang", conf.Lang)
	// req.Param("_", tool.GetWxTimeStamp())
	// resp, err := req.String()
	// if err != nil {
	// 	return "", errors.New("request error: " + err.Error())
	// }
	params := url.Values{}
	params.Set("appid", conf.AppID)
	params.Set("fun", "new")
	params.Set("lang", conf.Lang)
	params.Set("_", tool.GetWxTimeStamp())
	resp, err := wxwb.client.Get("https://login.weixin.qq.com/jslogin?" + params.Encode())
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
	// req := httplib.Get("https://login.weixin.qq.com/qrcode/" + uuid)
	// req.Param("t", "webwx")
	// req.Param("_", tool.GetWxTimeStamp())
	// _, err = req.String()
	params := url.Values{}
	params.Set("t", "webwx")
	params.Set("_", tool.GetWxTimeStamp())
	_, err = wxwb.client.Get("https://login.weixin.qq.com/qrcode/" + uuid + "?" + params.Encode())
	if err != nil {
		return err
	}
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
		// req := httplib.Get("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login")
		// req.Param("tip", scaned2TipMap[scaned])
		// req.Param("uuid", uuid)
		// req.Param("_", tool.GetWxTimeStamp())
		// resp, err := req.String()
		redirectURL, err = func() (redirectURL string, err error) {
			params := url.Values{}
			params.Set("tip", scaned2TipMap[scaned])
			params.Set("uuid", uuid)
			params.Set("_", tool.GetWxTimeStamp())
			resp, err := wxwb.client.Get("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?" + params.Encode())
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
	// u, err := url.Parse(redirectURL)
	// if err != nil {
	// 	return errors.New("redirect_url parse fail: " + err.Error())
	// }
	// query := u.Query()
	// req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	// req.Param("ticket", query.Get("ticket"))
	// req.Param("uuid", query.Get("uuid"))
	// req.Param("lang", conf.Lang)
	// req.Param("scan", query.Get("scan"))
	// req.Param("fun", "new")
	// req.SetUserAgent(wxwb.userAgent)
	// resp, err := req.Response()
	resp, err := wxwb.client.Get(redirectURL + "&fun=new")
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
	wxwb.refreshCookie(resp.Cookies())
	wxwb.PassTicket = bodyResp.PassTicket
	wxwb.sKey = bodyResp.Skey
	wxwb.baseRequest = &datastruct.BaseRequest{
		Uin:      wxwb.cookie.Wxuin,
		Sid:      wxwb.cookie.Wxsid,
		Skey:     wxwb.sKey,
		DeviceID: wxwb.deviceID,
	}
	return nil
}

func (wxwb *WechatWeb) wxInit() (err error) {
	// req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit")
	// req.Header("Content-Type", "application/json;charset=UTF-8")
	// req.Param("r", tool.GetWxTimeStamp())
	// setWechatCookie(req, wxwb.cookie)
	// req.Body(data)
	// err = req.ToJSON(&respStruct)
	// r, err := req.Bytes()
	data, err := json.Marshal(datastruct.WxInitRequestBody{
		BaseRequest: wxwb.baseRequest,
	})
	if err != nil {
		return errors.New("json.Marshal error: " + err.Error())
	}
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	resp, err := wxwb.client.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?"+params.Encode(),
		"application/json;charset=UTF-8",
		bytes.NewReader(data))

	respStruct := datastruct.WxInitRespond{}
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return errors.New("respond json Unmarshal to struct fail: " + err.Error())
	}
	if respStruct.BaseResponse.Ret != 0 {
		return fmt.Errorf("respond ret error: %d", respStruct.BaseResponse.Ret)
	}
	wxwb.user = respStruct.User
	wxwb.syncKey = respStruct.SyncKey
	// wxwb.sKey = respStruct.SKey
	return nil
}

func (wxwb *WechatWeb) getContactList() (err error) {
	// req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact")
	// req.Param("r", tool.GetWxTimeStamp())
	// setWechatCookie(req, wxwb.cookie)
	// r, err := req.Bytes()
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
	wxwb.contactList = respStruct.MemberList
	return nil
}

// Login 登陆方法总成
func (wxwb *WechatWeb) Login() (err error) {
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
	err = wxwb.StatusNotify(wxwb.user.UserName, wxwb.user.UserName, 3)
	if err != nil {
		return errors.New("StatusNotify error: " + err.Error())
	}
	err = wxwb.getContactList()
	if err != nil {
		return errors.New("getContactList error: " + err.Error())
	}
	log.Printf("User %s has Login Success, total %d contacts\n", wxwb.user.NickName, len(wxwb.contactList))
	return nil
}
