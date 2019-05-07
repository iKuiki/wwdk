package wwdk

import (
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/pkg/errors"
)

// 此文件内的方法主要为WechatWeb暴露给外部调用获取信息的方法

// GetUser 获取当前登陆用户
func (wxwb *WechatWeb) GetUser() (user datastruct.User, err error) {
	if wxwb.userInfo.user == nil {
		err = errors.New("User not found")
	} else {
		user = *wxwb.userInfo.user
	}
	return
}

// GetContact 根据username获取联系人
func (wxwb *WechatWeb) GetContact(username string) (contact datastruct.Contact, err error) {
	contact, ok := wxwb.userInfo.contactList[username]
	if !ok {
		err = errors.New("User not found")
	}
	return
}

// GetContactByAlias 根据Alias获取联系人
func (wxwb *WechatWeb) GetContactByAlias(alias string) (contact datastruct.Contact, err error) {
	found := false
	for _, v := range wxwb.userInfo.contactList {
		if v.Alias == alias {
			contact = v
			found = true
		}
	}
	if !found {
		err = errors.New("User not found")
	}
	return
}

// GetContactByNickname 根据昵称获取用户名
func (wxwb *WechatWeb) GetContactByNickname(nickname string) (contact datastruct.Contact, err error) {
	found := false
	for _, v := range wxwb.userInfo.contactList {
		if v.NickName == nickname {
			contact = v
			found = true
		}
	}
	if !found {
		err = errors.New("User not found")
	}
	return
}

// GetContactByRemarkName 根据备注获取用户名
func (wxwb *WechatWeb) GetContactByRemarkName(remarkName string) (contact datastruct.Contact, err error) {
	found := false
	for _, v := range wxwb.userInfo.contactList {
		if v.RemarkName == remarkName {
			contact = v
			found = true
		}
	}
	if !found {
		err = errors.New("User not found")
	}
	return
}

// GetContactList 获取联系人列表
func (wxwb *WechatWeb) GetContactList() (contacts []datastruct.Contact) {
	for _, v := range wxwb.userInfo.contactList {
		contacts = append(contacts, v)
	}
	return
}

// GetRunInfo 获取运行计数器信息
func (wxwb *WechatWeb) GetRunInfo() (runinfo WechatRunInfo) {
	return wxwb.runInfo
}
