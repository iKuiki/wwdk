package wxweb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/datastruct/appmsg"
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
	filename = msgId + ".jpg"
	err = req.ToFile(filename)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	return filename, nil
}

func (this *WechatWeb) SaveMessageVideo(msg datastruct.Message) (filename string, err error) {
	if msg.MsgType != datastruct.LITTLE_VIDEO_MSG {
		return "", errors.New("Message type wrong")
	}
	var videoContent appmsg.VideoMsgContent
	err = xml.Unmarshal([]byte(msg.Content), &videoContent)
	if err != nil {
		return "", errors.New("Unmarshal message content to struct: " + err.Error())
	}
	req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvideo")
	req.Param("msgid", msg.MsgID)
	req.Param("skey", this.sKey)
	setWechatCookie(req, this.cookie)
	req.Header("Range", "bytes=0-")
	filename = msg.MsgID + ".mp4"
	// err = req.ToFile(filename)
	resp, err := req.Response()
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	length, err := strconv.ParseInt(videoContent.VideoMsg.Length, 10, 64)
	if err != nil {
		return "", errors.New("Parse Video Content Length error: " + err.Error())
	}
	if resp.ContentLength != length {
		return "", errors.New("Respond content length wrong")
	}
	n, err := tool.WriteToFile(filename, resp.Body)
	if err != nil {
		return "", errors.New("WriteToFile error: " + err.Error())
	}
	if int64(n) != length {
		return filename, errors.New("File size wrong")
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
				err = this.messageProcesser(msg)
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
