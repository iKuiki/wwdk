package datastruct

import (
	"errors"
	"strings"
)

// ContactFlag 联系人标志
type ContactFlag int64

const (
	// ContactFlagStar 星标联系人
	ContactFlagStar ContactFlag = 64
	// ContactFlagNoShareMoments 不分享朋友圈的联系人
	ContactFlagNoShareMoments ContactFlag = 256
	// ContactFlagQuiet 屏蔽新消息提示的联系人
	ContactFlagQuiet ContactFlag = 512
	// ContactFlagTop 置顶联系人
	ContactFlagTop ContactFlag = 2048
	// ContactFlagNoFollowMoments 不查看其朋友圈的联系人
	ContactFlagNoFollowMoments ContactFlag = 65536
)

// Member 群成员
type Member struct {
	AttrStatus      int64  `json:"AttrStatus"`
	DisplayName     string `json:"DisplayName"`
	KeyWord         string `json:"KeyWord"`
	MemberStatus    int64  `json:"MemberStatus"`
	NickName        string `json:"NickName"`
	PYInitial       string `json:"PYInitial"`
	PYQuanPin       string `json:"PYQuanPin"`
	RemarkPYInitial string `json:"RemarkPYInitial"`
	RemarkPYQuanPin string `json:"RemarkPYQuanPin"`
	Uin             int64  `json:"Uin"`
	UserName        string `json:"UserName"`
}

// Contact 联系人结构
type Contact struct {
	Alias            string   `json:"Alias"` // 微信号
	AppAccountFlag   int64    `json:"AppAccountFlag"`
	AttrStatus       int64    `json:"AttrStatus"`
	ChatRoomID       int64    `json:"ChatRoomId"`
	City             string   `json:"City"`
	ContactFlag      int64    `json:"ContactFlag"`
	DisplayName      string   `json:"DisplayName"`
	EncryChatRoomID  string   `json:"EncryChatRoomId"`
	HeadImgURL       string   `json:"HeadImgUrl"`
	HideInputBarFlag int64    `json:"HideInputBarFlag"`
	IsOwner          int64    `json:"IsOwner"`
	KeyWord          string   `json:"KeyWord"`
	MemberCount      int64    `json:"MemberCount"`
	MemberList       []Member `json:"MemberList"`
	NickName         string   `json:"NickName"` // 用户昵称
	OwnerUin         int64    `json:"OwnerUin"`
	PYInitial        string   `json:"PYInitial"`
	PYQuanPin        string   `json:"PYQuanPin"`
	Province         string   `json:"Province"`
	RemarkName       string   `json:"RemarkName"`      // 备注名称
	RemarkPYInitial  string   `json:"RemarkPYInitial"` // 拼音首字母
	RemarkPYQuanPin  string   `json:"RemarkPYQuanPin"` // 拼音全拼
	Sex              int64    `json:"Sex"`             // 性别，1男2女
	Signature        string   `json:"Signature"`       // 个性签名
	SnsFlag          int64    `json:"SnsFlag"`
	StarFriend       int64    `json:"StarFriend"` // 是否星标好友，1是0否
	Statues          int64    `json:"Statues"`
	Uin              int64    `json:"Uin"`
	UniFriend        int64    `json:"UniFriend"`
	UserName         string   `json:"UserName"` // 用户标识，收发信息都以此为依据，个体用户以@开头(包括公众号)，群组以@@开头
	VerifyFlag       int64    `json:"VerifyFlag"`
}

// IsStar 返回联系人是否为星标联系人
func (contact Contact) IsStar() bool {
	return contact.ContactFlag&int64(ContactFlagStar) > 0
}

// IsNoShareMoments 返回是否为不分享朋友圈的联系人
func (contact Contact) IsNoShareMoments() bool {
	return contact.ContactFlag&int64(ContactFlagNoShareMoments) > 0
}

// IsQuiet 返回是否屏蔽该联系人的新消息提醒
func (contact Contact) IsQuiet() bool {
	return contact.ContactFlag&int64(ContactFlagQuiet) > 0
}

// IsTop 返回该联系人是否置顶
func (contact Contact) IsTop() bool {
	return contact.ContactFlag&int64(ContactFlagTop) > 0
}

// IsNoFollowMoments 返回是否不查看该联系人的朋友圈
func (contact Contact) IsNoFollowMoments() bool {
	return contact.ContactFlag&int64(ContactFlagStar) > 0
}

// IsChatroom 返回联系人是否为群组
func (contact Contact) IsChatroom() bool {
	return strings.HasPrefix(contact.UserName, "@@")
}

// GetMember 根据userName获取成员
func (contact Contact) GetMember(userName string) (member Member, err error) {
	found := false
	for _, m := range contact.MemberList {
		if m.UserName == userName {
			member = m
			found = true
		}
	}
	if !found {
		err = errors.New("member not found")
	}
	return
}
