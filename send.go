package wwdk

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
)

// StatusNotify 消息已读通知
// TODO: delete
func (wxwb *WechatWeb) StatusNotify(fromUserName, toUserName string, code int64) (err error) {
	msgID, _ := strconv.ParseInt(tool.GetWxTimeStamp(), 10, 64)
	reqBody := datastruct.StatusNotifyRequest{
		BaseRequest:  wxwb.baseRequest(),
		ClientMsgID:  msgID,
		Code:         code,
		FromUserName: fromUserName,
		ToUserName:   toUserName,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return errors.New("Marshal request body to json fail: " + err.Error())
	}
	params := url.Values{}
	params.Set("pass_ticket", wxwb.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+wxwb.apiRuntime.apiDomain+"/cgi-bin/mmwebwx-bin/webwxstatusnotify?"+params.Encode(), bytes.NewReader(data))
	if err != nil {
		return errors.New("create request error: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := wxwb.request(req)
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var snResp datastruct.StatusNotifyRespond
	err = json.Unmarshal(body, &snResp)
	if err != nil {
		return errors.New("Unmarshal respond json fail: " + err.Error())
	}
	if snResp.BaseResponse.Ret != 0 {
		return errors.Errorf("respond error ret(%d): %s", snResp.BaseResponse.Ret, snResp.BaseResponse.ErrMsg)
	}
	return nil
}

// SendTextMessage 发送消息
// TODO: delete
func (wxwb *WechatWeb) SendTextMessage(toUserName, content string) (sendMessageRespond *datastruct.SendMessageRespond, err error) {
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: wxwb.baseRequest(),
		Msg: &datastruct.SendMessage{

			ClientMsgID:  tool.GetWxTimeStamp(),
			Content:      content,
			FromUserName: wxwb.userInfo.user.UserName,
			LocalID:      tool.GetWxTimeStamp(),
			ToUserName:   toUserName,
			Type:         datastruct.TextMsg,
		},
	}
	body, err := json.Marshal(msgReq)
	if err != nil {
		return nil, errors.New("Marshal body to json fail: " + err.Error())
	}
	params := url.Values{}
	params.Set("pass_ticket", wxwb.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+wxwb.apiRuntime.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsendmsg?"+params.Encode(), bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.request(req)
	if err != nil {
		return nil, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	var smResp datastruct.SendMessageRespond
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read response body error: " + err.Error())
	}
	err = json.Unmarshal(respBody, &smResp)
	if err != nil {
		return nil, errors.New("UnMarshal respond json fail: " + err.Error())
	}
	if smResp.BaseResponse.Ret != 0 {
		return nil, errors.Errorf("Respond error ret(%d): %s", smResp.BaseResponse.Ret, smResp.BaseResponse.ErrMsg)
	}
	wxwb.runInfo.MessageCount++
	wxwb.runInfo.MessageSentCount++
	return &smResp, nil
}

// SendRevokeMessage 撤回消息
// TODO: delete
func (wxwb *WechatWeb) SendRevokeMessage(svrMsgID, clientMsgID, toUserName string) (revokeMessageRespond *datastruct.RevokeMessageRespond, err error) {
	srmReq := datastruct.RevokeMessageRequest{
		BaseRequest: wxwb.baseRequest(),
		ClientMsgID: clientMsgID,
		SvrMsgID:    svrMsgID,
		ToUserName:  toUserName,
	}
	body, err := json.Marshal(srmReq)
	if err != nil {
		return nil, errors.New("Marshal body to json fail: " + err.Error())
	}
	req, err := http.NewRequest("POST", "https://"+wxwb.apiRuntime.apiDomain+"/cgi-bin/mmwebwx-bin/webwxrevokemsg", bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.request(req)
	if err != nil {
		return nil, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	var rmResp datastruct.RevokeMessageRespond
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read response body error: " + err.Error())
	}
	err = json.Unmarshal(respBody, &rmResp)
	if err != nil {
		return nil, errors.New("UnMarshal respond json fail: " + err.Error())
	}
	if rmResp.BaseResponse.Ret != 0 {
		return nil, errors.Errorf("Respond error ret(%d): %s", rmResp.BaseResponse.Ret, rmResp.BaseResponse.ErrMsg)
	}
	wxwb.runInfo.MessageRevokeCount++
	wxwb.runInfo.MessageRevokeSentCount++
	return &rmResp, nil
}

// ModifyUserRemakName 修改用户备注
// TODO: delete
func (wxwb *WechatWeb) ModifyUserRemakName(userName, remarkName string) (revokeMessageRespond *datastruct.ModifyRemarkRespond, err error) {
	murReq := datastruct.ModifyRemarkRequest{
		BaseRequest: wxwb.baseRequest(),
		CmdID:       2,
		RemarkName:  remarkName,
		UserName:    userName,
	}
	body, err := json.Marshal(murReq)
	if err != nil {
		return nil, errors.New("Marshal body to json fail: " + err.Error())
	}
	req, err := http.NewRequest("POST", "https://"+wxwb.apiRuntime.apiDomain+"/cgi-bin/mmwebwx-bin/webwxoplog", bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.request(req)
	if err != nil {
		return nil, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	var murResp datastruct.ModifyRemarkRespond
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read response body error: " + err.Error())
	}
	err = json.Unmarshal(respBody, &murResp)
	if err != nil {
		return nil, errors.New("UnMarshal respond json fail: " + err.Error())
	}
	if murResp.BaseResponse.Ret != 0 {
		return nil, errors.Errorf("Respond error ret(%d): %s", murResp.BaseResponse.Ret, murResp.BaseResponse.ErrMsg)
	}
	return &murResp, nil
}

// ModifyChatRoomTopic 修改群名
// TODO: delete
func (wxwb *WechatWeb) ModifyChatRoomTopic(userName, newTopic string) (revokeMessageRespond *datastruct.ModifyChatRoomTopicRespond, err error) {
	mctReq := datastruct.ModifyChatRoomTopicRequest{
		BaseRequest:  wxwb.baseRequest(),
		NewTopic:     newTopic,
		ChatRoomName: userName,
	}
	body, err := json.Marshal(mctReq)
	if err != nil {
		return nil, errors.New("Marshal body to json fail: " + err.Error())
	}
	req, err := http.NewRequest("POST", "https://"+wxwb.apiRuntime.apiDomain+"/cgi-bin/mmwebwx-bin/webwxupdatechatroom?fun=modtopic", bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.request(req)
	if err != nil {
		return nil, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	var mctResp datastruct.ModifyChatRoomTopicRespond
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read response body error: " + err.Error())
	}
	err = json.Unmarshal(respBody, &mctResp)
	if err != nil {
		return nil, errors.New("UnMarshal respond json fail: " + err.Error())
	}
	if mctResp.BaseResponse.Ret != 0 {
		return nil, errors.Errorf("Respond error ret(%d): %s", mctResp.BaseResponse.Ret, mctResp.BaseResponse.ErrMsg)
	}
	return &mctResp, nil
}
