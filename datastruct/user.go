package datastruct

// User 当前登陆用户的结构体
type User struct {
	AppAccountFlag    int64  `json:"AppAccountFlag"`
	ContactFlag       int64  `json:"ContactFlag"`
	HeadImgFlag       int64  `json:"HeadImgFlag"`
	HeadImgURL        string `json:"HeadImgUrl"`
	HideInputBarFlag  int64  `json:"HideInputBarFlag"`
	NickName          string `json:"NickName"`
	PYInitial         string `json:"PYInitial"`
	PYQuanPin         string `json:"PYQuanPin"`
	RemarkName        string `json:"RemarkName"`
	RemarkPYInitial   string `json:"RemarkPYInitial"`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin"`
	Sex               int64  `json:"Sex"`
	Signature         string `json:"Signature"`
	SnsFlag           int64  `json:"SnsFlag"`
	StarFriend        int64  `json:"StarFriend"`
	Uin               int64  `json:"Uin"`
	UserName          string `json:"UserName"`
	VerifyFlag        int64  `json:"VerifyFlag"`
	WebWxPluginSwitch int64  `json:"WebWxPluginSwitch"`
}
