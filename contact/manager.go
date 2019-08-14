package contact

import (
	stdErrors "errors"
	"strings"
	"sync"

	"github.com/ikuiki/go-component/language"
	"github.com/ikuiki/wwdk/datastruct"
)

// Manager 联系人管理器
// 负责托管web微信的联系人
type Manager interface {
	// SetContact 设置联系人
	// @param contacts 联系人列表(Contact结构体形式)
	// @return err 没有错误则为空
	SetContact(contacts ...datastruct.Contact) (err error)
	// DelContact 根据userName删除联系人(也会同时从好友列表中移除)
	// @param userNames 要删除的联系人的userNames
	// @return deletedContacts 被删除的联系人(如存在)
	// @return err 没有错误则为空
	DelContact(userNames ...string) (deletedContacts []datastruct.Contact, err error)
	// SetFriend 设置好友(设置联系人并且添加到好友列表)
	// @param contacts 联系人列表(Contact结构体形式)
	// @return err 没有错误则为空
	SetFriend(contacts ...datastruct.Contact) (err error)
	// DelFriend 根据userName删除好友(不删除联系人)
	// @param userNames 要删除的联系人的userNames
	// @return deletedContacts 被删除的联系人(如存在)
	// @return err 没有错误则为空
	DelFriend(userNames ...string) (deletedContacts []datastruct.Contact, err error)
	// GetFriendList 获取好友列表
	// @return list 好友列表(Contact结构体形式)
	GetFriendList() (list []datastruct.Contact, err error)
	// GetFriend 通过userName获取contact
	// @param userName 要获取的联系人的userName
	// @return contact 联系人(Contact结构体形式)
	// @return err 如果查无此人返回ErrNotFound
	GetFriend(userName string) (contact datastruct.Contact, err error)
	// GetFriends 通过userName获取contact
	// @param userName 要获取的联系人的userName列表
	// @return contact 联系人列表(Contact结构体形式)
	// @return allFound 所有联系人都查找到了
	// @return err 没有错误则为空
	GetFriends(userNames ...string) (contacts []datastruct.Contact, allFound bool, err error)
	// GetFriendByNickName 通过nickname获取contact
	// @return contact 联系人(Contact结构体形式)
	// @return err 如果查无此人返回ErrNotFound
	GetFriendByNickName(nickname string) (contact datastruct.Contact, err error)
	// GetFriendByRemarkName 通过remarkName获取contact
	// @return contact 联系人(Contact结构体形式)
	// @return err 如果查无此人返回ErrNotFound
	GetFriendByRemarkName(remarkName string) (contact datastruct.Contact, err error)
}

// NewManager 创建联系人管理器
func NewManager() (contactManager Manager, err error) {
	contactManager = &manager{
		contactMap: make(map[string]datastruct.Contact),
	}
	return
}

var (
	// ErrNotFound 错误：未找到
	ErrNotFound = stdErrors.New("not found")
)

type manager struct {
	contactMap    map[string]datastruct.Contact // 联系人列表
	contactLocker sync.RWMutex                  // contactMap的读写锁
	friendList    []string                      // 好友列表
	friendLocker  sync.RWMutex                  // friendList的读写锁
}

// SetContact 设置联系人
// @param contacts 联系人列表(Contact结构体形式)
// @return err 没有错误则为空
func (mgr *manager) SetContact(contacts ...datastruct.Contact) (err error) {
	if len(contacts) == 0 {
		return
	}
	mgr.contactLocker.Lock()
	defer mgr.contactLocker.Unlock()
	for _, contact := range contacts {
		mgr.contactMap[contact.UserName] = contact
	}
	return
}

// DelContact 根据userName删除联系人(也会同时从好友列表中移除)
// @param userNames 要删除的联系人的userNames
// @return deletedContacts 被删除的联系人(如存在)
// @return err 没有错误则为空
func (mgr *manager) DelContact(userNames ...string) (deletedContacts []datastruct.Contact, err error) {
	if len(userNames) == 0 {
		return
	}
	mgr.contactLocker.Lock()
	defer mgr.contactLocker.Unlock()
	for _, userName := range userNames {
		if contact, ok := mgr.contactMap[userName]; ok {
			deletedContacts = append(deletedContacts, contact)
		}
		delete(mgr.contactMap, userName)
	}
	// 也要删除好友
	mgr.friendLocker.Lock()
	defer mgr.friendLocker.Unlock()
	mgr.friendList = language.ArrayDiff(mgr.friendList, userNames).([]string)
	return
}

// SetFriend 设置好友(设置联系人并且添加到好友列表)
// @param contacts 联系人列表(Contact结构体形式)
// @return err 没有错误则为空
func (mgr *manager) SetFriend(contacts ...datastruct.Contact) (err error) {
	if len(contacts) == 0 {
		return
	}
	mgr.contactLocker.Lock()
	defer mgr.contactLocker.Unlock()
	mgr.friendLocker.Lock()
	defer mgr.friendLocker.Unlock()
	for _, contact := range contacts {
		mgr.contactMap[contact.UserName] = contact
		mgr.friendList = append(mgr.friendList, contact.UserName)
	}
	mgr.friendList = language.ArrayUnique(mgr.friendList).([]string)
	return
}

// DelFriend 根据userName删除好友(不删除联系人)
// @param userNames 要删除的好友的userNames
// @return deletedContacts 被删除的好友(如存在)
// @return err 没有错误则为空
func (mgr *manager) DelFriend(userNames ...string) (deletedContacts []datastruct.Contact, err error) {
	if len(userNames) == 0 {
		return
	}
	mgr.contactLocker.RLock()
	defer mgr.contactLocker.RUnlock()
	for _, userName := range userNames {
		if contact, ok := mgr.contactMap[userName]; ok {
			deletedContacts = append(deletedContacts, contact)
		}
	}
	// 删除好友
	mgr.friendLocker.Lock()
	defer mgr.friendLocker.Unlock()
	mgr.friendList = language.ArrayDiff(mgr.friendList, userNames).([]string)
	return
}

// 根据contactMap填充聊天室中的成员
// 注意！调用此方法前请自行Lock用到的资源：contactMap
// @param chatroom 要填充的聊天室
// @return 填充好的聊天室
func (mgr *manager) fillChatroomMember(chatroom datastruct.Contact) datastruct.Contact {
	for k, member := range chatroom.MemberList {
		if contact, ok := mgr.contactMap[member.UserName]; ok {
			chatroom.MemberList[k] = datastruct.Member{
				AttrStatus:      contact.AttrStatus,
				DisplayName:     contact.DisplayName,
				KeyWord:         contact.KeyWord,
				NickName:        contact.NickName,
				PYInitial:       contact.PYInitial,
				PYQuanPin:       contact.PYQuanPin,
				RemarkPYInitial: contact.RemarkPYInitial,
				RemarkPYQuanPin: contact.RemarkPYQuanPin,
				Uin:             contact.Uin,
				UserName:        contact.UserName,
				// 此MemberStatus不知为何,还是根据原样填回
				MemberStatus: member.MemberStatus,
			}
		}
	}
	return chatroom
}

// GetFriendList 获取好友列表
// @return list 好友列表(Contact结构体形式)
func (mgr *manager) GetFriendList() (list []datastruct.Contact, err error) {
	mgr.friendLocker.RLock()
	defer mgr.friendLocker.RUnlock()
	mgr.contactLocker.RLock()
	defer mgr.contactLocker.RUnlock()
	for _, userName := range mgr.friendList {
		if contact, ok := mgr.contactMap[userName]; ok {
			if contact.IsChatroom() {
				contact = mgr.fillChatroomMember(contact)
			}
			list = append(list, contact)
		}
	}
	return
}

// GetFriend 通过userName获取contact
// @param userName 要获取的联系人的userName
// @return contact 联系人(Contact结构体形式)
// @return err 如果查无此人返回ErrNotFound
func (mgr *manager) GetFriend(userName string) (contact datastruct.Contact, err error) {
	mgr.friendLocker.RLock()
	defer mgr.friendLocker.RUnlock()
	if language.ArrayIn(mgr.friendList, userName) == -1 {
		// 查无此好友
		err = ErrNotFound
		return
	}
	mgr.contactLocker.RLock()
	defer mgr.contactLocker.RUnlock()
	contact, ok := mgr.contactMap[userName]
	if !ok {
		// 查无此联系人
		err = ErrNotFound
		return
	}
	if contact.IsChatroom() {
		contact = mgr.fillChatroomMember(contact)
	}
	return
}

// GetFriends 通过userName获取contact
// @param userName 要获取的好友的userName列表
// @return contact 联系人列表(Contact结构体形式)
// @return allFound 所有好友都查找到了
// @return err 没有错误则为空
func (mgr *manager) GetFriends(userNames ...string) (contacts []datastruct.Contact, allFound bool, err error) {
	allFound = true
	if len(userNames) == 0 {
		return
	}
	mgr.friendLocker.RLock()
	defer mgr.friendLocker.RUnlock()
	mgr.contactLocker.RLock()
	defer mgr.contactLocker.RUnlock()
	for _, userName := range userNames {
		if language.ArrayIn(mgr.friendList, userName) != -1 {
			if contact, ok := mgr.contactMap[userName]; ok {
				if contact.IsChatroom() {
					contact = mgr.fillChatroomMember(contact)
				}
				contacts = append(contacts, contact)
			}
		}
	}
	// 检查是否全部发现
	allFound = len(contacts) == len(userNames)
	return
}

// GetFriendByNickName 通过nickname查找好友
// @return contact 联系人(Contact结构体形式)
// @return err 如果查无此人返回ErrNotFound
func (mgr *manager) GetFriendByNickName(nickname string) (contact datastruct.Contact, err error) {
	mgr.friendLocker.RLock()
	defer mgr.friendLocker.RUnlock()
	mgr.contactLocker.RLock()
	defer mgr.contactLocker.RUnlock()
	for _, userName := range mgr.friendList { // 之所以对friendList执行for each而非对contactMap执行，是因为contactList太大了
		if contact, ok := mgr.contactMap[userName]; ok {
			if strings.Contains(strings.ToLower(contact.NickName), strings.ToLower(nickname)) {
				return contact, nil
			}
		}
	}
	err = ErrNotFound
	return
}

// GetFriendByRemarkName 通过remarkName查找好友
// @return contact 联系人(Contact结构体形式)
// @return err 如果查无此人返回ErrNotFound
func (mgr *manager) GetFriendByRemarkName(remarkName string) (contact datastruct.Contact, err error) {
	mgr.friendLocker.RLock()
	defer mgr.friendLocker.RUnlock()
	mgr.contactLocker.RLock()
	defer mgr.contactLocker.RUnlock()
	for _, userName := range mgr.friendList { // 之所以对friendList执行for each而非对contactMap执行，是因为contactList太大了
		if contact, ok := mgr.contactMap[userName]; ok {
			if strings.Contains(strings.ToLower(contact.RemarkName), strings.ToLower(remarkName)) {
				return contact, nil
			}
		}
	}
	err = ErrNotFound
	return
}
