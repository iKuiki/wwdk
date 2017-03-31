package wxweb

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func assembleSyncKey(syncKey *datastruct.SyncKey) string {
	keys := make([]string, 0)
	for _, v := range syncKey.List {
		keys = append(keys, strconv.FormatInt(v.Key, 10)+"_"+strconv.FormatInt(v.Val, 10))
	}
	ret := strings.Join(keys, "|")
	return url.QueryEscape(ret)
}

func analysisSyncResp(syncResp string) (result datastruct.SyncCheckRespond) {
	syncResp = strings.TrimPrefix(syncResp, "{")
	syncResp = strings.TrimSuffix(syncResp, "}")
	arr := strings.Split(syncResp, ",")
	result = datastruct.SyncCheckRespond{}
	for _, v := range arr {
		if strings.HasPrefix(v, "retcode") {
			result.Retcode = strings.TrimPrefix(strings.TrimSuffix(v, `"`), `retcode:"`)
		}
		if strings.HasPrefix(v, "selector") {
			result.Selector = strings.TrimPrefix(strings.TrimSuffix(v, `"`), `selector:"`)
		}
	}
	return result
}

func (this *WechatWeb) syncCheck() (selector string, err error) {
	req := httplib.Get("https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck")
	req.Param("r", tool.GetWxTimeStamp())
	req.Param("skey", this.sKey)
	req.Param("sid", this.cookie.Wxsid)
	req.Param("uin", this.cookie.Wxuin)
	req.Param("deviceid", this.deviceId)
	req.Param("synckey", assembleSyncKey(this.syncKey))
	req.Param("_", tool.GetWxTimeStamp())
	setWechatCookie(req, this.cookie)
	resp, err := req.String()
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	retArr := tool.AnalysisWxWindowRespond(resp)

	ret := analysisSyncResp(retArr["window.synccheck"])
	if ret.Retcode != "0" {
		return "", errors.New("respond Retcode " + ret.Retcode)
	}
	return ret.Selector, nil
}

func (this *WechatWeb) getMessage() (gmResp datastruct.GetMessageRespond, err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsync")
	req.Param("sid", this.cookie.Wxsid)
	req.Param("skey", this.cookie.Skey)
	req.Param("pass_ticket", this.cookie.PassTicket)
	setWechatCookie(req, this.cookie)
	gmResp = datastruct.GetMessageRespond{}
	reqBody := datastruct.GetMessageRequest{
		BaseRequest: getBaseRequest(this.cookie, this.deviceId),
		SyncKey:     this.syncKey,
		Rr:          ^time.Now().Unix() + 1,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return gmResp, errors.New("Marshal request body to json fail: " + err.Error())
	}
	req.Body(body)
	resp, err := req.Bytes()
	if err != nil {
		return gmResp, errors.New("request error: " + err.Error())
	}
	err = json.Unmarshal(resp, &gmResp)
	if err != nil {
		return gmResp, errors.New("Unmarshal respond json fail: " + err.Error())
	}
	if gmResp.BaseResponse.Ret != 0 {
		return gmResp, errors.New("respond error ret: " + strconv.FormatInt(gmResp.BaseResponse.Ret, 10))
	}
	// if gmResp.AddMsgCount > 0 {
	// 	fmt.Println(string(resp))
	// 	panic(nil)
	// }
	return gmResp, nil
}

func (this *WechatWeb) SaveMessageImage(msgId string) (filename string, err error) {
	req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetmsgimg")
	req.Param("MsgID", msgId)
	req.Param("skey", this.sKey)
	// req.Param("type", "slave")
	setWechatCookie(req, this.cookie)
	filename = msgId
	err = req.ToFile(filename)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	return filename, nil
}

func (this *WechatWeb) StartServe() {
	for true {
		selector, err := this.syncCheck()
		if err != nil {
			log.Printf("SyncCheck error: %s\n", err.Error())
			continue
		}
		switch selector {
		case "7":
			// log.Println("SyncCheck 7")
			gmResp, err := this.getMessage()
			if err != nil {
				log.Printf("GetMessage error: %s\n", err.Error())
				continue
			}
			this.syncKey = gmResp.SyncKey
			for _, msg := range gmResp.AddMsgList {
				from, err := this.getContact(msg.FromUserName)
				if err != nil {
					log.Printf("getContact error: %s\n", err.Error())
					continue
				}
				err = this.messageProcesser(msg, from)
				if err != nil {
					log.Printf("MessageProcesser error: %s\n", err.Error())
					continue
				}
			}
		default:
			log.Printf("SyncCheck Unknow selector: %s\n", selector)
		}
		time.Sleep(100 * time.Microsecond)
	}
}
