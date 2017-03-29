package datastruct

type Contact struct {
	Alias            string        `json:"Alias"`
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
	NickName         string        `json:"NickName"`
	OwnerUin         int64         `json:"OwnerUin"`
	PYInitial        string        `json:"PYInitial"`
	PYQuanPin        string        `json:"PYQuanPin"`
	Province         string        `json:"Province"`
	RemarkName       string        `json:"RemarkName"`
	RemarkPYInitial  string        `json:"RemarkPYInitial"`
	RemarkPYQuanPin  string        `json:"RemarkPYQuanPin"`
	Sex              int64         `json:"Sex"`
	Signature        string        `json:"Signature"`
	SnsFlag          int64         `json:"SnsFlag"`
	StarFriend       int64         `json:"StarFriend"`
	Statues          int64         `json:"Statues"`
	Uin              int64         `json:"Uin"`
	UniFriend        int64         `json:"UniFriend"`
	UserName         string        `json:"UserName"`
	VerifyFlag       int64         `json:"VerifyFlag"`
}
