package api_test

import (
	"github.com/ikuiki/go-component/language"
	"testing"
	"time"
)

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
	chatRoomUserName, _, err := client.CreateChatRoom("", memberUsernames)
	checkErrorIsNil(err)
	// 此处无需期待联系人变动通知，如果没有新消息到来，此处不会出现联系人变动通知的
	t.Log("Create chatroom successful, chatRoomUserName is ", chatRoomUserName)
	// 获取多一个联系人来测试加人进群的操作
	contact3, skip := getTestContact("TestUpdateChatroom get third member", friendContact)
	if skip {
		t.SkipNow()
	}
	// 期待联系人变动通知
	// 之所以放到获取多一个联系人后再操作是因为获取联系人一定会触发新消息
CHECK_CREATE:
	for {
		select {
		case mCon := <-modContactChan:
			if mCon.UserName == chatRoomUserName {
				var mList []string // 用来记录这个群里的联系人的UserName
				for _, member := range mCon.MemberList {
					mList = append(mList, member.UserName)
				}
				if len(language.ArrayDiff(memberUsernames, mList).([]string)) == 0 {
					break CHECK_CREATE // 创建时指定的成员都在内
				}
				// 如果member列表未见，则认为添加失败
				t.Fatalf("modify user remark fail, except %s, got %s\n",
					"test chatroom",
					mCon.NickName,
				)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("wait for modify notify timeout")
		}
	}
	_, err = client.UpdateChatRoomAddMember(chatRoomUserName, contact3.UserName)
	checkErrorIsNil(err)
	// 期待联系人变动通知
CHECK_ADD:
	for {
		select {
		case mCon := <-modContactChan:
			if mCon.UserName == chatRoomUserName {
				for _, member := range mCon.MemberList {
					if member.UserName == contact3.UserName {
						t.Logf("valid member %s has added\n", contact3.NickName)
						break CHECK_ADD
					}
				}
				// 如果member列表未见，则认为添加失败
				t.Fatalf("modify user remark fail, except %s, got %s\n",
					"test chatroom",
					mCon.NickName,
				)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("wait for modify notify timeout")
		}
	}
	t.Log("UpdateChatroom AddMember successful")
	// 等待2秒，以免截获添加群成员的联系人变动
	time.Sleep(2 * time.Second)
	// 测试修改一个群聊的标题
	_, err = client.UpdateChatRoomTopic(chatRoomUserName, "test chatroom")
	checkErrorIsNil(err)
	// 检查是否收到联系人变动通知
CHECK_TOPIC:
	for {
		select {
		case mCon := <-modContactChan:
			if mCon.UserName == chatRoomUserName {
				if mCon.NickName != "test chatroom" {
					t.Fatalf("modify user remark fail, except %s, got %s\n",
						"test chatroom",
						mCon.NickName,
					)
				}
				t.Logf("valid chatroom topic has change to %s\n", mCon.NickName)
				break CHECK_TOPIC
			}
		case <-time.After(5 * time.Second):
			t.Fatal("wait for modify notify timeout")
		}
	}
}
