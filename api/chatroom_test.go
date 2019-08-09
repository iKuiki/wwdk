package api_test

import (
	"fmt"
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

// 测试拉群、在群内拉新人、移除群成员的操作
func TestUpdateChatroom(t *testing.T) {
	// 测试根据联系人的userNames拉一个新群的功能
	// 获取2个联系人来拉群
	contact1, skip := getTestContact("TestUpdateChatroom get first member", friendContact)
	if skip {
		t.SkipNow()
	}
	contact2, skip := getTestContact("TestUpdateChatroom get second member", friendContact)
	if skip {
		t.SkipNow()
	}
	memberUsernames := []string{
		contact1.UserName,
		contact2.UserName,
	}
	chatRoomUserName, _, err := client.CreateChatRoom("test chatroom", memberUsernames)
	checkErrorIsNil(err)
	t.Log("Create chatroom successful, chatRoomUserName is ", chatRoomUserName)
	// 获取多一个联系人来测试加人进群的操作
	contact3, skip := getTestContact("TestUpdateChatroom get third member", friendContact)
	if skip {
		t.SkipNow()
	}
	_, err = client.UpdateChatRoomAddMember(chatRoomUserName, contact3.UserName)
	checkErrorIsNil(err)
	fmt.Println("sleep 5 second...")
	// 先休息5秒再移除
	time.Sleep(5 * time.Second)
	t.Log("UpdateChatroom AddMember successful")
	// 测试将刚刚拉进来的人踢出去
	_, err = client.UpdateChatRoomDelMember(chatRoomUserName, contact3.UserName)
	checkErrorIsNil(err)
}
