package wwdk

import (
	"github.com/pkg/errors"

	"github.com/ikuiki/wwdk/datastruct"
)

// SaveMessageImage 保存消息图片到指定位置
func (wxwb *WechatWeb) SaveMessageImage(msg datastruct.Message) (filename string, err error) {
	d, err := wxwb.api.SaveMessageImage(msg.MsgID)
	if err != nil {
		return
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
	d, err := wxwb.api.SaveMessageVoice(msg.MsgID)
	if err != nil {
		return
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
	d, err := wxwb.api.SaveMessageVideo(msg.MsgID)
	if err != nil {
		return
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
	d, err := wxwb.api.SaveContactImg(contact.HeadImgURL)
	if err != nil {
		return
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
	d, err := wxwb.api.SaveContactImg(user.HeadImgURL)
	if err != nil {
		return
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
// TODO: delete
func (wxwb *WechatWeb) SaveMemberImg(member datastruct.Member, chatroomID string) (filename string, err error) {
	d, err := wxwb.api.SaveMemberImg(member.UserName, chatroomID)
	if err != nil {
		return
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
