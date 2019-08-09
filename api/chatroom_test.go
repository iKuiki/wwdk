package api_test

import (
	"testing"
	"time"
)

// 测试修改一个群聊的标题
// 为了不影响他人，需要先生成一个随机数，让对应联系人发过来确认
// 期望：返回err=nil
func TestModifyChatroomTopic(t *testing.T) {
	contact, skip := getTestContact("TestModifyChatroomTopic", chatroomContact)
	if skip {
		t.SkipNow()
	}
	_, err := client.UpdateChatRoomTopic(contact.UserName, contact.NickName+"2")
	checkErrorIsNil(err)
	for {
		select {
		case mCon := <-modContactChan:
			if mCon.UserName == contact.UserName {
				if mCon.NickName != contact.NickName+"2" {
					t.Fatalf("modify user remark fail, except %s, got %s\n",
						contact.NickName+"2",
						mCon.NickName,
					)
				}
				t.Logf("valid chatroom topic has change to %s\n", mCon.NickName)
				return
			}
		case <-time.After(5 * time.Second):
			t.Fatal("wait for modify notify timeout")
		}
	}

}
