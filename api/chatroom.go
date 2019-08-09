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

// CreateChatRoom 好友拉群创建聊天室
// @param topic 新聊天室的标题
// @param userNames 要添加到新聊天室的用户名列表
func (api *wechatwebAPI) CreateChatRoom(topic string, userNames []string) (chatroomUserName string, body []byte, err error) {
	if len(userNames) == 0 {
		err = errors.New("userName list empty")
	}
	var memberList []datastruct.MemberListItem
	for _, userName := range userNames {
		memberList = append(memberList, datastruct.MemberListItem{
			UserName: userName,
		})
	}
	ccRequest := datastruct.CreateChatRoomRequest{
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
	var ccResp datastruct.CreateChatRoomResponse
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

// UpdateChatRoomTopic 修改聊天室标题
// @param userName 要修改的群的UserName
// @param remarkName 新的群名
func (api *wechatwebAPI) UpdateChatRoomTopic(userName, newTopic string) (body []byte, err error) {
	uctReq := datastruct.ModifyChatRoomTopicRequest{
		BaseRequest:  api.baseRequest(),
		NewTopic:     newTopic,
		ChatRoomName: userName,
	}
	reqBody, err := json.Marshal(uctReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("fun", "modtopic")
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxupdatechatroom?"+params.Encode(), bytes.NewReader(reqBody))
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
	var uctResp datastruct.UpdateChatRoomResponse
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &uctResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if uctResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", uctResp.BaseResponse.Ret, uctResp.BaseResponse.ErrMsg)
		return
	}
	return
}

// UpdateChatRoomAddMember 更新聊天室：添加成员
// @param chatroomUserName 要更新的聊天室的UserName
// @param memberUserName 要添加的聊天室成员
func (api *wechatwebAPI) UpdateChatRoomAddMember(chatroomUserName, memberUserName string) (body []byte, err error) {
	ucaReq := datastruct.UpdateChatRoomAddMemberRequest{
		BaseRequest:   api.baseRequest(),
		ChatRoomName:  chatroomUserName,
		AddMemberList: memberUserName,
	}
	reqBody, err := json.Marshal(ucaReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("fun", "addmember")
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxupdatechatroom?"+params.Encode(), bytes.NewReader(reqBody))
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
	var ucaResp datastruct.UpdateChatRoomResponse
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &ucaResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if ucaResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", ucaResp.BaseResponse.Ret, ucaResp.BaseResponse.ErrMsg)
		return
	}
	return
}

// UpdateChatRoomDelMember 更新聊天室：移除成员
// @param chatroomUserName 要更新的聊天室的UserName
// @param memberUserName 要移除的聊天室成员
func (api *wechatwebAPI) UpdateChatRoomDelMember(chatroomUserName, memberUserName string) (body []byte, err error) {
	ucdReq := datastruct.UpdateChatRoomDelMemberRequest{
		BaseRequest:   api.baseRequest(),
		ChatRoomName:  chatroomUserName,
		DelMemberList: memberUserName,
	}
	reqBody, err := json.Marshal(ucdReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("fun", "delmember")
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxupdatechatroom?"+params.Encode(), bytes.NewReader(reqBody))
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
	var ucdResp datastruct.UpdateChatRoomResponse
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &ucdResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if ucdResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", ucdResp.BaseResponse.Ret, ucdResp.BaseResponse.ErrMsg)
		return
	}
	return
}
