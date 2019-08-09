package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"github.com/pkg/errors"
)

// CreateChatroom 好友拉群创建聊天室
// @param topic 新聊天室的标题
// @param userNames 要添加到新聊天室的用户名列表
func (api *wechatwebAPI) CreateChatroom(topic string, userNames []string) (chatroomUserName string, body []byte, err error) {
	if len(userNames) == 0 {
		err = errors.New("userName list empty")
	}
	var memberList []datastruct.CreateChatroomRequestMemberList
	for _, userName := range userNames {
		memberList = append(memberList, datastruct.CreateChatroomRequestMemberList{
			UserName: userName,
		})
	}
	ccRequest := datastruct.CreateChatroomRequest{
		BaseRequest: api.baseRequest(),
		Topic:       topic,
		MemberCount: int64(len(memberList)),
		MemberList:  memberList,
	}
	reqBody, err := json.Marshal(ccRequest)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxcreatechatroom?"+params.Encode(), bytes.NewReader(reqBody))
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
	var ccResp datastruct.CreateChatroomResponse
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &ccResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if ccResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", ccResp.BaseResponse.Ret, ccResp.BaseResponse.ErrMsg)
		return
	}
	chatroomUserName = ccResp.ChatRoomName
	return
}

// ModifyChatRoomTopic 修改聊天室标题
// @param userName 要修改的群的UserName
// @param remarkName 新的群名
func (api *wechatwebAPI) ModifyChatRoomTopic(userName, newTopic string) (body []byte, err error) {
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
