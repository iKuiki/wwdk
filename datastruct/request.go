package datastruct

type BaseRequest struct {
	DeviceID string `json:"DeviceID"`
	Sid      string `json:"Sid"`
	Skey     string `json:"Skey"`
	Uin      string `json:"Uin"`
}

type WxInitRequestBody struct {
	BaseRequest *BaseRequest `json:"BaseRequest"`
}

type GetMessageRequest struct {
	BaseRequest *BaseRequest `json:"BaseRequest"`
	SyncKey     *SyncKey     `json:"SyncKey"`
	Rr          int64        `json:"rr"`
}

type StatusNotifyRequest struct {
	BaseRequest  *BaseRequest `json:"BaseRequest"`
	ClientMsgID  int64        `json:"ClientMsgId"`
	Code         int64        `json:"Code"`
	FromUserName string       `json:"FromUserName"`
	ToUserName   string       `json:"ToUserName"`
}
