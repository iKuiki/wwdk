package wwdk

import (
	"reflect"

	"github.com/getsentry/sentry-go"
	"github.com/ikuiki/wwdk/api"
	"github.com/ikuiki/wwdk/contactmgr"
	"github.com/ikuiki/wwdk/runinfo"
	"github.com/pkg/errors"

	"github.com/ikuiki/storer"
	"github.com/kataras/golog"

	"time"

	"github.com/ikuiki/wwdk/datastruct"
)

// WechatWeb 微信网页版客户端实例
type WechatWeb struct {
	api            api.WechatwebAPI       // 微信网页版的api实现
	runInfo        runinfo.WechatRunInfo  // 运行统计信息
	loginStorer    storer.Storer          // 存储器，如果有赋值，则用于记录登录信息
	logger         *golog.Logger          // 日志输出器
	mediaStorer    MediaStorer            // 媒体存储器，用于处理微信的媒体信息（如用户头像、发送的图片、视频、音频等
	syncChannel    chan<- SyncChannelItem // 同步通道，方便除sync方法外发生同步
	user           *datastruct.User       // 当前登录用户，未登录时为nil
	contactManager contactmgr.Manager     // 联系人管理器
	sentryHub      *sentry.Hub            // 用来进行错误追踪的hub，bindClient后生效
}

// NewWechatWeb 生成微信网页版客户端实例
func NewWechatWeb(configs ...interface{}) (wxweb *WechatWeb, err error) {
	a, err := api.NewWechatwebAPI()
	if err != nil {
		return nil, err
	}
	w := &WechatWeb{
		api: a,
		runInfo: runinfo.WechatRunInfo{
			StartAt: time.Now(),
		},
		logger:      golog.Default.Clone(),
		mediaStorer: NewLocalMediaStorer("./"),
		sentryHub:   sentry.NewHub(nil, sentry.NewScope()),
	}
	{
		w.sentryHub.Scope().SetTags(map[string]string{
			"wwdk_version": Version,
		})
		// 设置sentry的tag与extra
		w.sentryHub.Scope().SetExtras(map[string]interface{}{
			"package":      "github.com/iKuiki/wwdk",
			"wwdk_version": Version,
		})
	}
	for _, c := range configs {
		switch c.(type) {
		case storer.Storer:
			w.sentryHub.Scope().SetExtra("loginStore", reflect.TypeOf(c).String())
			w.loginStorer = c.(storer.Storer)
		case *golog.Logger:
			w.logger = c.(*golog.Logger)
		case MediaStorer:
			w.sentryHub.Scope().SetExtra("mediaStorer", reflect.TypeOf(c).String())
			w.mediaStorer = c.(MediaStorer)
		case *sentry.Client:
			w.sentryHub.BindClient(c.(*sentry.Client))
		case contactmgr.Manager:
			w.contactManager = c.(contactmgr.Manager)
		default:
			err = errors.Errorf("unknown config type(%s): %#v", reflect.TypeOf(c).String(), c)
			w.captureException(err, "Unknown wwdk config", sentry.LevelWarning)
			return nil, err
		}
	}
	// 对wxwb做最后处理
	w.logger.NewLine = true // 日志有换行
	if w.contactManager == nil {
		// 如果联系人管理器为空则创建一个
		w.contactManager, err = contactmgr.NewManager()
		if err != nil {
			return
		}
	}
	return w, nil
}
