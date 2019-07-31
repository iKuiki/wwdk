package wwdk

// StatusNotify 消息已读通知
func (wxwb *WechatWeb) StatusNotify(toUserName string, code int64) (err error) {
	_, err = wxwb.api.StatusNotify(wxwb.userInfo.user.UserName, toUserName, code)
	return
}

// SendTextMessage 发送消息
func (wxwb *WechatWeb) SendTextMessage(toUserName, content string) (msgID, localID string, err error) {
	msgID, localID, _, err = wxwb.api.SendTextMessage(wxwb.userInfo.user.UserName, toUserName, content)
	if err == nil {
		wxwb.runInfo.MessageCount++
		wxwb.runInfo.MessageSentCount++
	}
	return
}

// SendRevokeMessage 撤回消息
func (wxwb *WechatWeb) SendRevokeMessage(svrMsgID, clientMsgID, toUserName string) (err error) {
	_, err = wxwb.api.SendRevokeMessage(toUserName, svrMsgID, clientMsgID)
	if err == nil {
		wxwb.runInfo.MessageRevokeCount++
		wxwb.runInfo.MessageRevokeSentCount++
	}
	return
}

// ModifyUserRemakName 修改用户备注
func (wxwb *WechatWeb) ModifyUserRemakName(userName, remarkName string) (err error) {
	_, err = wxwb.api.ModifyUserRemakName(userName, remarkName)
	return
}

// ModifyChatRoomTopic 修改群名
func (wxwb *WechatWeb) ModifyChatRoomTopic(userName, newTopic string) (err error) {
	_, err = wxwb.api.ModifyChatRoomTopic(userName, newTopic)
	return
}
