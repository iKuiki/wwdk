package wxweb

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"strconv"
)

func sendTextMessage(cookie *wechatCookie, deviceId string, userUserName, toUserName, content string) (err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg")
	req.Param("pass_ticket", cookie.PassTicket)
	setWechatCookie(req, cookie)
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: getBaseRequest(cookie, deviceId),
		Msg: &datastruct.SendMessage{

			ClientMsgID:  tool.GetWxTimeStamp(),
			Content:      content,
			FromUserName: userUserName,
			LocalID:      tool.GetWxTimeStamp(),
			ToUserName:   toUserName,
			Type:         datastruct.TEXT_MSG,
		},
	}
	body, err := json.Marshal(msgReq)
	if err != nil {
		return errors.New("Marshal body to json fail: " + err.Error())
	}
	req.Body(body)
	resp, err := req.Bytes()
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	var smResp datastruct.SendMessageRespond
	err = json.Unmarshal(resp, &smResp)
	if err != nil {
		return errors.New("UnMarshal respond json fail: " + err.Error())
	}
	if smResp.BaseResponse.Ret != 0 {
		return errors.New("Respond error ret: " + strconv.FormatInt(smResp.BaseResponse.Ret, 10))
	}
	return nil
}

func (this *WechatWeb) SendMessage(toUserName, content string) (err error) {
	return sendTextMessage(this.cookie, this.deviceId, this.user.UserName, toUserName, content)
}
