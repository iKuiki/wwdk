package api_test

import (
	"testing"
)

// 测试修改一个群聊的标题
// 为了不影响他人，需要先生成一个随机数，让对应联系人发过来确认
// 期望：返回err=nil
func TestModifyChatroomTopic(t *testing.T) {
	contact, skip := getTestContact("TestModifyChatroomTopic", true)
	if !skip {
		_, err := client.ModifyChatRoomTopic(contact.UserName, contact.NickName+"2")
		checkErrorIsNil(err)
	}
}
