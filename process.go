package wxweb

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"log"
	"strconv"
)

func statusNotify(cookie *wechatCookie, deviceId string, fromUserName, toUserName string) (err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify")
	req.Param("pass_ticket", cookie.PassTicket)
	setWechatCookie(req, cookie)
	msgId, _ := strconv.ParseInt(tool.GetWxTimeStamp(), 10, 64)
	reqBody := datastruct.StatusNotifyRequest{
		BaseRequest: &datastruct.BaseRequest{
			DeviceID: deviceId,
			Sid:      cookie.Wxsid,
			Skey:     cookie.Skey,
			Uin:      cookie.Wxuin,
		},
		ClientMsgID:  msgId,
		Code:         1,
		FromUserName: fromUserName,
		ToUserName:   toUserName,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return errors.New("Marshal request body to json fail: " + err.Error())
	}
	req.Body(body)
	resp, err := req.Bytes()
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	var snResp datastruct.StatusNotifyRespond
	err = json.Unmarshal(resp, &snResp)
	if err != nil {
		return errors.New("Unmarshal respond json fail: " + err.Error())
	}
	if snResp.BaseResponse.Ret != 0 {
		return errors.New("respond error ret: " + strconv.FormatInt(snResp.BaseResponse.Ret, 10))
	}
	return nil
}

func messageProcesser(cookie *wechatCookie, deviceId string, msg *datastruct.Message, from *datastruct.Contact) (err error) {
	switch msg.MsgType {
	case 1:
		log.Printf("Recived a msg from %s: %s", from.NickName, msg.Content)
		// Set message to readed at phone
		err := statusNotify(cookie, deviceId, msg.ToUserName, msg.FromUserName)
		if err != nil {
			return errors.New("StatusNotify error: " + err.Error())
		}
	default:
		return errors.New("Unknown MsgType: " + strconv.FormatInt(msg.MsgType, 10))
	}
	return nil
}
