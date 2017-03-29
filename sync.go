package wxweb

import (
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

func syncCheck(sKey, deviceId string, cookie *wechatCookie, syncKey *datastruct.SyncKey) (selector string, err error) {
	req := httplib.Get("https://webpush2.weixin.qq.com/cgi-bin/mmwebwx-bin/synccheck")
	req.Param("r", tool.GetWxTimeStamp())
	req.Param("skey", sKey)
	req.Param("sid", cookie.Wxsid)
	req.Param("uin", cookie.Wxuin)
	req.Param("deviceid", deviceId)
	req.Param("synckey", assembleSyncKey(syncKey))
	req.Param("_", tool.GetWxTimeStamp())
	setWechatCookie(req, cookie)
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

func (this *WechatWeb) StartServe() {
	for true {
		selector, err := syncCheck(this.sKey, this.deviceId, this.cookie, this.syncKey)
		if err != nil {
			log.Printf("SyncCheck error: %s\n", err.Error())
			continue
		}
		switch selector {
		default:
			log.Printf("SyncCheck Unknow selector: %s\n", selector)
		}
		time.Sleep(100 * time.Microsecond)
	}
}
