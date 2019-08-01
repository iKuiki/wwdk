package wwdk

import (
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"html"
	"runtime/debug"
	"strings"

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
	// 收到信息分3种情况
	// 收到文字信息：无需处理
	// 收到撤回信息：更新的是撤回计数器
	// 收到其他消息：解码Content后再放入channel
	switch msg.MsgType {
	case datastruct.TextMsg:
		wxwb.runInfo.MessageRecivedCount++
		syncChannel <- SyncChannelItem{
			Code:    SyncStatusNewMessage,
			Message: msg,
		}
	case datastruct.ImageMsg:
		fallthrough
	case datastruct.AnimationEmotionsMsg:
		fallthrough
	case datastruct.LittleVideoMsg:
		fallthrough
	case datastruct.VoiceMsg:
		wxwb.runInfo.MessageRecivedCount++
		msg.Content = strings.Replace(html.UnescapeString(msg.Content), "<br/>", "", -1)
		syncChannel <- SyncChannelItem{
			Code:    SyncStatusNewMessage,
			Message: msg,
		}
	case datastruct.RevokeMsg:
		wxwb.runInfo.MessageRevokeRecivedCount++
		msg.Content = strings.Replace(html.UnescapeString(msg.Content), "<br/>", "", -1)
		syncChannel <- SyncChannelItem{
			Code:    SyncStatusNewMessage,
			Message: msg,
		}
	default:
		return errors.Errorf("Unknown MsgType %v: %#v", msg.MsgType, msg)
	}
	return nil
}
