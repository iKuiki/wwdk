package wxweb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
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
	req := httplib.Get("https://login.weixin.qq.com/jslogin")
	req.Param("appid", conf.APP_ID)
	req.Param("fun", "new")
	req.Param("lang", "zh_CN")
	req.Param("_", tool.GetWxTimeStamp())
	resp, err := req.String()
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	ret := tool.AnalysisWxWindowRespond(resp)
	if ret["window.QRLogin.code"] != "200" {
		return "", errors.New("window.QRLogin.code = " + ret["window.QRLogin.code"])
	}
	return ret["window.QRLogin.uuid"], nil
}

// getQrCode 通过uuid生成二维码并输出到控制台
func (wxwb *WechatWeb) getQrCode(uuid string) (err error) {
	req := httplib.Post("https://login.weixin.qq.com/qrcode/" + uuid)
	req.Param("t", "webwx")
	req.Param("_", tool.GetWxTimeStamp())
	_, err = req.String()
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
OutLoop:
	for true {
		req := httplib.Get("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login")
		req.Param("tip", scaned2TipMap[scaned])
		req.Param("uuid", uuid)
		req.Param("_", tool.GetWxTimeStamp())
		resp, err := req.String()
		if err != nil {
			log.Println("waitForScan request error: " + err.Error())
			continue
		}
		ret = tool.AnalysisWxWindowRespond(resp)
		switch ret["window.code"] {
		case "200": // 确认登陆
			log.Println("Login success")
			break OutLoop
		case "201": // 用户已扫码
			scaned = true
			log.Println("Scan success, waiting for login")
			continue
		case "400": // 登陆失败(二维码失效)
			return "", errors.New("Login fail: qrcode has run out")
		case "408": // 等待登陆
			time.Sleep(500 * time.Microsecond)
		default:
			return "", errors.New("Login fail: unknown response code: " + ret["window.code"])
		}
	}
	return ret["window.redirect_uri"], nil
}

func (wxwb *WechatWeb) getCookie(redirectURL, userAgent string) (err error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return errors.New("redirect_url parse fail: " + err.Error())
	}
	query := u.Query()
	req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	req.Param("ticket", query.Get("ticket"))
	req.Param("uuid", query.Get("uuid"))
	req.Param("lang", "zh_CN")
	req.Param("scan", query.Get("scan"))
	req.Param("fun", "new")
	req.SetUserAgent(userAgent)
	resp, err := req.Response()
	if err != nil {
		return errors.New("getCookie request error: " + err.Error())
	}
	cookies := make(map[string]string)
	for _, c := range resp.Cookies() {
		cookies[c.Name] = c.Value
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
	wxwb.cookie = &wechatCookie{
		Wxuin:      cookies["wxuin"],
		Wxsid:      cookies["wxsid"],
		Uvid:       cookies["webwxuvid"],
		DataTicket: cookies["webwx_data_ticket"],
		AuthTicket: cookies["webwx_auth_ticket"],
		PassTicket: bodyResp.PassTicket,
	}
	wxwb.sKey = bodyResp.Skey
	return nil
}

func (wxwb *WechatWeb) wxInit() (err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit")
	body := datastruct.WxInitRequestBody{
		BaseRequest: getBaseRequest(wxwb.cookie, wxwb.deviceID),
	}
	req.Header("Content-Type", "application/json")
	req.Header("charset", "UTF-8")
	req.Param("r", tool.GetWxTimeStamp())
	setWechatCookie(req, wxwb.cookie)
	resp := datastruct.WxInitRespond{}
	data, err := json.Marshal(body)
	if err != nil {
		return errors.New("json.Marshal error: " + err.Error())
	}
	req.Body(data)
	// err = req.ToJSON(&resp)
	r, err := req.Bytes()
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return errors.New("respond json Unmarshal to struct fail: " + err.Error())
	}
	if resp.BaseResponse.Ret != 0 {
		return fmt.Errorf("respond ret error: %d", resp.BaseResponse.Ret)
	}
	wxwb.user = resp.User
	wxwb.syncKey = resp.SyncKey
	// wxwb.sKey = resp.SKey
	return nil
}

func (wxwb *WechatWeb) getContactList() (err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact")
	req.Param("r", tool.GetWxTimeStamp())
	setWechatCookie(req, wxwb.cookie)
	req.Body([]byte("{}"))
	resp := datastruct.GetContactRespond{}
	r, err := req.Bytes()
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return errors.New("respond json Unmarshal to struct fail: " + err.Error())
	}
	if resp.BaseResponse.Ret != 0 {
		return fmt.Errorf("respond ret error: %d", resp.BaseResponse.Ret)
	}
	wxwb.contactList = resp.MemberList
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
	err = wxwb.getCookie(redirectURL, wxwb.userAgent)
	if err != nil {
		return errors.New("getCookie error: " + err.Error())
	}
	err = wxwb.wxInit()
	if err != nil {
		return errors.New("wxInit error: " + err.Error())
	}
	err = wxwb.getContactList()
	if err != nil {
		return errors.New("getContactList error: " + err.Error())
	}
	log.Printf("User %s has Login Success, total %d contacts\n", wxwb.user.NickName, len(wxwb.contactList))
	return nil
}
