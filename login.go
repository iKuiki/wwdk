package wxweb

import (
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/mdp/qrterminal"
	"github.com/yinhui87/wechat-web/conf"
	"github.com/yinhui87/wechat-web/tool"
	"log"
	"net/url"
	"os"
	"time"
)

func getUUID() (uuid string, err error) {
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

func getQrCode(uuid string) (err error) {
	req := httplib.Post("https://login.weixin.qq.com/qrcode/" + uuid)
	req.Param("t", "webwx")
	req.Param("_", tool.GetWxTimeStamp())
	_, err = req.String()
	if err != nil {
		return err
	}
	qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)
	return nil
}

func waitForScan(uuid string) (redirectUrl string, err error) {
	var ret map[string]string
	req := httplib.Get("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login")
	req.Param("tip", "1")
	req.Param("uuid", uuid)
	req.Param("_", tool.GetWxTimeStamp())
	_, err = req.String()
	if err != nil {
		return "", errors.New("waitForScan request error: " + err.Error())
	}
	log.Println("Scan success, waiting for login")
	for true {
		req := httplib.Get("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login")
		req.Param("tip", "1")
		req.Param("uuid", uuid)
		req.Param("_", tool.GetWxTimeStamp())
		resp, err := req.String()
		if err != nil {
			log.Println("waitForScan request error: " + err.Error())
			continue
		}
		ret = tool.AnalysisWxWindowRespond(resp)
		if ret["window.code"] != "200" {
			time.Sleep(500 * time.Microsecond)
			continue
		}
		break
	}
	return ret["window.redirect_uri"], nil
}

func getCookie(redirectUrl string) (cookie wechatCookie, err error) {
	u, err := url.Parse(redirectUrl)
	if err != nil {
		return wechatCookie{}, errors.New("redirect_url parse fail: " + err.Error())
	}
	query := u.Query()
	req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	req.Param("ticket", query.Get("ticket"))
	req.Param("uuid", query.Get("uuid"))
	req.Param("lang", "zh_CN")
	req.Param("scan", query.Get("scan"))
	req.Param("fun", "new")
	req.SetUserAgent(conf.USER_AGENT)
	resp, err := req.Response()
	if err != nil {
		return wechatCookie{}, errors.New("getCookie request error: " + err.Error())
	}
	cookies := make(map[string]string)
	for _, c := range resp.Cookies() {
		cookies[c.Name] = c.Value
	}
	cookie = wechatCookie{
		Wxuin:      cookies["wxuin"],
		Wxsid:      cookies["wxsid"],
		Uvid:       cookies["webwxuvid"],
		DataTicket: cookies["webwx_data_ticket"],
		AuthTicket: cookies["webwx_auth_ticket"],
	}
	return cookie, nil
}

func (this *WechatWeb) Login() (err error) {
	uuid, err := getUUID()
	if err != nil {
		return errors.New("Get UUID fail: " + err.Error())
	}
	err = getQrCode(uuid)
	if err != nil {
		return errors.New("Get QrCode fail: " + err.Error())
	}
	redirectUrl, err := waitForScan(uuid)
	if err != nil {
		return errors.New("waitForScan error: " + err.Error())
	}
	cookie, err := getCookie(redirectUrl)
	if err != nil {
		return errors.New("getCookie error: " + err.Error())
	}
	this.cookie = cookie
	return nil
}
