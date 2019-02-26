package wxweb

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ikuiki/wechat-web/datastruct"
	"github.com/ikuiki/wechat-web/tool"
)

var syncHosts = []string{
	"webpush.wx.qq.com",
	"wx2.qq.com",
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
	log.Println("choose sync host...")
	for _, host := range syncHosts {
		wxwb.syncHost = host
		code, _, _ := wxwb.syncCheck()
		if code == `0` {
			log.Printf("sync host [%s] avaliable", host)
			return true
		}
	}
	return false
}

// syncCheck 同步状态
// 轮询微信服务器，如果有新的状态，会通过此接口返回需要同步的信息
func (wxwb *WechatWeb) syncCheck() (retCode, selector string, err error) {
	if wxwb.syncHost == "" {
		return "", "", errors.New("sync host empty")
	}
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	params.Set("sid", wxwb.loginInfo.cookie.Wxsid)
	params.Set("uin", wxwb.loginInfo.cookie.Wxuin)
	params.Set("deviceid", wxwb.deviceID)
	params.Set("synckey", assembleSyncKey(wxwb.loginInfo.syncKey))
	params.Set("_", tool.GetWxTimeStamp())
	req, err := http.NewRequest("GET", "https://"+wxwb.syncHost+"/cgi-bin/mmwebwx-bin/synccheck?"+params.Encode(), nil)
	if err != nil {
		return "", "", errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.request(req)
	if err != nil {
		return "", "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	retArr := tool.AnalysisWxWindowRespond(string(body))

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
		log.Println("all sync host unavaliable, exit...")
		return
	}
	getMessage := func() {
		gmResp, err := wxwb.getMessage()
		if err != nil {
			log.Printf("GetMessage error: %s\n", err.Error())
			return
		}
		if gmResp.SyncCheckKey != nil {
			wxwb.loginInfo.syncKey = gmResp.SyncCheckKey
		} else {
			wxwb.loginInfo.syncKey = gmResp.SyncKey
		}
		// 处理新增联系人
		for _, contact := range gmResp.ModContactList {
			log.Println("Modify contact: ", contact.NickName)
			oldContact := wxwb.contactList[contact.UserName]
			wxwb.contactProcesser(&oldContact, &contact)
			wxwb.contactList[contact.UserName] = contact
		}
		// 新消息
		for _, msg := range gmResp.AddMsgList {
			err = wxwb.messageProcesser(&msg)
			if err != nil {
				log.Printf("MessageProcesser error: %s\n", err.Error())
				continue
			}
		}
	}
	for {
		isBreaked := func() (isBreaked bool) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in StartServe loop: ", r)
				}
			}()
			code, selector, err := wxwb.syncCheck()
			if err != nil {
				log.Printf("SyncCheck error: %s\n", err.Error())
				return false
			}
			if code != "0" {
				switch code {
				case "1101":
					log.Println("User has logout web wechat, exit...")
					return true
				case "1100":
					log.Println("sync host unavaliable, choose a new one...")
					avaliable = wxwb.chooseSyncHost()
					if !avaliable {
						log.Println("all sync host unavaliable, exit...")
						return true
					}
					return false
				}
			}
			// log.Println("selector: ", selector)
			switch selector {
			case "0":
				// log.Println("SyncCheck 0")
				// normal
				// log.Println("no new message")
			case "6":
				log.Printf("selector is 6")
				getMessage()
			case "7":
				log.Printf("selector is 7")
				getMessage()
			case "1":
				log.Printf("selector is 1")
				getMessage()
			case "3":
				log.Printf("selector is 3")
				getMessage()
			case "4":
				log.Printf("selector is 4")
				getMessage()
			case "5":
				log.Printf("selector is 5")
				getMessage()
			case "2":
				// log.Println("SyncCheck 2")
				getMessage()
			default:
				log.Printf("SyncCheck Unknow selector: %s\n", selector)
			}
			time.Sleep(1000 * time.Millisecond)
			return false
		}()
		if isBreaked {
			break
		}
	}
}
