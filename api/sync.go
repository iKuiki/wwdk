package api

import (
	"bytes"
	"encoding/json"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	// ErrLogout 错误：已经登出
	ErrLogout = errors.New("Logout")
)

// SyncCheck 检查同步
// 轮询微信服务器，如果有新的状态，会通过此接口返回需要同步的信息
// @return retCode 状态码，正常为0
// @return selector 状态机，根据此项判定是否需要拉去新信息
func (api *WechatwebAPI) SyncCheck() (retCode, selector string, body []byte, err error) {
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	params.Set("sid", api.loginInfo.Wxsid)
	params.Set("uin", api.loginInfo.Wxuin)
	params.Set("deviceid", api.deviceID)
	params.Set("synckey", tool.AssembleSyncKey(api.loginInfo.SyncKey))
	params.Set("_", tool.GetWxTimeStamp())
	req, err := http.NewRequest("GET", "https://webpush."+api.apiDomain+"/cgi-bin/mmwebwx-bin/synccheck?"+params.Encode(), nil)
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	retArr := tool.ExtractWxWindowRespond(string(body))

	ret := tool.AnalysisSyncResp(retArr["window.synccheck"])
	retCode, selector = ret.Retcode, ret.Selector
	if retCode != "0" {
		if retCode == "1101" {
			err = ErrLogout
			return
		}
		err = errors.New("respond Retcode " + retCode)
		return
	}
	return
}

// WebwxSync 同步消息
// 如果检查同步接口返回有新消息需要同步，通过此接口从服务器中获取新消息
// @return syncResp 同步状态返回值
func (api *WechatwebAPI) WebwxSync() (modContacts []datastruct.Contact,
	delContacts []datastruct.WebwxSyncRespondDelContactListItem,
	addMessages []datastruct.Message,
	body []byte, err error) {
	syncResp := datastruct.WebwxSyncRespond{}
	reqBody, err := json.Marshal(datastruct.WebwxSyncRequest{
		BaseRequest: api.baseRequest(),
		SyncKey:     api.loginInfo.SyncKey,
		Rr:          ^time.Now().Unix() + 1,
	})
	if err != nil {
		err = errors.New("Marshal request body to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("sid", api.loginInfo.Wxsid)
	params.Set("skey", api.loginInfo.SKey)
	// params.Set("pass_ticket", api.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsync?"+params.Encode(), bytes.NewReader(reqBody))
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	// resp, err := api.client.Post("https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsync?"+params.Encode(),
	// 	"application/json;charset=UTF-8",
	// 	bytes.NewReader(reqBody))
	// if err != nil {
	// 	return syncResp, errors.New("request error: " + err.Error())
	// }
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &syncResp)
	if err != nil {
		err = errors.New("Unmarshal respond json fail: " + err.Error())
		return
	}
	// 更新SyncKey
	if syncResp.SyncCheckKey != nil {
		api.loginInfo.SyncKey = syncResp.SyncCheckKey
	} else {
		api.loginInfo.SyncKey = syncResp.SyncKey
	}
	if syncResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("respond error ret(%d): %s", syncResp.BaseResponse.Ret, syncResp.BaseResponse.ErrMsg)
		return
	}
	// 赋值结果
	modContacts, delContacts, addMessages = syncResp.ModContactList, syncResp.DelContactList, syncResp.AddMsgList
	return
}
