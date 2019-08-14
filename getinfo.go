package wwdk

import (
	"github.com/ikuiki/wwdk/contactmgr"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/runinfo"
	"github.com/pkg/errors"
)

// 此文件内的方法主要为WechatWeb暴露给外部调用获取信息的方法

// GetUser 获取当前登陆用户
func (wxwb *WechatWeb) GetUser() (user datastruct.User, err error) {
	if wxwb.user == nil {
		err = errors.New("User not found")
	} else {
		user = *wxwb.user
	}
	return
}

// GetContact 根据username获取联系人
func (wxwb *WechatWeb) GetContact(username string) (contact datastruct.Contact, err error) {
	contact, err = wxwb.contactManager.GetFriend(username)
	if err == contactmgr.ErrNotFound {
		// 尝试获取一次
		contactList, _, _ := wxwb.api.BatchGetContact([]datastruct.BatchGetContactRequestListItem{
			datastruct.BatchGetContactRequestListItem{
				UserName: username, // 因为只填写了UserName,所以获取的一定是Friend
			},
		})
		wxwb.contactManager.SetFriend(contactList...)
		contact, err = wxwb.contactManager.GetFriend(username)
		if err == nil {
			wxwb.logger.Infof("User %s not found, but BatchGetContact got that", contact.NickName)
			wxwb.syncChannel <- SyncChannelItem{
				Code:    SyncStatusModifyContact,
				Contact: &contact,
			}
		}
	}
	return
}

// GetContactByAlias 根据Alias获取联系人
func (wxwb *WechatWeb) GetContactByAlias(alias string) (contact datastruct.Contact, err error) {
	return wxwb.contactManager.GetFriendByAlias(alias)
}

// GetContactByNickname 根据昵称获取用户名
func (wxwb *WechatWeb) GetContactByNickname(nickname string) (contact datastruct.Contact, err error) {
	return wxwb.contactManager.GetFriendByNickName(nickname)
}

// GetContactByRemarkName 根据备注获取用户名
func (wxwb *WechatWeb) GetContactByRemarkName(remarkName string) (contact datastruct.Contact, err error) {
	return wxwb.contactManager.GetFriendByRemarkName(remarkName)
}

// GetContactList 获取联系人列表
func (wxwb *WechatWeb) GetContactList() (contacts []datastruct.Contact, err error) {
	return wxwb.contactManager.GetFriendList()
}

// GetRunInfo 获取运行计数器信息
func (wxwb *WechatWeb) GetRunInfo() (runinfo runinfo.WechatRunInfo) {
	return wxwb.runInfo
}
