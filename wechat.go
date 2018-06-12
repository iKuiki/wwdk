package wxweb

import (
	// "crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/cookiejar"

	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
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
	cookie         wechatCookie
	userAgent      string
	deviceID       string // 由客户端生成，为e+15位随机数
	contactList    map[string]datastruct.Contact
	user           *datastruct.User
	syncHost       string
	syncKey        *datastruct.SyncKey
	sKey           string
	PassTicket     string
	messageHook    map[datastruct.MessageType][]interface{}
	modContactHook []interface{}
	client         *http.Client
}

// NewWechatWeb 生成微信网页版客户端实例
func NewWechatWeb() (wxweb *WechatWeb, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &WechatWeb{}, err
	}
	return &WechatWeb{
		contactList: make(map[string]datastruct.Contact),
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
			Jar:     jar,
			Timeout: 1 * time.Minute,
		},
	}, nil
}

func (wxwb *WechatWeb) baseRequest() (baseRequest *datastruct.BaseRequest) {
	return &datastruct.BaseRequest{
		Uin:      wxwb.cookie.Wxuin,
		Sid:      wxwb.cookie.Wxsid,
		Skey:     wxwb.sKey,
		DeviceID: wxwb.deviceID,
	}
}

// GetContact 根据username获取联系人
func (wxwb *WechatWeb) GetContact(username string) (contact datastruct.Contact, err error) {
	contact, ok := wxwb.contactList[username]
	if !ok {
		err = errors.New("User not found")
	}
	return
}

// GetContactByAlias 根据备注获取联系人
func (wxwb *WechatWeb) GetContactByAlias(alias string) (contact datastruct.Contact, err error) {
	found := false
	for _, v := range wxwb.contactList {
		if v.Alias == alias {
			contact = v
			found = true
		}
	}
	if !found {
		err = errors.New("User not found")
	}
	return
}

// GetContactByNickname 根据昵称获取用户名
func (wxwb *WechatWeb) GetContactByNickname(nickname string) (contact datastruct.Contact, err error) {
	found := false
	for _, v := range wxwb.contactList {
		if v.NickName == nickname {
			contact = v
			found = true
		}
	}
	if !found {
		err = errors.New("User not found")
	}
	return
}

// GetContactList 获取联系人列表
func (wxwb *WechatWeb) GetContactList() (contacts []datastruct.Contact) {
	for _, v := range wxwb.contactList {
		contacts = append(contacts, v)
	}
	return
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

// 统一请求
func (wxwb *WechatWeb) request(req *http.Request) (resp *http.Response, err error) {
	resp, err = wxwb.client.Do(req)
	wxwb.refreshCookie(resp.Cookies())
	return
}
