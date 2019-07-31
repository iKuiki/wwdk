package datastruct

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// MessageType 消息类型
type MessageType int64

const (
	// TextMsg 文字消息
	TextMsg MessageType = 1
	// ImageMsg 图片消息
	ImageMsg MessageType = 3
	// VoiceMsg 音频消息
	VoiceMsg MessageType = 34
	// ContactCardMsg 名片
	ContactCardMsg MessageType = 42
	// LittleVideoMsg 小视频消息
	LittleVideoMsg MessageType = 43
	// AnimationEmotionsMsg 动画表情
	AnimationEmotionsMsg MessageType = 47
	// LinkMsg 链接消息类型，已知有转账、开始共享实时位置、合并转发聊天记录
	LinkMsg MessageType = 49
	// AppendMsg 拓展消息类型，已知有红包、停止共享实时位置
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

// Message 微信消息结构体
type Message struct {
	AppInfo *struct {
		AppID string `json:"AppID"`
		Type  int64  `json:"Type"`
	} `json:"AppInfo"`
	AppMsgType    AppMessageType `json:"AppMsgType"`
	Content       string         `json:"Content"`
	CreateTime    int64          `json:"CreateTime"`
	FileName      string         `json:"FileName"`
	FileSize      string         `json:"FileSize"`
	ForwardFlag   int64          `json:"ForwardFlag"`
	FromUserName  string         `json:"FromUserName"`
	HasProductID  int64          `json:"HasProductId"`
	ImgHeight     int64          `json:"ImgHeight"`
	ImgStatus     int64          `json:"ImgStatus"`
	ImgWidth      int64          `json:"ImgWidth"`
	MediaID       string         `json:"MediaId"`
	MsgID         string         `json:"MsgId"`
	MsgType       MessageType    `json:"MsgType"`
	NewMsgID      int64          `json:"NewMsgId"`
	OriContent    string         `json:"OriContent"`
	PlayLength    int64          `json:"PlayLength"`
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

// IsChatroom 返回消息是否为群组消息
func (msg Message) IsChatroom() bool {
	return strings.HasPrefix(msg.FromUserName, "@@")
}

// GetMemberUserName 获取群组消息的发件人
func (msg Message) GetMemberUserName() (userName string, err error) {
	if msg.IsChatroom() {
		splitIndex := strings.Index(msg.Content, ":<br/>")
		if splitIndex == -1 {
			err = errors.New("userName not found")
		} else {
			if match, _ := regexp.MatchString("^@\\w+$", msg.Content[:splitIndex]); match {
				userName = msg.Content[:splitIndex]
			}
		}
	} else {
		err = errors.New("this message is not chatroon message")
	}
	return
}

// GetMemberMsgContent 获取群组消息的内容
func (msg Message) GetMemberMsgContent() (content string, err error) {
	if msg.IsChatroom() {
		splitIndex := strings.Index(msg.Content, ":<br/>")
		if splitIndex == -1 {
			err = errors.New("content not found")
		} else {
			if match, _ := regexp.MatchString("^@\\w+$", msg.Content[:splitIndex]); match {
				content = msg.Content[splitIndex+1:]
			}
		}
	} else {
		err = errors.New("this message is not chatroon message")
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
