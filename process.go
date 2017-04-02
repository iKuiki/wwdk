package wxweb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/datastruct/appmsg"
	"github.com/yinhui87/wechat-web/tool"
	"html"
	"strconv"
	"strings"
)

func (this *WechatWeb) StatusNotify(fromUserName, toUserName string) (err error) {
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

func (this *WechatWeb) messageProcesser(msg *datastruct.Message) (err error) {
	context := Context{App: this, hasStop: false}
	switch msg.MsgType {
	case datastruct.TEXT_MSG:
		for _, v := range this.messageHook[datastruct.TEXT_MSG] {
			if f, ok := v.(TextMessageHook); ok {
				f(&context, *msg)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.IMAGE_MSG:
		msg.Content = strings.Replace(html.UnescapeString(msg.Content), "<br/>", "", -1)
		var imgContent appmsg.ImageMsgContent
		err = xml.Unmarshal([]byte(msg.Content), &imgContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range this.messageHook[datastruct.IMAGE_MSG] {
			if f, ok := v.(ImageMessageHook); ok {
				f(&context, *msg, imgContent)
			}
			if context.hasStop {
				break
			}
		}
	default:
		return errors.New(fmt.Sprintf("Unknown MsgType %v: %#v", msg.MsgType, msg))
	}
	return nil
}
