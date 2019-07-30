package wwdk

import (
	"fmt"

	"github.com/ikuiki/storer"
	"github.com/kataras/golog"
	"github.com/pkg/errors"

	// "crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"

	"time"

	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
)

// wechatCookie 微信登陆后的cookie凭据，登陆后的消息同步等操作需要此凭据
// TODO: delete
type wechatCookie struct {
	Wxsid      string
	Wxuin      string // 应该是用户的唯一识别号，同一个用户每次登陆此字段都相同
	Uvid       string
	DataTicket string
	AuthTicket string
	// loadTime   string // 登陆时间(10位时间戳字符串)
}

// TODO: delete
type wechatLoginInfo struct {
	cookie     wechatCookie
	syncKey    *datastruct.SyncKey
	sKey       string
	PassTicket string
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

// userInfo 微信用户信息，包含用户、联系人列表等信息
type userInfo struct {
	user        *datastruct.User
	contactList map[string]datastruct.Contact
}

// apiRuntime 微信web客户端运行时信息
type apiRuntime struct {
	userAgent string
	apiDomain string // 当前的apiDomain，从用户扫码登陆后返回的RedirectURL中解析
	client    *http.Client
	deviceID  string // 由客户端生成，为e+15位随机数
}

// WechatWeb 微信网页版客户端实例
type WechatWeb struct {
	userInfo    userInfo        // 用户信息
	apiRuntime  apiRuntime      // wechat客户端运行时信息
	loginInfo   wechatLoginInfo // 登陆信息
	runInfo     WechatRunInfo   // 运行统计信息
	loginStorer storer.Storer   // 存储器，如果有赋值，则用于记录登录信息
	logger      *golog.Logger   // 日志输出器
	mediaStorer MediaStorer     // 媒体存储器，用于处理微信的媒体信息（如用户头像、发送的图片、视频、音频等
}

// NewWechatWeb 生成微信网页版客户端实例
func NewWechatWeb(configs ...interface{}) (wxweb *WechatWeb, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &WechatWeb{}, err
	}
	w := &WechatWeb{
		userInfo: userInfo{
			contactList: make(map[string]datastruct.Contact),
		},
		apiRuntime: apiRuntime{
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
					// 	InsecureSkipVerify: true,
					// },
					// Proxy: func(_ *http.Request) (*url.URL, error) {
					// 	return url.Parse("http://127.0.0.1:8888") //根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
					// },
				},
				Jar:     jar,
				Timeout: 1 * time.Minute,
			},
		},
		runInfo: WechatRunInfo{
			StartAt: time.Now(),
		},
		logger:      golog.Default.Clone(),
		mediaStorer: NewLocalMediaStorer("./"),
	}
	for _, c := range configs {
		switch c.(type) {
		case storer.Storer:
			w.loginStorer = c.(storer.Storer)
		case *golog.Logger:
			w.logger = c.(*golog.Logger)
		case MediaStorer:
			w.mediaStorer = c.(MediaStorer)
		default:
			return &WechatWeb{}, fmt.Errorf("unknown config: %#v", c)
		}
	}
	return w, nil
}

// TODO: delete
func (wxwb *WechatWeb) baseRequest() (baseRequest *datastruct.BaseRequest) {
	return &datastruct.BaseRequest{
		Uin:      wxwb.loginInfo.cookie.Wxuin,
		Sid:      wxwb.loginInfo.cookie.Wxsid,
		Skey:     wxwb.loginInfo.sKey,
		DeviceID: wxwb.apiRuntime.deviceID,
	}
}

// refreshCookie 根据response更新cookie
// TODO: delete
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

// 统一请求
// TODO: delete
func (wxwb *WechatWeb) request(req *http.Request) (resp *http.Response, err error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	resp, err = wxwb.apiRuntime.client.Do(req)
	if err == nil {
		wxwb.refreshCookie(resp.Cookies())
	}
	return
}
