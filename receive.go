package wwdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ikuiki/wwdk/datastruct"
)

// getMessage 同步消息
// 如果同步状态接口返回有新消息需要同步，通过此接口从服务器中获取新消息
func (wxwb *WechatWeb) getMessage() (gmResp datastruct.GetMessageRespond, err error) {
	gmResp = datastruct.GetMessageRespond{}
	data, err := json.Marshal(datastruct.GetMessageRequest{
		BaseRequest: wxwb.baseRequest(),
		SyncKey:     wxwb.loginInfo.syncKey,
		Rr:          ^time.Now().Unix() + 1,
	})
	if err != nil {
		return gmResp, errors.New("Marshal request body to json fail: " + err.Error())
	}
	params := url.Values{}
	params.Set("sid", wxwb.loginInfo.cookie.Wxsid)
	params.Set("skey", wxwb.loginInfo.sKey)
	// params.Set("pass_ticket", wxwb.PassTicket)
	resp, err := wxwb.apiRuntime.client.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsync?"+params.Encode(),
		"application/json;charset=UTF-8",
		bytes.NewReader(data))
	if err != nil {
		return gmResp, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &gmResp)
	if err != nil {
		return gmResp, errors.New("Unmarshal respond json fail: " + err.Error())
	}
	if gmResp.BaseResponse.Ret != 0 {
		return gmResp, errors.New("respond error ret: " + strconv.FormatInt(gmResp.BaseResponse.Ret, 10))
	}
	// if gmResp.AddMsgCount > 0 {
	// 	wxwb.logger.Debug(string(resp)+"\n")
	// 	panic(nil)
	// }
	return gmResp, nil
}

// SaveMessageImage 保存消息图片到指定位置
func (wxwb *WechatWeb) SaveMessageImage(msg datastruct.Message) (filename string, err error) {
	params := url.Values{}
	params.Set("MsgID", msg.MsgID)
	params.Set("skey", wxwb.loginInfo.sKey)
	// params.Set("type", "slave")
	resp, err := wxwb.apiRuntime.client.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetmsgimg?" + params.Encode())
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Read io.ReadCloser error: " + err.Error())
	}
	filename, err = wxwb.mediaStorer.Storer(MediaFile{
		MediaType:     MediaTypeMessageImage,
		FileName:      msg.MsgID + ".png",
		BinaryContent: d,
	})
	if err != nil {
		return "", errors.New("mediaStorer.Storer error: " + err.Error())
	}
	return filename, nil
}

// SaveMessageVoice 保存消息声音到指定位置
func (wxwb *WechatWeb) SaveMessageVoice(msg datastruct.Message) (filename string, err error) {
	if msg.MsgType != datastruct.VoiceMsg {
		return "", errors.New("Message type wrong")
	}
	params := url.Values{}
	params.Set("MsgID", msg.MsgID)
	params.Set("skey", wxwb.loginInfo.sKey)
	// params.Set("type", "slave")
	resp, err := wxwb.apiRuntime.client.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvoice?" + params.Encode())
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Read io.ReadCloser error: " + err.Error())
	}
	filename, err = wxwb.mediaStorer.Storer(MediaFile{
		MediaType:     MediaTypeMessageVoice,
		FileName:      msg.MsgID + ".mp3",
		BinaryContent: d,
	})
	if err != nil {
		return "", errors.New("mediaStorer.Storer error: " + err.Error())
	}
	return filename, nil
}

// SaveMessageVideo 保存消息视频到指定位置
func (wxwb *WechatWeb) SaveMessageVideo(msg datastruct.Message) (filename string, err error) {
	if msg.MsgType != datastruct.LittleVideoMsg {
		return "", errors.New("Message type wrong")
	}
	params := url.Values{}
	params.Set("msgid", msg.MsgID)
	params.Set("skey", wxwb.loginInfo.sKey)
	req, err := http.NewRequest("GET", "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvideo?"+params.Encode(), strings.NewReader(""))
	if err != nil {
		return "", errors.New("create request error: " + err.Error())
	}
	req.Header.Set("Range", "bytes=0-")
	resp, err := wxwb.apiRuntime.client.Do(req)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Read io.ReadCloser error: " + err.Error())
	}
	filename, err = wxwb.mediaStorer.Storer(MediaFile{
		MediaType:     MediaTypeMessageVideo,
		FileName:      msg.MsgID + ".mp4",
		BinaryContent: d,
	})
	if err != nil {
		return "", errors.New("mediaStorer.Storer error: " + err.Error())
	}
	return filename, nil
}

// SaveContactImg 保存联系人头像
func (wxwb *WechatWeb) SaveContactImg(contact datastruct.Contact) (filename string, err error) {
	req, err := http.NewRequest("GET", "https://wx2.qq.com"+contact.HeadImgURL+wxwb.loginInfo.sKey, nil)
	if err != nil {
		return "", errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.apiRuntime.client.Do(req)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Read io.ReadCloser error: " + err.Error())
	}
	filename, err = wxwb.mediaStorer.Storer(MediaFile{
		MediaType:     MediaTypeContactHeadImg,
		FileName:      contact.UserName + ".png",
		BinaryContent: d,
	})
	if err != nil {
		return "", errors.New("mediaStorer.Storer error: " + err.Error())
	}
	return filename, nil
}

// SaveUserImg 保存登陆用户的头像
func (wxwb *WechatWeb) SaveUserImg(user datastruct.User) (filename string, err error) {
	req, err := http.NewRequest("GET", "https://wx2.qq.com"+user.HeadImgURL, nil)
	if err != nil {
		return "", errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.apiRuntime.client.Do(req)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Read io.ReadCloser error: " + err.Error())
	}
	filename, err = wxwb.mediaStorer.Storer(MediaFile{
		MediaType:     MediaTypeUserHeadImg,
		FileName:      user.UserName + ".png",
		BinaryContent: d,
	})
	if err != nil {
		return "", errors.New("mediaStorer.Storer error: " + err.Error())
	}
	return filename, nil
}

// SaveMemberImg 保存群成员的头像
func (wxwb *WechatWeb) SaveMemberImg(member datastruct.Member, chatroomID string) (filename string, err error) {
	req, err := http.NewRequest("GET", "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgeticon?seq=0&username="+member.UserName+"&chatroomid="+chatroomID+"&skey=", nil)
	if err != nil {
		return "", errors.New("create request error: " + err.Error())
	}
	resp, err := wxwb.apiRuntime.client.Do(req)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Read io.ReadCloser error: " + err.Error())
	}
	filename, err = wxwb.mediaStorer.Storer(MediaFile{
		MediaType:     MediaTypeMemberHeadImg,
		FileName:      member.UserName + ".png",
		BinaryContent: d,
	})
	if err != nil {
		return "", errors.New("mediaStorer.Storer error: " + err.Error())
	}
	return filename, nil
}
