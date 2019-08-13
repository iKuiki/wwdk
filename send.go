package wwdk

import (
	"github.com/getsentry/sentry-go"
)

// StatusNotify 消息已读通知
func (wxwb *WechatWeb) StatusNotify(toUserName string, code int64) (err error) {
	body, err := wxwb.api.StatusNotify(wxwb.userInfo.user.UserName, toUserName, code)
	if err != nil {
		wxwb.captureException(err, "fatal", sentry.LevelError, extraData{"body", string(body)})
		return
	}
	return
}

// SendTextMessage 发送消息
func (wxwb *WechatWeb) SendTextMessage(toUserName, content string) (msgID, localID string, err error) {
	msgID, localID, body, err := wxwb.api.SendTextMessage(wxwb.userInfo.user.UserName, toUserName, content)
	if err != nil {
		wxwb.captureException(err, "SendTextMessage fatal", sentry.LevelError, extraData{"body", string(body)})
		return
	}
	wxwb.runInfo.MessageCount++
	wxwb.runInfo.MessageSentCount++
	return
}

// SendRevokeMessage 撤回消息
func (wxwb *WechatWeb) SendRevokeMessage(svrMsgID, clientMsgID, toUserName string) (err error) {
	body, err := wxwb.api.SendRevokeMessage(toUserName, svrMsgID, clientMsgID)
	if err != nil {
		wxwb.captureException(err, "SendRevokeMessage fatal", sentry.LevelError, extraData{"body", string(body)})
		return
	}
	wxwb.runInfo.MessageRevokeCount++
	wxwb.runInfo.MessageRevokeSentCount++
	return
}

// UpdateUserRemakName 修改用户备注
func (wxwb *WechatWeb) UpdateUserRemakName(userName, remarkName string) (err error) {
	body, err := wxwb.api.UpdateUserRemakName(userName, remarkName)
	if err != nil {
		wxwb.captureException(err, "UpdateUserRemakName fatal", sentry.LevelError, extraData{"body", string(body)})
		return
	}
	return
}

// UpdateChatRoomTopic 修改群名
func (wxwb *WechatWeb) UpdateChatRoomTopic(userName, newTopic string) (err error) {
	body, err := wxwb.api.UpdateChatRoomTopic(userName, newTopic)
	if err != nil {
		wxwb.captureException(err, "UpdateChatRoomTopic fatal", sentry.LevelError, extraData{"body", string(body)})
		return
	}
	return
}
