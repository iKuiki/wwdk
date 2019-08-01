package api_test

import (
	"fmt"
	"testing"
	"time"
)

// 测试发送一次消息并撤回
func TestSendAndRevoke(t *testing.T) {
	contact, skip := getTestContact("TestSendAndRevoke", false)
	if !skip {
		msgID, localID, _, err := client.SendTextMessage(user.UserName, contact.UserName, "test message will be revoke")
		checkErrorIsNil(err)
		fmt.Println("sleep 2 second")
		time.Sleep(2 * time.Second)
		_, err = client.SendRevokeMessage(contact.UserName, msgID, localID)
		checkErrorIsNil(err)
	}
}
