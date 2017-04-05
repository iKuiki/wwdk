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
	wx := wxweb.NewWechatWeb()
	t := testServ{}
	err := wx.RegisterMessageHook(wxweb.TextMessageHook(t.ProcessTextMessage))
	if err != nil {
		panic("RegisterMessageHook TextMessageHook: " + err.Error())
	}
	err = wx.RegisterMessageHook(wxweb.ImageMessageHook(t.ProcessImageMessage))
	if err != nil {
		panic("RegisterMessageHook ImageMessageHook: " + err.Error())
	}
	err = wx.RegisterMessageHook(wxweb.EmotionMessageHook(ProcessEmojiMessage))
	if err != nil {
		panic("RegisterMessageHook EmotionMessageHook: " + err.Error())
	}
	err = wx.RegisterMessageHook(wxweb.RevokeMessageHook(ProcessRevokeMessage))
	if err != nil {
		panic("RegisterMessageHook RevokeMessageHook: " + err.Error())
	}
	err = wx.RegisterMessageHook(wxweb.VideoMessageHook(ProcessVideoMessage))
	if err != nil {
		panic("RegisterMessageHook VideoMessageHook: " + err.Error())
	}
	err = wx.RegisterMessageHook(wxweb.VoiceMessageHook(ProcessVoiceMessage))
	if err != nil {
		panic("RegisterMessageHook VoiceMessageHook: " + err.Error())
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

func (this *testServ) ProcessTextMessage(ctx *wxweb.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
	}
	log.Printf("Recived a text msg from %s: %s", from.NickName, msg.Content)
	// reply the same message
	smResp, err := ctx.App.SendTextMessage(msg.FromUserName, msg.Content)
	if err != nil {
		log.Println("sendTextMessage error: " + err.Error())
	}
	log.Println("messageSent, msgId: " + smResp.MsgID + ", Local ID: " + smResp.LocalID)
	// Set message to readed at phone
	err = ctx.App.StatusNotify(msg.ToUserName, msg.FromUserName)
	if err != nil {
		log.Println("StatusNotify error: " + err.Error())
	}
	go func() {
		time.Sleep(10 * time.Second)
		ctx.App.SendRevokeMessage(smResp.MsgID, smResp.LocalID, msg.FromUserName)
	}()
}

func (this *testServ) ProcessImageMessage(ctx *wxweb.Context, msg datastruct.Message, imgContent appmsg.ImageMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
	}
	log.Printf("Recived a image msg from %s\n", from.NickName)
	fmt.Println("aeskey: ", imgContent.Img.AesKey)
}

func ProcessEmojiMessage(ctx *wxweb.Context, msg datastruct.Message, emojiContent appmsg.EmotionMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
	}
	log.Printf("Recived a emotion from %s url: %s\n", from.NickName, emojiContent.Emoji.CdnUrl)
}

func ProcessRevokeMessage(ctx *wxweb.Context, msg datastruct.Message, revokeContent appmsg.RevokeMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
	}
	log.Printf("With %s chat: %s", from.NickName, revokeContent.RevokeMsg.ReplaceMsg)
}

func ProcessVideoMessage(ctx *wxweb.Context, msg datastruct.Message, videoContent appmsg.VideoMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
	}
	log.Printf("Recived video from %s: %s", from.NickName, videoContent.VideoMsg.AesKey)
}

func ProcessVoiceMessage(ctx *wxweb.Context, msg datastruct.Message, voiceContent appmsg.VoiceMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
	}
	log.Printf("Recived voice from %s, length %s ms", from.NickName, voiceContent.VoiceMsg.VoiceLength)
}
