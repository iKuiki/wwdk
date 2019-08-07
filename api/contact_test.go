package api_test

import (
	"github.com/ikuiki/wwdk/datastruct"
	"testing"
	"time"
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
