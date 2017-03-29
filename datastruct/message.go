package datastruct

type MessageTypd int64

const ()

type Message struct {
	AppInfo *struct {
		AppID string `json:"AppID"`
		Type  int64  `json:"Type"`
	} `json:"AppInfo"`
	AppMsgType    int64  `json:"AppMsgType"`
	Content       string `json:"Content"`
	CreateTime    int64  `json:"CreateTime"`
	FileName      string `json:"FileName"`
	FileSize      string `json:"FileSize"`
	ForwardFlag   int64  `json:"ForwardFlag"`
	FromUserName  string `json:"FromUserName"`
	HasProductID  int64  `json:"HasProductId"`
	ImgHeight     int64  `json:"ImgHeight"`
	ImgStatus     int64  `json:"ImgStatus"`
	ImgWidth      int64  `json:"ImgWidth"`
	MediaID       string `json:"MediaId"`
	MsgID         string `json:"MsgId"`
	MsgType       int64  `json:"MsgType"`
	NewMsgID      int64  `json:"NewMsgId"`
	OriContent    string `json:"OriContent"`
	PlayLength    int64  `json:"PlayLength"`
	RecommendInfo *struct {
		Alias      string `json:"Alias"`
		AttrStatus int64  `json:"AttrStatus"`
		City       string `json:"City"`
		Content    string `json:"Content"`
		NickName   string `json:"NickName"`
		OpCode     int64  `json:"OpCode"`
		Province   string `json:"Province"`
		QQNum      int64  `json:"QQNum"`
		Scene      int64  `json:"Scene"`
		Sex        int64  `json:"Sex"`
		Signature  string `json:"Signature"`
		Ticket     string `json:"Ticket"`
		UserName   string `json:"UserName"`
		VerifyFlag int64  `json:"VerifyFlag"`
	} `json:"RecommendInfo"`
	Status               int64  `json:"Status"`
	StatusNotifyCode     int64  `json:"StatusNotifyCode"`
	StatusNotifyUserName string `json:"StatusNotifyUserName"`
	SubMsgType           int64  `json:"SubMsgType"`
	Ticket               string `json:"Ticket"`
	ToUserName           string `json:"ToUserName"`
	URL                  string `json:"Url"`
	VoiceLength          int64  `json:"VoiceLength"`
}
