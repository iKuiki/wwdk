package datastruct

type ContactFlag int64

const (
	ContactFlagStar            ContactFlag = 64
	ContactFlagNoShareMoments  ContactFlag = 256
	ContactFlagQuiet           ContactFlag = 512
	ContactFlagTop             ContactFlag = 2048
	ContactFlagNoFollowMoments ContactFlag = 65536
)

type Contact struct {
	Alias            string        `json:"Alias"` // 微信号
	AppAccountFlag   int64         `json:"AppAccountFlag"`
	AttrStatus       int64         `json:"AttrStatus"`
	ChatRoomID       int64         `json:"ChatRoomId"`
	City             string        `json:"City"`
	ContactFlag      int64         `json:"ContactFlag"`
	DisplayName      string        `json:"DisplayName"`
	EncryChatRoomID  string        `json:"EncryChatRoomId"`
	HeadImgURL       string        `json:"HeadImgUrl"`
	HideInputBarFlag int64         `json:"HideInputBarFlag"`
	IsOwner          int64         `json:"IsOwner"`
	KeyWord          string        `json:"KeyWord"`
	MemberCount      int64         `json:"MemberCount"`
	MemberList       []interface{} `json:"MemberList"`
	NickName         string        `json:"NickName"` // 用户昵称
	OwnerUin         int64         `json:"OwnerUin"`
	PYInitial        string        `json:"PYInitial"`
	PYQuanPin        string        `json:"PYQuanPin"`
	Province         string        `json:"Province"`
	RemarkName       string        `json:"RemarkName"`      // 备注名称
	RemarkPYInitial  string        `json:"RemarkPYInitial"` // 拼音首字母
	RemarkPYQuanPin  string        `json:"RemarkPYQuanPin"` // 拼音全拼
	Sex              int64         `json:"Sex"`             // 性别，1男2女
	Signature        string        `json:"Signature"`       // 个性签名
	SnsFlag          int64         `json:"SnsFlag"`
	StarFriend       int64         `json:"StarFriend"` // 是否星标好友，1是0否
	Statues          int64         `json:"Statues"`
	Uin              int64         `json:"Uin"`
	UniFriend        int64         `json:"UniFriend"`
	UserName         string        `json:"UserName"`
	VerifyFlag       int64         `json:"VerifyFlag"`
}

func (this Contact) IsStar() bool {
	return this.ContactFlag&int64(ContactFlagStar) > 0
}
func (this Contact) IsNoShareMoments() bool {
	return this.ContactFlag&int64(ContactFlagNoShareMoments) > 0
}
func (this Contact) IsQuiet() bool {
	return this.ContactFlag&int64(ContactFlagQuiet) > 0
}
func (this Contact) IsTop() bool {
	return this.ContactFlag&int64(ContactFlagTop) > 0
}
func (this Contact) IsNoFollowMoments() bool {
	return this.ContactFlag&int64(ContactFlagStar) > 0
}
