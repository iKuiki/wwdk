package main

import (
	"encoding/xml"
	"fmt"
	"github.com/ikuiki/wwdk"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/datastruct/appmsg"
	"github.com/ikuiki/wwdk/storer"
	"github.com/mdp/qrterminal"
	"github.com/pkg/errors"
	"log"
	"os"
)

func main() {
	// 初始化一个文件loginStorer
	storer := storer.MustNewFileStorer("loginInfo.txt")
	wx, err := wwdk.NewWechatWeb(storer)
	if err != nil {
		panic("Get new wechatweb client error: " + err.Error())
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
	syncChannel := make(chan wwdk.SyncChannelItem)
	wx.StartServe(syncChannel)
	for item := range syncChannel {
		switch item.Code {
		case wwdk.SyncStatusModifyContact:
			fmt.Println("Modify contact: ", item.Contact)
		case wwdk.SyncStatusNewMessage:
			msg := item.Message
			switch msg.MsgType {
			case datastruct.TextMsg:
				processTextMessage(wx, msg)
			case datastruct.ImageMsg:
				processImageMessage(wx, msg)
			case datastruct.AnimationEmotionsMsg:
				var emojiContent appmsg.EmotionMsgContent
				err := xml.Unmarshal([]byte(msg.GetContent()), &emojiContent)
				if err != nil {
					panic(errors.New("Unmarshal message content to struct: " + err.Error()))
				}
				processEmojiMessage(wx, msg, emojiContent)
			case datastruct.RevokeMsg:
				var revokeContent appmsg.RevokeMsgContent
				err := xml.Unmarshal([]byte(msg.GetContent()), &revokeContent)
				if err != nil {
					panic(errors.New("Unmarshal message content to struct: " + err.Error()))
				}
				processRevokeMessage(wx, msg, revokeContent)
			case datastruct.LittleVideoMsg:
				processVideoMessage(wx, msg)
			case datastruct.VoiceMsg:
				processVoiceMessage(wx, msg)
			}
		case wwdk.SyncStatusErrorOccurred:
			fmt.Printf("error occurred at sync: %+v\n", item.Err)
		case wwdk.SyncStatusPanic:
			fmt.Printf("sync panic: %+v\n", err)
			break
		}
	}
}

func processTextMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a text msg from %s: %s", from.NickName, msg.Content)
	// // reply the same message
	// smResp, err := app.SendTextMessage(msg.FromUserName, msg.Content)
	// if err != nil {
	// 	log.Println("sendTextMessage error: " + err.Error())
	// 	return
	// }
	// log.Println("messageSent, msgId: " + smResp.MsgID + ", Local ID: " + smResp.LocalID)
	// // Set message to readed at phone
	// err = app.StatusNotify(msg.ToUserName, msg.FromUserName, 1)
	// if err != nil {
	// 	log.Println("StatusNotify error: " + err.Error())
	// 	return
	// }
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	app.SendRevokeMessage(smResp.MsgID, smResp.LocalID, msg.FromUserName)
	// }()
}

// ProcessImageMessage set image message handle
func processImageMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a image msg from %s\n", from.NickName)
}

// ProcessEmojiMessage set Emoji message Handle
func processEmojiMessage(app *wwdk.WechatWeb, msg *datastruct.Message, emojiContent appmsg.EmotionMsgContent) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a emotion from %s url: %s\n", from.NickName, emojiContent.Emoji.CdnURL)
}

// ProcessRevokeMessage set revoke message handle
func processRevokeMessage(app *wwdk.WechatWeb, msg *datastruct.Message, revokeContent appmsg.RevokeMsgContent) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("With %s chat: %s", from.NickName, revokeContent.RevokeMsg.ReplaceMsg)
}

// ProcessVideoMessage set video message handle
func processVideoMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived video from %s", from.NickName)
}

// ProcessVoiceMessage set voice message handle
func processVoiceMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived voice from %s", from.NickName)
}
