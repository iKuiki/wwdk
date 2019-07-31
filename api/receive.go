package api

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// SaveMessageImage 下载图片消息
// 将消息的图片下载回来
// @param msgID 要下载的图片消息的MsgID
// @return imgData 下载到的图片的二进制数据
func (api *wechatwebAPI) SaveMessageImage(msgID string) (imgData []byte, err error) {
	params := url.Values{}
	params.Set("MsgID", msgID)
	params.Set("skey", api.loginInfo.SKey)
	// params.Set("type", "slave")
	req, err := http.NewRequest("GET", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxgetmsgimg?"+params.Encode(), nil)
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	// resp, err := api.client.Get("https://" + api.apiDomain + "/cgi-bin/mmwebwx-bin/webwxgetmsgimg?" + params.Encode())
	// if err != nil {
	// 	err = errors.New("request error: " + err.Error())
	// 	return
	// }
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("do request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	imgData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Read io.ReadCloser error: " + err.Error())
		return
	}
	return
}

// SaveMessageVoice 下载音频消息
// 将消息的音频下载回来
// @param msgID 要下载的音频消息的MsgID
// @return imgData 下载到的音频的二进制数据
func (api *wechatwebAPI) SaveMessageVoice(msgID string) (voiceData []byte, err error) {
	params := url.Values{}
	params.Set("MsgID", msgID)
	params.Set("skey", api.loginInfo.SKey)
	req, err := http.NewRequest("GET", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxgetvoice?"+params.Encode(), nil)
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	req.Header.Set("Range", "bytes=0-")
	// resp, err := api.client.Get("https://" + api.apiDomain + "/cgi-bin/mmwebwx-bin/webwxgetvoice?" + params.Encode())
	// if err != nil {
	// 	return "", errors.New("request error: " + err.Error())
	// }
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("do request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	voiceData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Read io.ReadCloser error: " + err.Error())
		return
	}
	return
}

// SaveMessageVideo 下载视频消息
// 将消息的视频下载回来
// @param msgID 要下载的视频消息的MsgID
// @return videoData 下载到的视频的二进制数据
func (api *wechatwebAPI) SaveMessageVideo(msgID string) (videoData []byte, err error) {
	params := url.Values{}
	params.Set("msgid", msgID)
	params.Set("skey", api.loginInfo.SKey)
	req, err := http.NewRequest("GET", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxgetvideo?"+params.Encode(), nil)
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	req.Header.Set("Range", "bytes=0-")
	resp, err := api.client.Do(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	videoData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Read io.ReadCloser error: " + err.Error())
		return
	}
	return
}

// SaveContactImg 保存联系人头像
// 根据头像地址保存头像
// @param headImgURL 联系人头像地址
// @return imgData 下载的头像的二进制图片数据
func (api *wechatwebAPI) SaveContactImg(headImgURL string) (imgData []byte, err error) {
	// 貌似时不需要带skey的，不如做个判断好了
	if strings.HasSuffix(headImgURL, "&skey=") {
		headImgURL = headImgURL + api.loginInfo.SKey
	}
	req, err := http.NewRequest("GET", "https://"+api.apiDomain+headImgURL, nil)
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	resp, err := api.client.Do(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	imgData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Read io.ReadCloser error: " + err.Error())
		return
	}
	return
}

// SaveMemberImg 保存群成员的头像
// @param userName 群成员的UserName
// @param chatroomID 群的EncryChatRoomId
// @return imgData 下载的头像的二进制图片数据
func (api *wechatwebAPI) SaveMemberImg(userName, chatroomID string) (imgData []byte, err error) {
	req, err := http.NewRequest("GET", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxgeticon?seq=0&username="+userName+"&chatroomid="+chatroomID+"&skey=", nil)
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	resp, err := api.client.Do(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	imgData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Read io.ReadCloser error: " + err.Error())
		return
	}
	return
}
