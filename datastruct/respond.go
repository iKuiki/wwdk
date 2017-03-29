package datastruct

type BaseResponse struct {
	ErrMsg string `json:"ErrMsg"`
	Ret    int64  `json:"Ret"`
}

type WxInitRespond struct {
	BaseResponse        BaseResponse  `json:"BaseResponse"`
	ChatSet             string        `json:"ChatSet"`
	ClickReportInterval int64         `json:"ClickReportInterval"`
	ClientVersion       int64         `json:"ClientVersion"`
	ContactList         []*Contact    `json:"ContactList"`
	Count               int64         `json:"Count"`
	GrayScale           int64         `json:"GrayScale"`
	InviteStartCount    int64         `json:"InviteStartCount"`
	MPSubscribeMsgCount int64         `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []interface{} `json:"MPSubscribeMsgList"`
	SKey                string        `json:"SKey"`
	SyncKey             *SyncKey      `json:"SyncKey"`
	SystemTime          int64         `json:"SystemTime"`
	User                *User         `json:"User"`
}

type GetContactRespond struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MemberCount  int64        `json:"MemberCount"`
	MemberList   []*Contact   `json:"MemberList"`
	Seq          int64        `json:"Seq"`
}

type SyncCheckRespond struct {
	Retcode  string `json:" retcode"`
	Selector string `json:"selector"`
}
