package wwdk

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
)

var syncHosts = []string{
	"wx2.qq.com",
	"webpush.wx.qq.com",
	"webpush.wx2.qq.com",
	"wx8.qq.com",
	"webpush.wx8.qq.com",
	"qq.com",
	"webpush.wx.qq.com",
	"web2.wechat.com",
	"webpush.web2.wechat.com",
	"wechat.com",
	"webpush.web.wechat.com",
	"webpush.weixin.qq.com",
	"webpush.wechat.com",
	"webpush1.wechat.com",
	"webpush2.wechat.com",
	"webpush2.wx.qq.com",
}

// assembleSyncKey 组装synckey
// 将同步需要的synckey组装为请求字符串
func assembleSyncKey(syncKey *datastruct.SyncKey) string {
	keys := make([]string, 0)
	for _, v := range syncKey.List {
		keys = append(keys, strconv.FormatInt(v.Key, 10)+"_"+strconv.FormatInt(v.Val, 10))
	}
	ret := strings.Join(keys, "|")
	// return url.QueryEscape(ret)
	return ret
}

// analysisSyncResp 解析同步状态返回值
// 同步状态返回的接口
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

func (wxwb *WechatWeb) chooseSyncHost() bool {
	wxwb.logger.Info("choose sync host...\n")
	for _, host := range syncHosts {
		wxwb.apiRuntime.syncHost = host
		code, _, _ := wxwb.syncCheck()
		if code == `0` {
			wxwb.logger.Infof("sync host [%s] avaliable\n", host)
			return true
		}
	}
	return false
}

// syncCheck 同步状态
// 轮询微信服务器，如果有新的状态，会通过此接口返回需要同步的信息
func (wxwb *WechatWeb) syncCheck() (retCode, selector string, err error) {
	if wxwb.apiRuntime.syncHost == "" {
		return "", "", errors.New("sync host empty")
	}
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	params.Set("sid", wxwb.loginInfo.cookie.Wxsid)
	params.Set("uin", wxwb.loginInfo.cookie.Wxuin)
	params.Set("deviceid", wxwb.apiRuntime.deviceID)
	params.Set("synckey", assembleSyncKey(wxwb.loginInfo.syncKey))
	params.Set("_", tool.GetWxTimeStamp())
	req, err := http.NewRequest("GET", "https://"+wxwb.apiRuntime.syncHost+"/cgi-bin/mmwebwx-bin/synccheck?"+params.Encode(), nil)
	if err != nil {
		return "", "", errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.request(req)
	if err != nil {
		return "", "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	retArr := tool.ExtractWxWindowRespond(string(body))

	ret := analysisSyncResp(retArr["window.synccheck"])
	return ret.Retcode, ret.Selector, nil

	// if ret.Retcode != "0" {
	// 	if ret.Retcode == "1101" {
	// 		return "Logout", nil
	// 	}
	// 	return "", errors.New("respond Retcode " + ret.Retcode)
	// }
	// return ret.Selector, nil
}

// StartServe 启动消息同步服务
func (wxwb *WechatWeb) StartServe() {
	avaliable := wxwb.chooseSyncHost()
	if !avaliable {
		wxwb.logger.Info("all sync host unavaliable, exit...\n")
		return
	}
	getMessage := func() {
		gmResp, err := wxwb.getMessage()
		if err != nil {
			wxwb.logger.Infof("GetMessage error: %s\n", err.Error())
			return
		}
		if gmResp.SyncCheckKey != nil {
			wxwb.loginInfo.syncKey = gmResp.SyncCheckKey
		} else {
			wxwb.loginInfo.syncKey = gmResp.SyncKey
		}
		// 处理新增联系人
		for _, contact := range gmResp.ModContactList {
			wxwb.runInfo.ContactModifyCount++
			wxwb.logger.Infof("Modify contact: %s\n", contact.NickName)
			oldContact := wxwb.contactList[contact.UserName]
			wxwb.contactProcesser(&oldContact, &contact)
			wxwb.contactList[contact.UserName] = contact
		}
		// 新消息
		for _, msg := range gmResp.AddMsgList {
			if msg.MsgType == datastruct.RevokeMsg {
				wxwb.runInfo.MessageRevokeCount++
			} else {
				wxwb.runInfo.MessageCount++
			}
			err = wxwb.messageProcesser(&msg)
			if err != nil {
				wxwb.logger.Infof("MessageProcesser error: %s\n", err.Error())
				continue
			}
		}
	}
	for {
		isBreaked := func() (isBreaked bool) {
			defer func() {
				if r := recover(); r != nil {
					wxwb.logger.Infof("Recovered in StartServe loop: %v\n", r)
					wxwb.runInfo.PanicCount++
				}
			}()
			code, selector, err := wxwb.syncCheck()
			if err != nil {
				wxwb.logger.Infof("SyncCheck error: %s\n", err.Error())
				return false
			}
			if code != "0" {
				switch code {
				case "1101":
					wxwb.logger.Info("User has logout web wechat, exit...\n")
					return true
				case "1100":
					wxwb.logger.Info("sync host unavaliable, choose a new one...\n")
					avaliable = wxwb.chooseSyncHost()
					if !avaliable {
						wxwb.logger.Info("all sync host unavaliable, exit...\n")
						return true
					}
					return false
				}
			}
			// wxwb.logger.Infof("selector: %v\n", selector)
			switch selector {
			case "0":
				// wxwb.logger.Info("SyncCheck 0\n")
				// normal
				// wxwb.logger.Info("no new message\n")
			case "6":
				wxwb.logger.Info("selector is 6\n")
				getMessage()
			case "7":
				wxwb.logger.Info("selector is 7\n")
				getMessage()
			case "1":
				wxwb.logger.Info("selector is 1\n")
				getMessage()
			case "3":
				wxwb.logger.Info("selector is 3\n")
				getMessage()
			case "4":
				wxwb.logger.Info("selector is 4\n")
				getMessage()
			case "5":
				wxwb.logger.Info("selector is 5\n")
				getMessage()
			case "2":
				// wxwb.logger.Info("SyncCheck 2\n")
				getMessage()
			default:
				wxwb.logger.Infof("SyncCheck Unknow selector: %s\n", selector)
			}
			wxwb.runInfo.SyncCount++
			time.Sleep(1000 * time.Millisecond)
			return false
		}()
		if isBreaked {
			break
		}
	}
}
