package api_test

import (
	"bytes"
	"fmt"
	"github.com/ikuiki/wwdk/datastruct"
	"io"
	"os"
	"testing"
	"time"
)

// 测试发送一次消息并撤回
func TestSendAndRevoke(t *testing.T) {
	contact, skip := getTestContact("TestSendAndRevoke", anyContact)
	if skip {
		t.SkipNow()
	}
	msgID, localID, _, err := client.SendTextMessage(user.UserName, contact.UserName, "test message will be revoke")
	checkErrorIsNil(err)
	fmt.Println("sleep 2 second")
	time.Sleep(2 * time.Second)
	_, err = client.SendRevokeMessage(contact.UserName, msgID, localID)
	checkErrorIsNil(err)
}

// 测试发送消息已读通知
func TestStatusNotify(t *testing.T) {
	contact, skip := getTestContact("TestStatusNotify", anyContact)
	if skip {
		t.SkipNow()
	}
	fmt.Println("sleep a second")
	time.Sleep(time.Second)
	fmt.Println("now notify readed")
	_, err := client.StatusNotify(user.UserName, contact.UserName, 1)
	checkErrorIsNil(err)
}

func saveFile(filename string, content []byte) {
	f, err := os.Create(filename)
	checkErrorIsNil(err)
	defer f.Close()
	io.Copy(f, bytes.NewBuffer(content))
}

// 测试接收、发送媒体消息
func TestReceiveAndSendMediaMessage(t *testing.T) {
	contact, skip := getTestContact("TestReceiveMediaMessage", anyContact)
	if skip {
		t.SkipNow()
	}
	// 测试接收图片消息
	fmt.Println("please send a image message via " + contact.NickName)
	for {
		msg := <-addMessageChan
		if msg.MsgType == datastruct.ImageMsg {
			body, err := client.SaveMessageImage(msg.MsgID)
			checkErrorIsNil(err)
			filename := msg.MsgID + ".png"
			saveFile(filename, body)
			fmt.Println("image save to " + filename)
			// 发送回去
			mediaID, _, err := client.UploadMedia(user.UserName, contact.UserName, filename, body)
			checkErrorIsNil(err)
			client.SendImageMessage(user.UserName, contact.UserName, mediaID)
			break
		}
	}
	// 测试接收音频消息
	fmt.Println("please send a voice message via " + contact.NickName)
	for {
		msg := <-addMessageChan
		if msg.MsgType == datastruct.VoiceMsg {
			body, err := client.SaveMessageVoice(msg.MsgID)
			checkErrorIsNil(err)
			filename := msg.MsgID + ".mp3"
			saveFile(filename, body)
			fmt.Println("voice save to " + filename)
			break
		}
	}
	// 测试接收视频消息
	fmt.Println("please send a video message via " + contact.NickName)
	for {
		msg := <-addMessageChan
		if msg.MsgType == datastruct.LittleVideoMsg {
			body, err := client.SaveMessageVideo(msg.MsgID)
			checkErrorIsNil(err)
			filename := msg.MsgID + ".mp4"
			saveFile(filename, body)
			fmt.Println("video save to " + filename)
			// 发送回去
			mediaID, _, err := client.UploadMedia(user.UserName, contact.UserName, filename, body)
			checkErrorIsNil(err)
			client.SendVideoMessage(user.UserName, contact.UserName, mediaID)
			break
		}
	}
	// 测试接收动图消息
	fmt.Println("please send a emoticon message via " + contact.NickName)
	for {
		msg := <-addMessageChan
		if msg.MsgType == datastruct.AnimationEmotionsMsg {
			body, err := client.SaveMessageImage(msg.MsgID)
			checkErrorIsNil(err)
			filename := msg.MsgID + ".gif"
			saveFile(filename, body)
			fmt.Println("emoticon save to " + filename)
			// 发送回去
			mediaID, _, err := client.UploadMedia(user.UserName, contact.UserName, filename, body)
			checkErrorIsNil(err)
			client.SendEmoticonMessage(user.UserName, contact.UserName, mediaID)
			break
		}
	}
	// 测试接收文件消息
	fmt.Println("please send a file message via " + contact.NickName)
	for {
		msg := <-addMessageChan
		if msg.MsgType == datastruct.AppMsg {
			body, err := client.SaveMessageFile(msg.FromUserName, msg.MediaID, msg.FileName)
			checkErrorIsNil(err)
			saveFile(msg.FileName, body)
			fmt.Println("file save to " + msg.FileName)
			// 发送回去
			mediaID, _, err := client.UploadMedia(user.UserName, contact.UserName, msg.FileName, body)
			checkErrorIsNil(err)
			client.SendFileMessage(user.UserName, contact.UserName, mediaID, msg.FileName, int64(len(body)))
			break
		}
	}
}
