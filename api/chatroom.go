package api

import (
	"bytes"
	"encoding/json"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// ModifyChatRoomTopic 修改聊天室标题
// @param userName 要修改的群的UserName
// @param remarkName 新的群名
func (api *WechatwebAPI) ModifyChatRoomTopic(userName, newTopic string) (body []byte, err error) {
	mctReq := datastruct.ModifyChatRoomTopicRequest{
		BaseRequest:  api.baseRequest(),
		NewTopic:     newTopic,
		ChatRoomName: userName,
	}
	reqBody, err := json.Marshal(mctReq)
	if err != nil {
		return nil, errors.New("Marshal reqBody to json fail: " + err.Error())
	}
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxupdatechatroom?fun=modtopic", bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.New("create request error: " + err.Error())
	}
	resp, err := api.request(req)
	if err != nil {
		return nil, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	var mctResp datastruct.ModifyChatRoomTopicRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read response body error: " + err.Error())
	}
	err = json.Unmarshal(body, &mctResp)
	if err != nil {
		return nil, errors.New("UnMarshal respond json fail: " + err.Error())
	}
	if mctResp.BaseResponse.Ret != 0 {
		return nil, errors.Errorf("Respond error ret(%d): %s", mctResp.BaseResponse.Ret, mctResp.BaseResponse.ErrMsg)
	}
	return
}
