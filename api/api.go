package api

import (
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// WechatwebAPI 微信网页版api
type WechatwebAPI struct {
	userAgent string
	apiDomain string // 当前的apiDomain，从用户扫码登陆后返回的RedirectURL中解析
	client    *http.Client
	deviceID  string // 由客户端生成，为e+15位随机数
	loginInfo LoginInfo
}

// NewWechatwebAPI 创建WechatwebAPI
func NewWechatwebAPI() (wechatAPI *WechatwebAPI, err error) {
	// 创建cookie jar用于持久化cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &WechatwebAPI{}, err
	}
	return &WechatwebAPI{
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		deviceID:  "e" + tool.GetRandomStringFromNum(15),
		apiDomain: "wx.qq.com", // 默认域名
		client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 1 * time.Minute,
				}).Dial,
				TLSHandshakeTimeout: 1 * time.Minute,
				// TLSClientConfig: &tls.Config{
				// 	InsecureSkipVerify: true, // 跳过https证书验证
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

// LoginInfo 登陆信息，登陆后可以获取到
type LoginInfo struct {
	SyncKey    *datastruct.SyncKey
	SKey       string // 更新时机：webwxinit\webwxsync
	PassTicket string
	// 以下更新时机都是刷新cookie时
	Wxsid      string
	Wxuin      string // 应该是用户的唯一识别号，同一个用户每次登陆此字段都相同
	Uvid       string
	AuthTicket string
	DataTicket string
	// loadTime   string // 登陆时间(10位时间戳字符串)
}

// 统一请求
func (api *WechatwebAPI) request(req *http.Request) (resp *http.Response, err error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	resp, err = api.client.Do(req)
	if err == nil {
		api.refreshCookie(resp.Cookies())
	}
	return
}

// refreshCookie 根据response更新cookie
func (api *WechatwebAPI) refreshCookie(cookies []*http.Cookie) {
	for _, c := range cookies {
		switch c.Name {
		case "wxuin":
			api.loginInfo.Wxuin = c.Value
		case "wxsid":
			api.loginInfo.Wxsid = c.Value
		case "webwxuvid":
			api.loginInfo.Uvid = c.Value
		case "webwx_data_ticket":
			api.loginInfo.DataTicket = c.Value
		case "webwx_auth_ticket":
			api.loginInfo.AuthTicket = c.Value
		}
	}
	// TODO: 在wwdk包做退出前保存
}

func (api *WechatwebAPI) baseRequest() (baseRequest *datastruct.BaseRequest) {
	return &datastruct.BaseRequest{
		Uin:      api.loginInfo.Wxuin,
		Sid:      api.loginInfo.Wxsid,
		Skey:     api.loginInfo.SKey,
		DeviceID: api.deviceID,
	}
}
