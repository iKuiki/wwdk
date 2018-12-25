package main

import (
	"fmt"
	"github.com/yinhui87/wechat-web"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/datastruct/appmsg"
	"log"
	"time"
)

func main() {
	wx, err := wxweb.NewWechatWeb()
	if err != nil {
		panic("Get new wechatweb client error: " + err.Error())
	}
	t := testServ{}
	err = wx.RegisterHook(wxweb.TextMessageHook(t.ProcessTextMessage))
	if err != nil {
		panic("RegisterHook TextMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wxweb.ImageMessageHook(t.ProcessImageMessage))
	if err != nil {
		panic("RegisterHook ImageMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wxweb.EmotionMessageHook(ProcessEmojiMessage))
	if err != nil {
		panic("RegisterHook EmotionMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wxweb.RevokeMessageHook(ProcessRevokeMessage))
	if err != nil {
		panic("RegisterHook RevokeMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wxweb.VideoMessageHook(ProcessVideoMessage))
	if err != nil {
		panic("RegisterHook VideoMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wxweb.VoiceMessageHook(ProcessVoiceMessage))
	if err != nil {
		panic("RegisterHook VoiceMessageHook: " + err.Error())
	}
	err = wx.Login()
	if err != nil {
		panic("WxWeb Login error: " + err.Error())
	}
	contacts := wx.GetContactList()
	for _, v := range contacts {
		if v.IsStar() {
			fmt.Println("Star Friend: " + v.NickName)
		}
		if v.IsTop() {
			fmt.Println("Top Friend: " + v.NickName)
		}
	}
	wx.StartServe()
}

type testServ struct {
}

func (serv *testServ) ProcessTextMessage(ctx *wxweb.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a text msg from %s: %s", from.NickName, msg.Content)
	// reply the same message
	smResp, err := ctx.App.SendTextMessage(msg.FromUserName, msg.Content)
	if err != nil {
		log.Println("sendTextMessage error: " + err.Error())
		return
	}
	log.Println("messageSent, msgId: " + smResp.MsgID + ", Local ID: " + smResp.LocalID)
	// Set message to readed at phone
	err = ctx.App.StatusNotify(msg.ToUserName, msg.FromUserName, 1)
	if err != nil {
		log.Println("StatusNotify error: " + err.Error())
		return
	}
	go func() {
		time.Sleep(10 * time.Second)
		ctx.App.SendRevokeMessage(smResp.MsgID, smResp.LocalID, msg.FromUserName)
	}()
}

// ProcessImageMessage set image message handle
func (serv *testServ) ProcessImageMessage(ctx *wxweb.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a image msg from %s\n", from.NickName)
}

// ProcessEmojiMessage set Emoji message Handle
func ProcessEmojiMessage(ctx *wxweb.Context, msg datastruct.Message, emojiContent appmsg.EmotionMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a emotion from %s url: %s\n", from.NickName, emojiContent.Emoji.CdnURL)
}

// ProcessRevokeMessage set revoke message handle
func ProcessRevokeMessage(ctx *wxweb.Context, msg datastruct.Message, revokeContent appmsg.RevokeMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("With %s chat: %s", from.NickName, revokeContent.RevokeMsg.ReplaceMsg)
}

// ProcessVideoMessage set video message handle
func ProcessVideoMessage(ctx *wxweb.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived video from %s", from.NickName)
}

// ProcessVoiceMessage set voice message handle
func ProcessVoiceMessage(ctx *wxweb.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived voice from %s", from.NickName)
}
