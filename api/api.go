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
type WechatwebAPI interface {
	// 登陆部分

	// JsLogin 获取uuid
	JsLogin() (uuid string, body []byte, err error)
	// Login 等待用户扫码登陆
	Login(tip, uuid string) (code, userAvatar, redirectURL string, body []byte, err error)
	// WebwxNewLoginPage 获取登陆凭据
	WebwxNewLoginPage(redirectURL string) (body []byte, err error)
	// WebwxInit 初始化微信
	WebwxInit() (user *datastruct.User, contactList []datastruct.Contact, body []byte, err error)

	// 联系人部分

	// GetContact 获取联系人
	GetContact() (contactList []datastruct.Contact, body []byte, err error)
	// BatchGetContact 获取群聊的成员
	BatchGetContact(contactItemList []datastruct.BatchGetContactRequestListItem) (contactList []datastruct.Contact, body []byte, err error)
	// ModifyUserRemakName 修改联系人备注
	ModifyUserRemakName(userName, remarkName string) (body []byte, err error)

	// 聊天室部分

	// ModifyChatRoomTopic 修改聊天室标题
	ModifyChatRoomTopic(userName, newTopic string) (body []byte, err error)

	// 同步部分

	// SyncCheck 检查同步
	SyncCheck() (retCode, selector string, body []byte, err error)
	// WebwxSync 同步消息
	WebwxSync() (modContacts []datastruct.Contact,
		delContacts []datastruct.WebwxSyncRespondDelContactListItem,
		addMessages []datastruct.Message,
		body []byte, err error)

	// 发送部分

	// StatusNotify 消息已读通知
	StatusNotify(fromUserName, toUserName string, code int64) (body []byte, err error)
	// SendTextMessage 发送消息
	SendTextMessage(fromUserName, toUserName, content string) (MsgID, LocalID string, body []byte, err error)
	// SendRevokeMessage 撤回消息
	SendRevokeMessage(toUserName, svrMsgID, clientMsgID string) (body []byte, err error)

	// 接收部分

	// SaveMessageImage 下载图片消息
	SaveMessageImage(msgID string) (imgData []byte, err error)
	// SaveMessageVoice 下载音频消息
	SaveMessageVoice(msgID string) (voiceData []byte, err error)
	// SaveMessageVideo 下载视频消息
	SaveMessageVideo(msgID string) (videoData []byte, err error)
	// SaveContactImg 保存联系人头像
	SaveContactImg(headImgURL string) (imgData []byte, err error)
	// SaveMemberImg 保存群成员的头像
	SaveMemberImg(userName, chatroomID string) (imgData []byte, err error)
}

// wechatwebAPI 微信网页版api
type wechatwebAPI struct {
	userAgent string
	apiDomain string // 当前的apiDomain，从用户扫码登陆后返回的RedirectURL中解析
	client    *http.Client
	deviceID  string // 由客户端生成，为e+15位随机数
	loginInfo LoginInfo
}

// NewWechatwebAPI 创建WechatwebAPI
func NewWechatwebAPI() (wechatAPI WechatwebAPI, err error) {
	// 创建cookie jar用于持久化cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &wechatwebAPI{}, err
	}
	return &wechatwebAPI{
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
func (api *wechatwebAPI) request(req *http.Request) (resp *http.Response, err error) {
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
func (api *wechatwebAPI) refreshCookie(cookies []*http.Cookie) {
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

func (api *wechatwebAPI) baseRequest() (baseRequest *datastruct.BaseRequest) {
	return &datastruct.BaseRequest{
		Uin:      api.loginInfo.Wxuin,
		Sid:      api.loginInfo.Wxsid,
		Skey:     api.loginInfo.SKey,
		DeviceID: api.deviceID,
	}
}
