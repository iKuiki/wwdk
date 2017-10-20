package wxweb

import (
	// "crypto/tls"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"net"
	"net/http"
	"net/http/cookiejar"
	// "net/url"
	"time"
)

// wechatCookie 微信登陆后的cookie凭据，登陆后的消息同步等操作需要此凭据
type wechatCookie struct {
	Wxsid      string
	Wxuin      string // 应该是用户的唯一识别号，同一个用户每次登陆此字段都相同
	Uvid       string
	DataTicket string
	AuthTicket string
	// loadTime   string // 登陆时间(10位时间戳字符串)
}

// WechatWeb 微信网页版客户端实例
type WechatWeb struct {
	cookie      wechatCookie
	userAgent   string
	deviceID    string // 由客户端生成，为e+15位随机数
	contactList []*datastruct.Contact
	user        *datastruct.User
	syncKey     *datastruct.SyncKey
	sKey        string
	PassTicket  string
	messageHook map[datastruct.MessageType][]interface{}
	baseRequest *datastruct.BaseRequest
	client      *http.Client
}

// NewWechatWeb 生成微信网页版客户端实例
func NewWechatWeb() (wxweb WechatWeb, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return WechatWeb{}, err
	}
	return WechatWeb{
		userAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		deviceID:    "e" + tool.GetRandomStringFromNum(15),
		messageHook: make(map[datastruct.MessageType][]interface{}),
		client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 1 * time.Minute,
				}).Dial,
				TLSHandshakeTimeout: 1 * time.Minute,
				// TLSClientConfig: &tls.Config{
				// 	InsecureSkipVerify: true,
				// },
				// Proxy: func(_ *http.Request) (*url.URL, error) {
				// 	return url.Parse("http://127.0.0.1:8888") //根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
				// },
			},
			Jar: jar,
		},
	}, nil
}

// setWechatCookie 为http request设置cookie登陆凭据
func setWechatCookie(request *httplib.BeegoHTTPRequest, cookie wechatCookie) {
	request.SetCookie(&http.Cookie{Name: "wxsid", Value: cookie.Wxsid})
	request.SetCookie(&http.Cookie{Name: "webwx_data_ticket", Value: cookie.DataTicket})
	request.SetCookie(&http.Cookie{Name: "webwxuvid", Value: cookie.Uvid})
	request.SetCookie(&http.Cookie{Name: "webwx_auth_ticket", Value: cookie.AuthTicket})
	request.SetCookie(&http.Cookie{Name: "wxuin", Value: cookie.Wxuin})
}

// GetContact 根据username获取联系人
func (wxwb *WechatWeb) GetContact(username string) (contact *datastruct.Contact, err error) {
	for _, v := range wxwb.contactList {
		if v.UserName == username {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

// GetContactByAlias 根据备注获取联系人
func (wxwb *WechatWeb) GetContactByAlias(alias string) (contact *datastruct.Contact, err error) {
	for _, v := range wxwb.contactList {
		if v.Alias == alias {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

// GetContactByNickname 根据昵称获取用户名
func (wxwb *WechatWeb) GetContactByNickname(nickname string) (contact *datastruct.Contact, err error) {
	for _, v := range wxwb.contactList {
		if v.NickName == nickname {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

// GetContactList 获取联系人列表
func (wxwb *WechatWeb) GetContactList() (contacts []*datastruct.Contact) {
	return wxwb.contactList
}

// refreshCookie 根据response更新cookie
func (wxwb *WechatWeb) refreshCookie(cookies []*http.Cookie) {
	for _, c := range cookies {
		switch c.Name {
		case "wxuin":
			wxwb.cookie.Wxuin = c.Value
		case "wxsid":
			wxwb.cookie.Wxsid = c.Value
		case "webwxuvid":
			wxwb.cookie.Uvid = c.Value
		case "webwx_data_ticket":
			wxwb.cookie.DataTicket = c.Value
		case "webwx_auth_ticket":
			wxwb.cookie.AuthTicket = c.Value
		}
	}
}
