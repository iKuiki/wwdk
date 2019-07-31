package wwdk

import (
	"fmt"
	"github.com/ikuiki/wwdk/api"

	"github.com/ikuiki/storer"
	"github.com/kataras/golog"

	"time"

	"github.com/ikuiki/wwdk/datastruct"
)

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

// WechatWeb 微信网页版客户端实例
type WechatWeb struct {
	userInfo    userInfo         // 用户信息
	api         api.WechatwebAPI // 微信网页版的api实现
	runInfo     WechatRunInfo    // 运行统计信息
	loginStorer storer.Storer    // 存储器，如果有赋值，则用于记录登录信息
	logger      *golog.Logger    // 日志输出器
	mediaStorer MediaStorer      // 媒体存储器，用于处理微信的媒体信息（如用户头像、发送的图片、视频、音频等
}

// NewWechatWeb 生成微信网页版客户端实例
func NewWechatWeb(configs ...interface{}) (wxweb *WechatWeb, err error) {
	a, err := api.NewWechatwebAPI()
	if err != nil {
		return nil, err
	}
	w := &WechatWeb{
		userInfo: userInfo{
			contactList: make(map[string]datastruct.Contact),
		},
		api: a,
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
