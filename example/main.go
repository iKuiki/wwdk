package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/ikuiki/storer"
	"github.com/ikuiki/wwdk"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/datastruct/appmsg"
	"github.com/mdp/qrterminal"
	"github.com/pkg/errors"
)

func main() {
	// 初始化一个文件loginStorer
	storer := storer.MustNewFileStorer("loginInfo.txt")
	// 将loginStorer作为配置传入构造函数，可以用来记录登陆状态
	wx, err := wwdk.NewWechatWeb(storer)
	if err != nil {
		panic("Get new wechatweb client error: " + err.Error())
	}
	// 创建登陆用channel用于回传登陆信息
	loginChan := make(chan wwdk.LoginChannelItem)
	wx.Login(loginChan)
	// 根据channel返回信息进行处理
	for item := range loginChan {
		switch item.Code {
		case wwdk.LoginStatusWaitForScan:
			// 返回了登陆二维码链接，输出到屏幕
			qrterminal.Generate(item.Msg, qrterminal.L, os.Stdout)
		case wwdk.LoginStatusScanedWaitForLogin:
			// 用户已扫码
			fmt.Println("scaned")
		case wwdk.LoginStatusScanedFinish:
			// 用户同意登陆
			fmt.Println("accepted")
		case wwdk.LoginStatusGotCookie:
			// 获取到cookie
			fmt.Println("got cookie")
		case wwdk.LoginStatusInitFinish:
			// 初始化完成
			fmt.Println("init finish")
		case wwdk.LoginStatusGotContact:
			// 获取联系人完成
			fmt.Println("got contact")
		case wwdk.LoginStatusBatchGotContact:
			// 获取群成员完成
			fmt.Println("got batch contact")
			break
		case wwdk.LoginStatusErrorOccurred:
			// 登陆失败
			panic(fmt.Sprintf("WxWeb Login error: %+v", item.Err))
		default:
			fmt.Printf("unknown code: %+v", item)
		}
	}
	// 获取联系人
	contacts := wx.GetContactList()
	// 创建联系人username->Contact映射
	contactMap := make(map[string]datastruct.Contact)
	for _, v := range contacts {
		contactMap[v.UserName] = v
		if v.IsStar() {
			fmt.Println("Star Friend: " + v.NickName)
		}
		if v.IsTop() {
			fmt.Println("Top Friend: " + v.NickName)
		}
	}
	// 创建同步channel
	syncChannel := make(chan wwdk.SyncChannelItem)
	// 将channel传入startServe方法，开始同步服务并且将新信息通过syncChannel传回
	wx.StartServe(syncChannel)
	// 处理syncChannel传回信息
	for item := range syncChannel {
		func() {
			// 声明一个匿名方法，并添加recover防止panic
			defer func() {
				if e := recover(); e != nil {
					fmt.Println("panic recovered at sync process func: ", e)
				}
			}()
			// 在子方法内执行逻辑
			switch item.Code {
			case wwdk.SyncStatusModifyContact:
				// 发生联系人变更，处理联系人变更
				if oldContact, ok := contactMap[item.Contact.UserName]; ok {
					// 旧联系人存在
					fmt.Println("Modify contact: ", item.Contact.NickName, ", oldContact is: ", oldContact.NickName)
				} else {
					// 旧联系人不存在，此为新联系人
					fmt.Println("New contact: ", item.Contact.NickName)
				}
				contactMap[item.Contact.UserName] = *item.Contact
			// 收到新信息
			case wwdk.SyncStatusNewMessage:
				// 根据收到的信息类型分别处理
				msg := item.Message
				switch msg.MsgType {
				case datastruct.TextMsg:
					// 处理文字信息
					processTextMessage(wx, msg)
				case datastruct.ImageMsg:
					// 处理图片信息
					processImageMessage(wx, msg)
				case datastruct.AnimationEmotionsMsg:
					// 处理表情信息
					processEmojiMessage(wx, msg)
				case datastruct.RevokeMsg:
					// 反序列化撤回消息附加信息
					var revokeContent appmsg.RevokeMsgContent
					err := xml.Unmarshal([]byte(msg.GetContent()), &revokeContent)
					if err != nil {
						panic(errors.New("Unmarshal message content to struct: " + err.Error()))
					}
					// 处理撤回消息
					processRevokeMessage(wx, msg, revokeContent)
				case datastruct.LittleVideoMsg:
					// 处理视频信息
					processVideoMessage(wx, msg)
				case datastruct.VoiceMsg:
					// 处理声音消息
					processVoiceMessage(wx, msg)
				}
			case wwdk.SyncStatusErrorOccurred:
				// 发生非致命性错误
				fmt.Printf("error occurred at sync: %+v\n", item.Err)
			case wwdk.SyncStatusPanic:
				// 发生致命错误，sync中断
				fmt.Printf("sync panic: %+v\n", err)
				break
			}
		}()
	}
}

func processTextMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Printf("msg[%d]%s getContact[%s] error: %v\n", msg.MsgType, msg.Content, msg.FromUserName, err)
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
		log.Printf("msg[%d] getContact[%s] error: %v\n", msg.MsgType, msg.FromUserName, err)
		return
	}
	filename, err := app.SaveMessageImage(*msg)
	if err != nil {
		log.Printf("Recived a image msg from %s but save fail: %v\n", from.NickName, err)
		return
	}
	log.Printf("Recived a image %s msg from %s\n", filename, from.NickName)
}

// ProcessEmojiMessage set Emoji message Handle
func processEmojiMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Printf("msg[%d] getContact[%s] error: %v\n", msg.MsgType, msg.FromUserName, err)
		return
	}
	log.Printf("Recived a emotion from %s content: %s\n", from.NickName, msg.Content)
}

// ProcessRevokeMessage set revoke message handle
func processRevokeMessage(app *wwdk.WechatWeb, msg *datastruct.Message, revokeContent appmsg.RevokeMsgContent) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Printf("msg[%d] getContact[%s] error: %v\n", msg.MsgType, msg.FromUserName, err)
		return
	}
	log.Printf("With %s chat: %s", from.NickName, revokeContent.RevokeMsg.ReplaceMsg)
}

// ProcessVideoMessage set video message handle
func processVideoMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Printf("msg[%d] getContact[%s] error: %v\n", msg.MsgType, msg.FromUserName, err)
		return
	}
	filename, err := app.SaveMessageVideo(*msg)
	if err != nil {
		log.Printf("Recived a video msg from %s but save fail: %v\n", from.NickName, err)
		return
	}
	log.Printf("Recived a video %s msg from %s\n", filename, from.NickName)
}

// ProcessVoiceMessage set voice message handle
func processVoiceMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Printf("msg[%d] getContact[%s] error: %v\n", msg.MsgType, msg.FromUserName, err)
		return
	}
	filename, err := app.SaveMessageVoice(*msg)
	if err != nil {
		log.Printf("Recived a voice msg from %s but save fail: %v\n", from.NickName, err)
		return
	}
	log.Printf("Recived a voice %s msg from %s\n", filename, from.NickName)
}
