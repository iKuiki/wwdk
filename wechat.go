package wwdk

import (
	"encoding/json"
	"fmt"
	"github.com/ikuiki/wwdk/storer"
	"log"
	// "crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/cookiejar"

	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"net/url"
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

type wechatLoginInfo struct {
	cookie     wechatCookie
	syncKey    *datastruct.SyncKey
	sKey       string
	PassTicket string
}

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

// WechatRunInfo 微信运行信息
type WechatRunInfo struct {
	// StartAt 程序启动的时间
	StartAt time.Time
	// LoginAt 程序登陆的时间
	LoginAt time.Time
	// SyncCount 同步次数
	SyncCount uint64
	// ContactModifyCount 联系人修改计数器
	ContactModifyCount uint64
	// MessageCount 消息计数器
	MessageCount uint64
	// MessageRecivedCount 收到消息计数器
	MessageRecivedCount uint64
	// MessageSentCount 发送消息计数器
	MessageSentCount uint64
	// MessageRevokeCount 撤回消息计数器
	MessageRevokeCount uint64
	// MessageRevokeRecivedCount 收到撤回消息计数器
	MessageRevokeRecivedCount uint64
	// MessageRevokeSentCount 发送撤回消息计数器
	MessageRevokeSentCount uint64
	// PanicCount panic计数器
	PanicCount uint64
}

// WechatWeb 微信网页版客户端实例
type WechatWeb struct {
	userAgent      string
	contactList    map[string]datastruct.Contact
	user           *datastruct.User
	syncHost       string
	messageHook    map[datastruct.MessageType][]interface{}
	modContactHook []interface{}
	client         *http.Client
	loginInfo      wechatLoginInfo // 登陆信息
	runInfo        WechatRunInfo   // 运行统计信息
	deviceID       string          // 由客户端生成，为e+15位随机数
	loginStorer    storer.Storer   // 存储器，如果有赋值，则用于记录登录信息
}

// NewWechatWeb 生成微信网页版客户端实例
func NewWechatWeb(configs ...interface{}) (wxweb *WechatWeb, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &WechatWeb{}, err
	}
	w := &WechatWeb{
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
		runInfo: WechatRunInfo{
			StartAt: time.Now(),
		},
	}
	for _, c := range configs {
		switch c.(type) {
		case storer.Storer:
			w.loginStorer = c.(storer.Storer)
		default:
			return &WechatWeb{}, fmt.Errorf("unknown config: %#v", c)
		}
	}
	return w, nil
}

func (wxwb *WechatWeb) baseRequest() (baseRequest *datastruct.BaseRequest) {
	return &datastruct.BaseRequest{
		Uin:      wxwb.loginInfo.cookie.Wxuin,
		Sid:      wxwb.loginInfo.cookie.Wxsid,
		Skey:     wxwb.loginInfo.sKey,
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

// GetContactByAlias 根据Alias获取联系人
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

// GetContactByRemarkName 根据备注获取用户名
func (wxwb *WechatWeb) GetContactByRemarkName(remarkName string) (contact datastruct.Contact, err error) {
	found := false
	for _, v := range wxwb.contactList {
		if v.RemarkName == remarkName {
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

// GetRunInfo 获取运行计数器信息
func (wxwb *WechatWeb) GetRunInfo() (runinfo WechatRunInfo) {
	return wxwb.runInfo
}

// refreshCookie 根据response更新cookie
func (wxwb *WechatWeb) refreshCookie(cookies []*http.Cookie) {
	for _, c := range cookies {
		switch c.Name {
		case "wxuin":
			wxwb.loginInfo.cookie.Wxuin = c.Value
		case "wxsid":
			wxwb.loginInfo.cookie.Wxsid = c.Value
		case "webwxuvid":
			wxwb.loginInfo.cookie.Uvid = c.Value
		case "webwx_data_ticket":
			wxwb.loginInfo.cookie.DataTicket = c.Value
		case "webwx_auth_ticket":
			wxwb.loginInfo.cookie.AuthTicket = c.Value
		}
	}
	// 如有必要，记录login信息到storer
	wxwb.writeLoginInfo()
}

// 重置登录信息
func (wxwb *WechatWeb) resetLoginInfo() error {
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in resetLoginInfo: ", r)
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
	wxwb.client.Jar = jar
	wxwb.loginInfo = wechatLoginInfo{}
	return nil
}

// 往storer中写入信息
func (wxwb *WechatWeb) writeLoginInfo() error {
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in writeLoginInfo: ", r)
				wxwb.runInfo.PanicCount++
			}
		}()
	}()
	cookieMap := make(map[string][]*http.Cookie)
	for _, host := range syncHosts {
		u, _ := url.Parse("https://" + host)
		cookieMap[host] = wxwb.client.Jar.Cookies(u)
	}
	if wxwb.loginStorer != nil {
		storeInfo := storeLoginInfo{
			Cookies:    cookieMap,
			Cookie:     wxwb.loginInfo.cookie,
			SyncKey:    wxwb.loginInfo.syncKey,
			SKey:       wxwb.loginInfo.sKey,
			PassTicket: wxwb.loginInfo.PassTicket,
			User:       wxwb.user,
			DeviceID:   wxwb.deviceID,
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
				log.Println("Recovered in readLoginInfo: ", r)
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
			wxwb.client.Jar.SetCookies(u, storeInfo.Cookies[host])
		}
		wxwb.loginInfo = wechatLoginInfo{
			cookie:     storeInfo.Cookie,
			syncKey:    storeInfo.SyncKey,
			sKey:       storeInfo.SKey,
			PassTicket: storeInfo.PassTicket,
		}
		wxwb.runInfo = storeInfo.RunInfo
		wxwb.user = storeInfo.User
		wxwb.deviceID = storeInfo.DeviceID
	}
	return nil
}

// 统一请求
func (wxwb *WechatWeb) request(req *http.Request) (resp *http.Response, err error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	resp, err = wxwb.client.Do(req)
	if err == nil {
		wxwb.refreshCookie(resp.Cookies())
	}
	return
}
