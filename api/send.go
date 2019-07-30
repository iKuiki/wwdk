package api

import (
	"bytes"
	"encoding/json"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// StatusNotify 消息已读通知
// @param fromUserName 自己的UserName
// @param toUserName 已读的联系人的UserName
// @param code 用途未知，目前发现的时登陆完毕调用时此处填3（from和to的username都一样）其余时候都是1
func (api *WechatwebAPI) StatusNotify(fromUserName, toUserName string, code int64) (body []byte, err error) {
	msgID, _ := strconv.ParseInt(tool.GetWxTimeStamp(), 10, 64)
	data := datastruct.StatusNotifyRequest{
		BaseRequest:  api.baseRequest(),
		ClientMsgID:  msgID,
		Code:         code,
		FromUserName: fromUserName,
		ToUserName:   toUserName,
	}
	reqBody, err := json.Marshal(data)
	if err != nil {
		err = errors.New("Marshal request body to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxstatusnotify?"+params.Encode(), bytes.NewReader(reqBody))
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	var snResp datastruct.StatusNotifyRespond
	err = json.Unmarshal(body, &snResp)
	if err != nil {
		err = errors.New("Unmarshal respond json fail: " + err.Error())
		return
	}
	if snResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("respond error ret(%d): %s", snResp.BaseResponse.Ret, snResp.BaseResponse.ErrMsg)
		return
	}
	return
}

// SendTextMessage 发送消息
// @param fromUserName 自己的UserName
// @param toUserName 要发送的目标联系人的UserName
// @param content 文字消息内容
// @return MsgID 消息的服务器ID（发送后由服务器生成）
// @return LocalID 消息本地ID（本地生成的）
func (api *WechatwebAPI) SendTextMessage(fromUserName, toUserName, content string) (MsgID, LocalID string, body []byte, err error) {
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: api.baseRequest(),
		Msg: &datastruct.SendMessage{
			ClientMsgID:  tool.GetWxTimeStamp(),
			Content:      content,
			FromUserName: fromUserName,
			LocalID:      tool.GetWxTimeStamp(),
			ToUserName:   toUserName,
			Type:         datastruct.TextMsg,
		},
	}
	reqBody, err := json.Marshal(msgReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsendmsg?"+params.Encode(), bytes.NewReader(reqBody))
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	var smResp datastruct.SendMessageRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &smResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if smResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", smResp.BaseResponse.Ret, smResp.BaseResponse.ErrMsg)
		return
	}
	MsgID, LocalID = smResp.MsgID, smResp.LocalID
	return
}

// SendRevokeMessage 撤回消息
// @param toUserName 要发送的目标联系人的UserName
// @param svrMsgID 消息的服务器ID（发送后由服务器生成）
// @param clientMsgID 消息本地ID（本地生成的）
func (api *WechatwebAPI) SendRevokeMessage(toUserName, svrMsgID, clientMsgID string) (body []byte, err error) {
	srmReq := datastruct.RevokeMessageRequest{
		BaseRequest: api.baseRequest(),
		ClientMsgID: clientMsgID,
		SvrMsgID:    svrMsgID,
		ToUserName:  toUserName,
	}
	reqBody, err := json.Marshal(srmReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxrevokemsg", bytes.NewReader(reqBody))
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	var rmResp datastruct.RevokeMessageRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &rmResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if rmResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", rmResp.BaseResponse.Ret, rmResp.BaseResponse.ErrMsg)
		return
	}
	return
}
