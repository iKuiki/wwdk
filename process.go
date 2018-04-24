package wxweb

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/datastruct/appmsg"
	"html"
	"strings"
)

func (wxwb *WechatWeb) messageProcesser(msg *datastruct.Message) (err error) {
	defer func() {
		// 防止外部方法导致的崩溃
		if err := recover(); err != nil {
			fmt.Println("messageProcesser panic: ", err)
			fmt.Println("message data: ", msg)
		}
	}()
	context := Context{App: wxwb, hasStop: false}
	switch msg.MsgType {
	case datastruct.TextMsg:
		for _, v := range wxwb.messageHook[datastruct.TextMsg] {
			if f, ok := v.(TextMessageHook); ok {
				f(&context, *msg)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.ImageMsg:
		msg.Content = strings.Replace(html.UnescapeString(msg.Content), "<br/>", "", -1)
		var imgContent appmsg.ImageMsgContent
		err = xml.Unmarshal([]byte(msg.Content), &imgContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.ImageMsg] {
			if f, ok := v.(ImageMessageHook); ok {
				f(&context, *msg, imgContent)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.AnimationEmotionsMsg:
		msg.Content = html.UnescapeString(msg.Content)
		var emojiContent appmsg.EmotionMsgContent
		err := xml.Unmarshal([]byte(msg.Content), &emojiContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.AnimationEmotionsMsg] {
			if f, ok := v.(EmotionMessageHook); ok {
				f(&context, *msg, emojiContent)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.RevokeMsg:
		msg.Content = html.UnescapeString(msg.Content)
		var revokeContent appmsg.RevokeMsgContent
		err := xml.Unmarshal([]byte(msg.Content), &revokeContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.RevokeMsg] {
			if f, ok := v.(RevokeMessageHook); ok {
				f(&context, *msg, revokeContent)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.LittleVideoMsg:
		msg.Content = strings.Replace(html.UnescapeString(msg.Content), "<br/>", "", -1)
		var videoContent appmsg.VideoMsgContent
		err := xml.Unmarshal([]byte(msg.Content), &videoContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.LittleVideoMsg] {
			if f, ok := v.(VideoMessageHook); ok {
				f(&context, *msg, videoContent)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.VoiceMsg:
		msg.Content = html.UnescapeString(msg.Content)
		var voiceContent appmsg.VoiceMsgContent
		err := xml.Unmarshal([]byte(msg.Content), &voiceContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.VoiceMsg] {
			if f, ok := v.(VoiceMessageHook); ok {
				f(&context, *msg, voiceContent)
			}
			if context.hasStop {
				break
			}
		}
	default:
		return fmt.Errorf("Unknown MsgType %v: %#v", msg.MsgType, msg)
	}
	return nil
}

func (wxwb *WechatWeb) contactProcesser(contact *datastruct.Contact) (err error) {
	defer func() {
		// 防止外部方法导致的崩溃
		if err := recover(); err != nil {
			fmt.Println("contactProcesser panic: ", err)
			fmt.Println("contact data: ", contact)
		}
	}()
	context := Context{App: wxwb, hasStop: false}
	for _, v := range wxwb.modContactHook {
		if f, ok := v.(ModContactHook); ok {
			f(&context, *contact)
		}
		if context.hasStop {
			break
		}
	}
	return
}
