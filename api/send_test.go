package api_test

import (
	"fmt"
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
