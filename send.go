package wxweb

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"strconv"
)

// SendTextMessage 发送消息
func (wxwb *WechatWeb) SendTextMessage(toUserName, content string) (sendMessageRespond *datastruct.SendMessageRespond, err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg")
	req.Param("pass_ticket", wxwb.cookie.PassTicket)
	setWechatCookie(req, wxwb.cookie)
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: getBaseRequest(wxwb.cookie, wxwb.deviceID),
		Msg: &datastruct.SendMessage{

			ClientMsgID:  tool.GetWxTimeStamp(),
			Content:      content,
			FromUserName: wxwb.user.UserName,
			LocalID:      tool.GetWxTimeStamp(),
			ToUserName:   toUserName,
			Type:         datastruct.TextMsg,
		},
	}
	body, err := json.Marshal(msgReq)
	if err != nil {
		return nil, errors.New("Marshal body to json fail: " + err.Error())
	}
	req.Body(body)
	resp, err := req.Bytes()
	if err != nil {
		return nil, errors.New("request error: " + err.Error())
	}
	var smResp datastruct.SendMessageRespond
	err = json.Unmarshal(resp, &smResp)
	if err != nil {
		return nil, errors.New("UnMarshal respond json fail: " + err.Error())
	}
	if smResp.BaseResponse.Ret != 0 {
		return nil, errors.New("Respond error ret: " + strconv.FormatInt(smResp.BaseResponse.Ret, 10))
	}
	return &smResp, nil
}

// SendRevokeMessage 撤回消息
func (wxwb *WechatWeb) SendRevokeMessage(svrMsgID, clientMsgID, toUserName string) (revokeMessageRespond *datastruct.RevokeMessageRespond, err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxrevokemsg")
	setWechatCookie(req, wxwb.cookie)
	srmReq := datastruct.RevokeMessageRequest{
		BaseRequest: getBaseRequest(wxwb.cookie, wxwb.deviceID),
		ClientMsgID: clientMsgID,
		SvrMsgID:    svrMsgID,
		ToUserName:  toUserName,
	}
	body, err := json.Marshal(srmReq)
	if err != nil {
		return nil, errors.New("Marshal body to json fail: " + err.Error())
	}
	req.Body(body)
	resp, err := req.Bytes()
	if err != nil {
		return nil, errors.New("request error: " + err.Error())
	}
	var rmResp datastruct.RevokeMessageRespond
	err = json.Unmarshal(resp, &rmResp)
	if err != nil {
		return nil, errors.New("UnMarshal respond json fail: " + err.Error())
	}
	if rmResp.BaseResponse.Ret != 0 {
		return nil, errors.New("Respond error ret: " + strconv.FormatInt(rmResp.BaseResponse.Ret, 10))
	}
	return &rmResp, nil
}
