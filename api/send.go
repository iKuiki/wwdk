package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/datastruct/msgcontent"
	"github.com/ikuiki/wwdk/tool"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// StatusNotify 消息已读通知
// @param fromUserName 自己的UserName
// @param toUserName 已读的联系人的UserName
// @param code 用途未知，目前发现的时登陆完毕调用时此处填3（from和to的username都一样）其余时候都是1
func (api *wechatwebAPI) StatusNotify(fromUserName, toUserName string, code int64) (body []byte, err error) {
	msgID, _ := strconv.ParseInt(tool.GetWxTimeStamp(), 10, 64)
	data := datastruct.StatusNotifyRequest{
		BaseRequest:  api.baseRequest(),
		ClientMsgID:  msgID,
		Code:         code,
		FromUserName: fromUserName,
		ToUserName:   toUserName,
	}
	reqBody, err := json.Marshal(data)
	if err != nil {
		err = errors.New("Marshal request body to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxstatusnotify?"+params.Encode(), bytes.NewReader(reqBody))
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	var snResp datastruct.StatusNotifyRespond
	err = json.Unmarshal(body, &snResp)
	if err != nil {
		err = errors.New("Unmarshal respond json fail: " + err.Error())
		return
	}
	if snResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("respond error ret(%d): %s", snResp.BaseResponse.Ret, snResp.BaseResponse.ErrMsg)
		return
	}
	return
}

// SendTextMessage 发送消息
// @param fromUserName 自己的UserName
// @param toUserName 要发送的目标联系人的UserName
// @param content 文字消息内容
// @return MsgID 消息的服务器ID（发送后由服务器生成）
// @return LocalID 消息本地ID（本地生成的）
func (api *wechatwebAPI) SendTextMessage(fromUserName, toUserName, content string) (MsgID, LocalID string, body []byte, err error) {
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: api.baseRequest(),
		Msg: &datastruct.SendMessage{
			ClientMsgID:  tool.GetWxTimeStamp(),
			Content:      content,
			FromUserName: fromUserName,
			LocalID:      tool.GetWxTimeStamp(),
			ToUserName:   toUserName,
			Type:         datastruct.TextMsg,
		},
	}
	reqBody, err := json.Marshal(msgReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsendmsg?"+params.Encode(), bytes.NewReader(reqBody))
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
	var smResp datastruct.SendMessageRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &smResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if smResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", smResp.BaseResponse.Ret, smResp.BaseResponse.ErrMsg)
		return
	}
	MsgID, LocalID = smResp.MsgID, smResp.LocalID
	return
}

// SendRevokeMessage 撤回消息
// @param toUserName 要发送的目标联系人的UserName
// @param svrMsgID 消息的服务器ID（发送后由服务器生成）
// @param clientMsgID 消息本地ID（本地生成的）
func (api *wechatwebAPI) SendRevokeMessage(toUserName, svrMsgID, clientMsgID string) (body []byte, err error) {
	srmReq := datastruct.RevokeMessageRequest{
		BaseRequest: api.baseRequest(),
		ClientMsgID: clientMsgID,
		SvrMsgID:    svrMsgID,
		ToUserName:  toUserName,
	}
	reqBody, err := json.Marshal(srmReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxrevokemsg", bytes.NewReader(reqBody))
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
	var rmResp datastruct.RevokeMessageRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &rmResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if rmResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", rmResp.BaseResponse.Ret, rmResp.BaseResponse.ErrMsg)
		return
	}
	return
}

// UploadMedia 上传媒体文件
// 在发送图片、视频、动图、文件消息前需要先将文件上传至微信服务器换取MediaID
// @param fromUserName 自己的userName
// @param toUserName 目标的userName
// @param fileName 文件的名字，其后缀很重要
// @param mediaData 要上传的文件的二进制数据
// @return mediaID 上传文件成功后返回的mediaID
func (api *wechatwebAPI) UploadMedia(fromUserName, toUserName, fileName string, mediaData []byte) (mediaID string, body []byte, err error) {
	var bodyBuf bytes.Buffer
	// 识别fileName的扩展名
	extName := path.Ext(fileName)
	if extName == "" {
		err = errors.Errorf("Unknown ext name of file %s", fileName)
	}
	mimeType := mime.TypeByExtension(extName)
	if mimeType == "" {
		err = errors.Errorf("Unknown mime type of ext %s", extName)
	}
	// 构造multipart
	multiWriter := multipart.NewWriter(&bodyBuf)
	// 设置multipart字段id
	multiWriter.WriteField("id",
		"WU_FILE_"+strconv.FormatInt(
			atomic.AddInt64(
				&api.uploadMediaIncrID, 1), 10))
	// 设置multipart字段name
	multiWriter.WriteField("name", fileName)
	// 设置multipart字段type
	multiWriter.WriteField("type",
		mime.TypeByExtension(extName))
	// 设置multipart字段lastModifiedDate
	multiWriter.WriteField("lastModifiedData",
		time.Now().Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"))
	// 设置multipart字段size
	fileSize := int64(len(mediaData))
	multiWriter.WriteField("size",
		strconv.FormatInt(
			fileSize,
			10,
		))
	// 设置multipart字段mediatype
	mediaType := "doc"
	switch strings.ToLower(extName) {
	case ".png":
		mediaType = "pic"
	case ".jpg":
		mediaType = "pic"
	case "mp4":
		mediaType = "video"
	}
	multiWriter.WriteField("mediatype", mediaType)
	// 设置multipart字段uploadmediarequest
	wxTimeStamp := tool.GetWxTimeStamp()
	clientMediaID, err := strconv.ParseInt(wxTimeStamp, 10, 64)
	if err != nil {
		err = errors.Errorf("parse wxTimeStamp %s fail: %v", wxTimeStamp, err)
		return
	}
	// 计算md5
	md5hash := md5.New()
	md5hash.Write(mediaData)
	md5sum := md5hash.Sum(nil)
	umReq := datastruct.UploadMediaRequest{
		BaseRequest:   api.baseRequest(),
		UploadType:    2, // 固定填2
		ClientMediaID: clientMediaID,
		TotalLen:      fileSize,
		StartPos:      0,
		DataLen:       fileSize,
		MediaType:     4, // 固定填4
		FromUserName:  fromUserName,
		ToUserName:    toUserName,
		FileMd5:       hex.EncodeToString(md5sum),
	}
	reqBody, err := json.Marshal(umReq)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	multiWriter.WriteField("uploadmediarequest", string(reqBody))
	// 设置multipart字段webwx_data_ticket
	multiWriter.WriteField("webwx_data_ticket",
		api.loginInfo.DataTicket)
	// 设置multipart字段pass_ticket
	multiWriter.WriteField("pass_ticket",
		api.loginInfo.PassTicket)
	// 设置文件字段
	fieldWrite, err := multiWriter.CreateFormFile("filename", fileName)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	_, err = io.Copy(fieldWrite, bytes.NewBuffer(mediaData))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	multiWriter.Close()
	// 构造请求
	req, err := http.NewRequest("POST", "https://file."+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json", &bodyBuf)
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	// 设置Content-Type
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())
	// 执行请求
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	var umResp datastruct.UploadMediaResponse
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &umResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if umResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", umResp.BaseResponse.Ret, umResp.BaseResponse.ErrMsg)
		return
	}
	mediaID = umResp.MediaID
	return
}

// SendImageMessage 发送图片消息
// @param fromUserName 自己的UserName
// @param toUserName 要发送的目标联系人的UserName
// @param mediaID 上传后的mediaID内容
// @return MsgID 消息的服务器ID（发送后由服务器生成）
// @return LocalID 消息本地ID（本地生成的）
func (api *wechatwebAPI) SendImageMessage(fromUserName, toUserName, mediaID string) (MsgID, LocalID string, body []byte, err error) {
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: api.baseRequest(),
		Msg: &datastruct.SendMessage{
			Type:         datastruct.ImageMsg,
			MediaID:      mediaID,
			FromUserName: fromUserName,
			ToUserName:   toUserName,
			LocalID:      tool.GetWxTimeStamp(),
			ClientMsgID:  tool.GetWxTimeStamp(),
		},
	}
	reqBody, err := json.Marshal(msgReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("fun", "async")
	params.Set("f", "json")
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsendmsgimg?"+params.Encode(), bytes.NewReader(reqBody))
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
	var smResp datastruct.SendMessageRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &smResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if smResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", smResp.BaseResponse.Ret, smResp.BaseResponse.ErrMsg)
		return
	}
	MsgID, LocalID = smResp.MsgID, smResp.LocalID
	return
}

// SendVideoMessage 发送视频消息
// @param fromUserName 自己的UserName
// @param toUserName 要发送的目标联系人的UserName
// @param mediaID 上传后的mediaID内容
// @return MsgID 消息的服务器ID（发送后由服务器生成）
// @return LocalID 消息本地ID（本地生成的）
func (api *wechatwebAPI) SendVideoMessage(fromUserName, toUserName, mediaID string) (MsgID, LocalID string, body []byte, err error) {
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: api.baseRequest(),
		Msg: &datastruct.SendMessage{
			Type:         datastruct.LittleVideoMsg,
			MediaID:      mediaID,
			FromUserName: fromUserName,
			ToUserName:   toUserName,
			LocalID:      tool.GetWxTimeStamp(),
			ClientMsgID:  tool.GetWxTimeStamp(),
		},
	}
	reqBody, err := json.Marshal(msgReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("fun", "async")
	params.Set("f", "json")
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsendvideomsg?"+params.Encode(), bytes.NewReader(reqBody))
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
	var smResp datastruct.SendMessageRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &smResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if smResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", smResp.BaseResponse.Ret, smResp.BaseResponse.ErrMsg)
		return
	}
	MsgID, LocalID = smResp.MsgID, smResp.LocalID
	return
}

// SendEmoticonMessage 发送动图消息
// @param fromUserName 自己的UserName
// @param toUserName 要发送的目标联系人的UserName
// @param mediaID 上传后的mediaID内容
// @return MsgID 消息的服务器ID（发送后由服务器生成）
// @return LocalID 消息本地ID（本地生成的）
func (api *wechatwebAPI) SendEmoticonMessage(fromUserName, toUserName, mediaID string) (MsgID, LocalID string, body []byte, err error) {
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: api.baseRequest(),
		Msg: &datastruct.SendMessage{
			Type:         datastruct.AnimationEmotionsMsg,
			MediaID:      mediaID,
			FromUserName: fromUserName,
			ToUserName:   toUserName,
			LocalID:      tool.GetWxTimeStamp(),
			ClientMsgID:  tool.GetWxTimeStamp(),
		},
	}
	reqBody, err := json.Marshal(msgReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("fun", "sys")
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsendemoticon?"+params.Encode(), bytes.NewReader(reqBody))
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
	var smResp datastruct.SendMessageRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &smResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if smResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", smResp.BaseResponse.Ret, smResp.BaseResponse.ErrMsg)
		return
	}
	MsgID, LocalID = smResp.MsgID, smResp.LocalID
	return
}

// SendFileMessage 发送文件消息
// @param fromUserName 自己的UserName
// @param toUserName 要发送的目标联系人的UserName
// @param mediaID 上传后的mediaID内容
// @return MsgID 消息的服务器ID（发送后由服务器生成）
// @return LocalID 消息本地ID（本地生成的）
func (api *wechatwebAPI) SendFileMessage(fromUserName, toUserName,
	mediaID, fileName string, fileSize int64) (MsgID, LocalID string, body []byte, err error) {
	msgContentReq := msgcontent.AppMsgContent{
		AppID: "wxeb7ec651dd0aefa9",
		Title: fileName,
		Type:  6,
		AppAttach: &msgcontent.AppMsgContentAppAttach{
			TotalLen: fileSize,
			AttachID: mediaID,
			FileExt:  path.Ext(fileName)[1:],
		},
	}
	msgContentBody, err := xml.Marshal(msgContentReq)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	msgReq := datastruct.SendMessageRequest{
		BaseRequest: api.baseRequest(),
		Msg: &datastruct.SendMessage{
			Type:         datastruct.AppMsg,
			Content:      string(msgContentBody),
			FromUserName: fromUserName,
			ToUserName:   toUserName,
			LocalID:      tool.GetWxTimeStamp(),
			ClientMsgID:  tool.GetWxTimeStamp(),
		},
	}
	reqBody, err := json.Marshal(msgReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("fun", "sys")
	params.Set("pass_ticket", api.loginInfo.PassTicket)
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxsendemoticon?"+params.Encode(), bytes.NewReader(reqBody))
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
	var smResp datastruct.SendMessageRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &smResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if smResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", smResp.BaseResponse.Ret, smResp.BaseResponse.ErrMsg)
		return
	}
	MsgID, LocalID = smResp.MsgID, smResp.LocalID
	return
}
