package datastruct

// BaseResponse 服务器响应的基本结构体，通用
type BaseResponse struct {
	ErrMsg string `json:"ErrMsg"`
	Ret    int64  `json:"Ret"`
}

// GetCookieRespond 获取Cookie的返回
type GetCookieRespond struct {
	Ret         int64  `xml:"ret"`
	Message     string `xml:"message"`
	Skey        string `xml:"skey"`
	Wxsid       string `xml:"wxsid"`
	Wxuin       string `xml:"wxuin"`
	PassTicket  string `xml:"pass_ticket"`
	Isgrayscale int64  `xml:"isgrayscale"`
}

// WxInitRespond 初始化请求的返回
type WxInitRespond struct {
	BaseResponse        *BaseResponse    `json:"BaseResponse"`
	ChatSet             string           `json:"ChatSet"`
	ClickReportInterval int64            `json:"ClickReportInterval"`
	ClientVersion       int64            `json:"ClientVersion"`
	ContactList         []Contact        `json:"ContactList"`
	Count               int64            `json:"Count"`
	GrayScale           int64            `json:"GrayScale"`
	InviteStartCount    int64            `json:"InviteStartCount"`
	MPSubscribeMsgCount int64            `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []MPSubscribeMsg `json:"MPSubscribeMsgList"`
	SKey                string           `json:"SKey"`
	SyncKey             *SyncKey         `json:"SyncKey"`
	SystemTime          int64            `json:"SystemTime"`
	User                *User            `json:"User"`
}

// GetContactRespond 获取联系人的返回
type GetContactRespond struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	MemberCount  int64         `json:"MemberCount"`
	MemberList   []Contact     `json:"MemberList"`
	Seq          int64         `json:"Seq"`
}

// GetBatchContactResponse 获取群组联系人的返回
type GetBatchContactResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	ContactList  []Contact     `json:"ContactList"`
	Count        int64         `json:"Count"`
}

// SyncCheckRespond 同步消息轮询的返回
type SyncCheckRespond struct {
	Retcode  string `json:" retcode"`
	Selector string `json:"selector"`
}

// GetMessageRespond 取回消息的返回
type GetMessageRespond struct {
	BaseResponse    *BaseResponse `json:"BaseResponse"`
	AddMsgCount     int64         `json:"AddMsgCount"`
	AddMsgList      []Message     `json:"AddMsgList"`
	ContinueFlag    int64         `json:"ContinueFlag"`
	DelContactCount int64         `json:"DelContactCount"`
	DelContactList  []struct {
		ContactFlag int64  `json:"ContactFlag"`
		UserName    string `json:"UserName"`
	} `json:"DelContactList"`
	ModChatRoomMemberCount int64         `json:"ModChatRoomMemberCount"`
	ModChatRoomMemberList  []interface{} `json:"ModChatRoomMemberList"`
	ModContactCount        int64         `json:"ModContactCount"`
	ModContactList         []Contact     `json:"ModContactList"`
	Profile                *struct {
		Alias     string `json:"Alias"`
		BindEmail struct {
			Buff string `json:"Buff"`
		} `json:"BindEmail"`
		BindMobile struct {
			Buff string `json:"Buff"`
		} `json:"BindMobile"`
		BindUin           int64  `json:"BindUin"`
		BitFlag           int64  `json:"BitFlag"`
		HeadImgUpdateFlag int64  `json:"HeadImgUpdateFlag"`
		HeadImgURL        string `json:"HeadImgUrl"`
		NickName          *struct {
			Buff string `json:"Buff"`
		} `json:"NickName"`
		PersonalCard int64  `json:"PersonalCard"`
		Sex          int64  `json:"Sex"`
		Signature    string `json:"Signature"`
		Status       int64  `json:"Status"`
		UserName     *struct {
			Buff string `json:"Buff"`
		} `json:"UserName"`
	} `json:"Profile"`
	SKey         string   `json:"SKey"`
	SyncCheckKey *SyncKey `json:"SyncCheckKey"`
	SyncKey      *SyncKey `json:"SyncKey"`
}

// StatusNotifyRespond 状态通知请求的返回
type StatusNotifyRespond struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	MsgID        string        `json:"MsgID"`
}

// SendMessageRespond 发送消息的返回
type SendMessageRespond struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	LocalID      string        `json:"LocalID"`
	MsgID        string        `json:"MsgID"`
}

// RevokeMessageRespond 撤回消息的返回
type RevokeMessageRespond struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	Introduction string        `json:"Introduction"`
	SysWording   string        `json:"SysWording"`
}

// ModifyRemarkRespond 修改用户备注的返回
type ModifyRemarkRespond struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
}

// ModifyChatRoomTopicRespond 修改群名的返回
type ModifyChatRoomTopicRespond struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	MemberCount  int64         `json:"MemberCount"`
	MemberList   []interface{} `json:"MemberList"`
}
