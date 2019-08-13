package runinfo

import "time"

// WechatRunInfo 微信运行信息
type WechatRunInfo struct {
	// StartAt 程序启动的时间
	StartAt time.Time
	// LoginAt 程序登陆的时间
	LoginAt time.Time
	// SyncCount 同步次数
	SyncCount uint64
	// ContactModifyCount 联系人修改计数器
	ContactModifyCount uint64
	// MessageCount 消息计数器
	MessageCount uint64
	// MessageRecivedCount 收到消息计数器
	MessageRecivedCount uint64
	// MessageSentCount 发送消息计数器
	MessageSentCount uint64
	// MessageRevokeCount 撤回消息计数器
	MessageRevokeCount uint64
	// MessageRevokeRecivedCount 收到撤回消息计数器
	MessageRevokeRecivedCount uint64
	// MessageRevokeSentCount 发送撤回消息计数器
	MessageRevokeSentCount uint64
	// PanicCount panic计数器
	PanicCount uint64
}
