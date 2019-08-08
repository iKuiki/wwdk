package api_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
)

// 测试获取联系人的方法
// 期望：获取到至少一个联系人
func TestGetContact(t *testing.T) {
	contacts, _, err := client.GetContact()
	checkErrorIsNil(err)
	if len(contacts) == 0 {
		t.Fatal("getContact result empty")
	}
	t.Logf("got %d contacts", len(contacts))
	for _, c := range contacts {
		// 日常维护contactMap
		contactMap[c.UserName] = c
	}
}

// 测试批量获取联系人
// 期望：获取到和提交的联系人数量一致
func TestBatchGetContact(t *testing.T) {
	if len(contactMap) == 0 {
		t.Fatal("contact length is 0")
	}
	var itemList []datastruct.BatchGetContactRequestListItem
	for _, contact := range contactMap {
		if len(itemList) > 7 {
			// 用7个人做测试就OK了
			break
		}
		itemList = append(itemList, datastruct.BatchGetContactRequestListItem{
			UserName: contact.UserName})
	}
	contactList, _, err := client.BatchGetContact(itemList)
	checkErrorIsNil(err)
	t.Logf("batch got %d contacts", len(contactList))
	if len(contactList) != len(itemList) {
		t.Fatalf("batchGetContact list len(%d) diff with contacts(%d)", len(contactList), len(itemList))
	}
	for _, c := range contactList {
		// 日常维护contactMap
		contactMap[c.UserName] = c
	}
}

// 测试修改一个联系人的备注
// 为了不影响他人，需要先生成一个随机数，让对应联系人发过来确认
// 期望：返回err=nil
func TestModifyContactRemark(t *testing.T) {
	contact, skip := getTestContact("TestModifyContactRemark", friendContact)
	if skip {
		t.SkipNow()
	}
	remark := contact.RemarkName
	if remark == "" {
		remark = contact.NickName
	}
	_, err := client.ModifyUserRemakName(contact.UserName, remark+"2")
	checkErrorIsNil(err)
	for {
		select {
		case mCon := <-modContactChan:
			// 日常维护contactMap
			contactMap[mCon.UserName] = mCon
			if mCon.UserName == contact.UserName {
				if mCon.RemarkName != remark+"2" {
					t.Fatalf("modify user remark fail, except %s, got %s\n",
						remark+"2",
						mCon.RemarkName,
					)
				}
				t.Logf("valid user remark has change to %s\n", mCon.RemarkName)
				return
			}
		case <-time.After(5 * time.Second):
			t.Fatal("wait for modify notify timeout")
		}
	}
}

// TestAcceptAddFriend 测试同意添加好友请求
// 首先需要收到好友请求并且获取其中的username与ticket
// 然后调用接受好友请求
// 然后期待接受到获取到新好友的消息
func TestAcceptAddFriend(t *testing.T) {
	validCode := tool.GetRandomStringFromNum(4)
	fmt.Println("please use a wechat add " + user.NickName + " to contact list")
	fmt.Printf("if you want to Skip test this func please send [skip %s]\n", validCode)
	var recommendInfo datastruct.MessageRecommendInfo
	for {
		// wait a add friend message
		msg := <-addMessageChan
		if msg.GetContent() == "skip "+validCode {
			t.SkipNow()
		} else {
			if msg.FromUserName == "fmessage" && msg.MsgType == datastruct.AddFriendMsg {
				recommendInfo = *msg.RecommendInfo
				break
			}
		}
	}
	fmt.Printf("recive addFriend msg [%s]: %s\n",
		recommendInfo.NickName,
		recommendInfo.Content)
	// Sleep 1 second for add friend message show on phone
	time.Sleep(time.Second)
	// now accept add friend
	_, err := client.AcceptAddFriend(recommendInfo.UserName,
		recommendInfo.Ticket)
	checkErrorIsNil(err)
	// 现在期待接到联系人添加的变更
	for {
		select {
		case c := <-modContactChan:
			// 日常维护contactMap
			contactMap[c.UserName] = c
			if c.UserName == recommendInfo.UserName {
				t.Logf("verify new contact [%s] added\n", c.NickName)
				return
			}
			// continue
		case <-time.After(5 * time.Second):
			t.Fatal("not receive contact modify msg")
		}
	}
}
