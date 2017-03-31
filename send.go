package wxweb

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"strconv"
)

func (this *WechatWeb) SendTextMessage(toUserName, content string) (sendMessageRespond *datastruct.SendMessageRespond, err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg")
	req.Param("pass_ticket", this.cookie.PassTicket)
	setWechatCookie(req, this.cookie)
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: getBaseRequest(this.cookie, this.deviceId),
		Msg: &datastruct.SendMessage{

			ClientMsgID:  tool.GetWxTimeStamp(),
			Content:      content,
			FromUserName: this.user.UserName,
			LocalID:      tool.GetWxTimeStamp(),
			ToUserName:   toUserName,
			Type:         datastruct.TEXT_MSG,
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
