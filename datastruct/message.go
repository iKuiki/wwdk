package datastruct

type MessageType int64

const (
	TEXT_MSG               MessageType = 1     // 文字消息
	PICTURE_MSG            MessageType = 3     // 图片消息
	AUDIO_MSG              MessageType = 34    // 音频消息
	CONTACT_CARD_MSG       MessageType = 42    // 名片
	LITTLE_VIDEO_MSG       MessageType = 43    // 小视频消息
	ANIMATION_EMOTIONS_MSG MessageType = 47    // 动画表情
	LINK_MSG               MessageType = 49    // 链接消息类型，已知有转账、开始共享实时位置、合并转发聊天记录
	APPEND_MSG             MessageType = 10000 // 拓展消息类型，已知有红包、停止共享实时位置
	RECALL_MSG             MessageType = 10002 // 撤回消息
)

type Message struct {
	AppInfo *struct {
		AppID string `json:"AppID"`
		Type  int64  `json:"Type"`
	} `json:"AppInfo"`
	AppMsgType    int64       `json:"AppMsgType"`
	Content       string      `json:"Content"`
	CreateTime    int64       `json:"CreateTime"`
	FileName      string      `json:"FileName"`
	FileSize      string      `json:"FileSize"`
	ForwardFlag   int64       `json:"ForwardFlag"`
	FromUserName  string      `json:"FromUserName"`
	HasProductID  int64       `json:"HasProductId"`
	ImgHeight     int64       `json:"ImgHeight"`
	ImgStatus     int64       `json:"ImgStatus"`
	ImgWidth      int64       `json:"ImgWidth"`
	MediaID       string      `json:"MediaId"`
	MsgID         string      `json:"MsgId"`
	MsgType       MessageType `json:"MsgType"`
	NewMsgID      int64       `json:"NewMsgId"`
	OriContent    string      `json:"OriContent"`
	PlayLength    int64       `json:"PlayLength"`
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
