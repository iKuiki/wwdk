package wxweb

import (
	"errors"
	"github.com/ikuiki/wechat-web/datastruct"
	"github.com/ikuiki/wechat-web/datastruct/appmsg"
)

// TextMessageHook 文字消息处理器接口
type TextMessageHook func(*Context, datastruct.Message)

// ImageMessageHook 图片消息处理器接口
type ImageMessageHook func(*Context, datastruct.Message)

// EmotionMessageHook 表情消息处理器接口
type EmotionMessageHook func(*Context, datastruct.Message, appmsg.EmotionMsgContent)

// RevokeMessageHook 撤回消息处理器接口
type RevokeMessageHook func(*Context, datastruct.Message, appmsg.RevokeMsgContent)

// VideoMessageHook 视频消息处理器接口
type VideoMessageHook func(*Context, datastruct.Message)

// VoiceMessageHook 语音消息处理器接口
type VoiceMessageHook func(*Context, datastruct.Message)

// ModContactHook 联系人变动处理接口
type ModContactHook func(context *Context, oldContact, newContact *datastruct.Contact)

// RegisterHook 注册处理器，需要传入处理器接口类型，会自动识别
func (wxwb *WechatWeb) RegisterHook(hook interface{}) error {
	switch hook.(type) {
	case TextMessageHook:
		wxwb.messageHook[datastruct.TextMsg] = append(wxwb.messageHook[datastruct.TextMsg], hook)
	case ImageMessageHook:
		wxwb.messageHook[datastruct.ImageMsg] = append(wxwb.messageHook[datastruct.ImageMsg], hook)
	case EmotionMessageHook:
		wxwb.messageHook[datastruct.AnimationEmotionsMsg] = append(wxwb.messageHook[datastruct.AnimationEmotionsMsg], hook)
	case RevokeMessageHook:
		wxwb.messageHook[datastruct.RevokeMsg] = append(wxwb.messageHook[datastruct.RevokeMsg], hook)
	case VideoMessageHook:
		wxwb.messageHook[datastruct.LittleVideoMsg] = append(wxwb.messageHook[datastruct.LittleVideoMsg], hook)
	case VoiceMessageHook:
		wxwb.messageHook[datastruct.VoiceMsg] = append(wxwb.messageHook[datastruct.VoiceMsg], hook)
	case ModContactHook:
		wxwb.modContactHook = append(wxwb.modContactHook, hook)
	default:
		return errors.New("Unknown hook function")
	}
	return nil
}
