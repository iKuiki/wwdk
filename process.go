package wwdk

import (
	"runtime/debug"

	"github.com/getsentry/sentry-go"

	"github.com/ikuiki/wwdk/datastruct"
)

func (wxwb *WechatWeb) messageProcesser(msg *datastruct.Message, syncChannel chan<- SyncChannelItem) (err error) {
	defer func() {
		// 防止外部方法导致的崩溃
		if e := recover(); e != nil {
			var eErr error
			if err, ok := e.(error); ok {
				eErr = err
			}
			wxwb.captureException(eErr, "MessageProcesser panic", sentry.LevelError,
				extraData{"panicItem", e},
				extraData{"msg", msg},
			)
			wxwb.logger.Errorf("messageProcesser panic: %v\n", e)
			wxwb.logger.Errorf("message data: %v\n", msg)
			wxwb.logger.Errorf("Stack: %s\n", string(debug.Stack()))
		}
	}()
	// 收到消息后，更新消息计数器再传入channel中
	// PS: 现在不再对消息Content做Unescape操作，而是由Message直接提供GetContentUnescape方法
	switch msg.MsgType {
	case datastruct.TextMsg:
		fallthrough
	case datastruct.ImageMsg:
		fallthrough
	case datastruct.AnimationEmotionsMsg:
		fallthrough
	case datastruct.LittleVideoMsg:
		fallthrough
	case datastruct.VoiceMsg:
		wxwb.runInfo.MessageRecivedCount++
	case datastruct.RevokeMsg:
		wxwb.runInfo.MessageRevokeRecivedCount++
	default:
		wxwb.captureException(nil, "Unknown MsgType", sentry.LevelWarning, extraData{"msgType", msg.MsgType}, extraData{"msg", msg})
		wxwb.logger.Infof("Unknown MsgType %v: %#v", msg.MsgType, msg)
	}
	syncChannel <- SyncChannelItem{
		Code:    SyncStatusNewMessage,
		Message: msg,
	}
	return nil
}
