package main

import (
	"fmt"
	"github.com/ikuiki/wwdk"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/datastruct/appmsg"
	"github.com/ikuiki/wwdk/storer"
	"github.com/mdp/qrterminal"
	"log"
	"os"
	"time"
)

func main() {
	// 初始化一个文件loginStorer
	storer := storer.MustNewFileStorer("loginInfo.txt")
	wx, err := wwdk.NewWechatWeb(storer)
	if err != nil {
		panic("Get new wechatweb client error: " + err.Error())
	}
	t := testServ{}
	err = wx.RegisterHook(wwdk.TextMessageHook(t.ProcessTextMessage))
	if err != nil {
		panic("RegisterHook TextMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wwdk.ImageMessageHook(t.ProcessImageMessage))
	if err != nil {
		panic("RegisterHook ImageMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wwdk.EmotionMessageHook(ProcessEmojiMessage))
	if err != nil {
		panic("RegisterHook EmotionMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wwdk.RevokeMessageHook(ProcessRevokeMessage))
	if err != nil {
		panic("RegisterHook RevokeMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wwdk.VideoMessageHook(ProcessVideoMessage))
	if err != nil {
		panic("RegisterHook VideoMessageHook: " + err.Error())
	}
	err = wx.RegisterHook(wwdk.VoiceMessageHook(ProcessVoiceMessage))
	if err != nil {
		panic("RegisterHook VoiceMessageHook: " + err.Error())
	}
	loginChan := make(chan wwdk.LoginChannelItem)
	wx.Login(loginChan)
	for item := range loginChan {
		switch item.Code {
		case wwdk.LoginStatusWaitForScan:
			qrterminal.Generate(item.Msg, qrterminal.L, os.Stdout)
		case wwdk.LoginStatusScanedWaitForLogin:
			fmt.Println("scaned")
		case wwdk.LoginStatusScanedFinish:
			fmt.Println("accepted")
		case wwdk.LoginStatusGotCookie:
			fmt.Println("got cookie")
		case wwdk.LoginStatusInitFinish:
			fmt.Println("init finish")
		case wwdk.LoginStatusGotContact:
			fmt.Println("got contact")
		case wwdk.LoginStatusGotBatchContact:
			fmt.Println("got batch contact")
		case wwdk.LoginStatusErrorOccurred:
			panic(fmt.Sprintf("WxWeb Login error: %+v", item.Err))
		default:
			fmt.Printf("unknown code: %+v", item)
		}
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

func (serv *testServ) ProcessTextMessage(ctx *wwdk.Context, msg datastruct.Message) {
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
func (serv *testServ) ProcessImageMessage(ctx *wwdk.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a image msg from %s\n", from.NickName)
}

// ProcessEmojiMessage set Emoji message Handle
func ProcessEmojiMessage(ctx *wwdk.Context, msg datastruct.Message, emojiContent appmsg.EmotionMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a emotion from %s url: %s\n", from.NickName, emojiContent.Emoji.CdnURL)
}

// ProcessRevokeMessage set revoke message handle
func ProcessRevokeMessage(ctx *wwdk.Context, msg datastruct.Message, revokeContent appmsg.RevokeMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("With %s chat: %s", from.NickName, revokeContent.RevokeMsg.ReplaceMsg)
}

// ProcessVideoMessage set video message handle
func ProcessVideoMessage(ctx *wwdk.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived video from %s", from.NickName)
}

// ProcessVoiceMessage set voice message handle
func ProcessVoiceMessage(ctx *wwdk.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived voice from %s", from.NickName)
}
