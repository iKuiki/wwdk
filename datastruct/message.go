package datastruct

import (
	"html"
	"regexp"
	"strings"

	stdErrors "errors"
)

// MessageType 消息类型
type MessageType int64

const (
	// TextMsg 文字消息
	TextMsg MessageType = 1
	// ImageMsg 图片消息
	ImageMsg MessageType = 3
	// AppMsg app消息(已知的有发送文件的消息
	AppMsg MessageType = 6
	// VoiceMsg 音频消息
	VoiceMsg MessageType = 34
	// AddFriendMsg 收到添加好友请求消息
	AddFriendMsg MessageType = 37
	// ContactCardMsg 名片
	ContactCardMsg MessageType = 42
	// LittleVideoMsg 小视频消息
	LittleVideoMsg MessageType = 43
	// AnimationEmotionsMsg 动画表情
	AnimationEmotionsMsg MessageType = 47
	// LinkMsg 链接消息类型，已知有转账、开始共享实时位置、合并转发聊天记录
	LinkMsg MessageType = 49
	// AppendMsg 灰色无边框文字消息类型，已知有红包、停止共享实时位置
	AppendMsg MessageType = 10000
	// RevokeMsg 撤回消息
	RevokeMsg MessageType = 10002
)

// AppMessageType 应用消息类型
type AppMessageType int64

const (
	// UnknownAppmsg 未知消息类型
	UnknownAppmsg AppMessageType = 0
	// ReciveFileAppmsg 接收文件
	ReciveFileAppmsg AppMessageType = 6
	// StartShareLocation 共享位置
	StartShareLocation AppMessageType = 17
	// ReciveTransferAppmsg 收到转账
	ReciveTransferAppmsg AppMessageType = 2000
)

// 错误声明
var (
	// ErrMemberContactNotFound 消息中未发现成员联系人(此情况一般是消息由用户本人发出)
	ErrMemberContactNotFound = stdErrors.New("member contact not found")
	// ErrNotChatroomMsg 对非chatroom消息调用了chatroom方法
	ErrNotChatroomMsg = stdErrors.New("this message is not chatroon message")
)

// Message 微信消息结构体
type Message struct {
	AppInfo *struct {
		AppID string `json:"AppID"`
		Type  int64  `json:"Type"`
	} `json:"AppInfo"`
	AppMsgType           AppMessageType        `json:"AppMsgType"`
	Content              string                `json:"Content"`
	CreateTime           int64                 `json:"CreateTime"`
	FileName             string                `json:"FileName"`
	FileSize             string                `json:"FileSize"`
	ForwardFlag          int64                 `json:"ForwardFlag"`
	FromUserName         string                `json:"FromUserName"`
	HasProductID         int64                 `json:"HasProductId"`
	ImgHeight            int64                 `json:"ImgHeight"`
	ImgStatus            int64                 `json:"ImgStatus"`
	ImgWidth             int64                 `json:"ImgWidth"`
	MediaID              string                `json:"MediaId"`
	MsgID                string                `json:"MsgId"`
	MsgType              MessageType           `json:"MsgType"`
	NewMsgID             int64                 `json:"NewMsgId"`
	OriContent           string                `json:"OriContent"`
	PlayLength           int64                 `json:"PlayLength"`
	RecommendInfo        *MessageRecommendInfo `json:"RecommendInfo"`
	Status               int64                 `json:"Status"`
	StatusNotifyCode     int64                 `json:"StatusNotifyCode"`
	StatusNotifyUserName string                `json:"StatusNotifyUserName"`
	SubMsgType           int64                 `json:"SubMsgType"`
	Ticket               string                `json:"Ticket"`
	ToUserName           string                `json:"ToUserName"`
	URL                  string                `json:"Url"`
	VoiceLength          int64                 `json:"VoiceLength"`
}

// MessageRecommendInfo 接到添加好友请求时其内的联系人信息
type MessageRecommendInfo struct {
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
}

// IsChatroom 返回消息是否为群组消息
func (msg Message) IsChatroom() bool {
	return strings.HasPrefix(msg.FromUserName, "@@")
}

// GetMemberUserName 获取群组消息的发件人
func (msg Message) GetMemberUserName() (userName string, err error) {
	if msg.IsChatroom() {
		re := regexp.MustCompile("^@\\w+:")
		userName = re.FindString(msg.Content)
		if userName == "" {
			err = ErrMemberContactNotFound
		} else {
			userName = userName[:len(userName)-1]
		}
	} else {
		err = ErrNotChatroomMsg
	}
	return
}

// GetMemberMsgContent 获取群组消息的内容
func (msg Message) GetMemberMsgContent() (content string, err error) {
	if msg.IsChatroom() {
		memberUserName, err := msg.GetMemberUserName()
		if err != nil {
			content = msg.Content
		} else {
			content = strings.TrimPrefix(msg.Content, memberUserName+":")
		}
	} else {
		err = ErrNotChatroomMsg
	}
	return
}

// GetMemberMsgContentUnescape 获取群组消息的内容并解码
func (msg Message) GetMemberMsgContentUnescape() (content string, err error) {
	content, err = msg.GetMemberMsgContent()
	if err == nil {
		content = strings.Replace(html.UnescapeString(content), "<br/>", "\n", -1)
	}
	return
}

// GetContent 获取消息，如果为群消息，则自动尝试获取真实消息本体
func (msg Message) GetContent() (content string) {
	content = msg.Content
	if msg.IsChatroom() {
		content, _ = msg.GetMemberMsgContent()
	}
	return
}

// GetContentUnescape 获取消息本地并解码，如果为群消息则尝试自动获取真实消息本体
func (msg Message) GetContentUnescape() (content string) {
	content = strings.Replace(html.UnescapeString(msg.Content), "<br/>", "\n", -1)
	return
}
