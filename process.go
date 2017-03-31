package wxweb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"log"
	"strconv"
)

func (this *WechatWeb) statusNotify(fromUserName, toUserName string) (err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify")
	req.Param("pass_ticket", this.cookie.PassTicket)
	setWechatCookie(req, this.cookie)
	msgId, _ := strconv.ParseInt(tool.GetWxTimeStamp(), 10, 64)
	reqBody := datastruct.StatusNotifyRequest{
		BaseRequest:  getBaseRequest(this.cookie, this.deviceId),
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

func (this *WechatWeb) messageProcesser(msg *datastruct.Message, from *datastruct.Contact) (err error) {
	switch msg.MsgType {
	case datastruct.TEXT_MSG:
		log.Printf("Recived a text msg from %s: %s", from.NickName, msg.Content)
		// reply the same message
		smResp, err := this.SendTextMessage(msg.FromUserName, msg.Content)
		if err != nil {
			return errors.New("sendTextMessage error: " + err.Error())
		}
		log.Println("messageSent, msgId: " + smResp.MsgID + ", Local ID: " + smResp.LocalID)
		// Set message to readed at phone
		err = this.statusNotify(msg.ToUserName, msg.FromUserName)
		if err != nil {
			return errors.New("StatusNotify error: " + err.Error())
		}
	default:
		return errors.New(fmt.Sprintf("Unknown MsgType %v: %#v", msg.MsgType, msg))
	}
	return nil
}
