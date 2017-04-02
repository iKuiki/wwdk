package appmsg

import (
	"encoding/xml"
)

type EmotionMsgContentEmoji struct {
	Emoji        xml.Name `xml:"emoji"`
	FromUserName string   `xml:"fromusername,attr"`
	ToUserName   string   `xml:"tousername,attr"`
	Type         string   `xml:"type,attr"`
	IdBuffer     string   `xml:"idbuffer,attr"`
	Md5          string   `xml:"md5,attr"`
	Len          string   `xml:"len,attr"`
	ProductId    string   `xml:"productid,attr"`
	AndroidMd5   string   `xml:"androidmd5,attr"`
	AndroidLen   string   `xml:"androidlen,attr"`
	S60v3Md5     string   `xml:"s60v3md5,attr"`
	S60v3Len     string   `xml:"s60v3len,attr"`
	S60v5Md5     string   `xml:"s60v5md5,attr"`
	S60v5Len     string   `xml:"s60v5len,attr"`
	CdnUrl       string   `xml:"cdnurl,attr"`
	DesignerId   string   `xml:"designerid,attr"`
	ThumbUrl     string   `xml:"thumburl,attr"`
	EncryptUrl   string   `xml:"encrypturl,attr"`
	AesKey       string   `xml:"aeskey,attr"`
	Width        string   `xml:"width,attr"`
	Height       string   `xml:"height,attr"`
}

type EmotionMsgContentGameext struct {
	Gameext xml.Name `xml:"gameext"`
	Type    string   `xml:"type,attr"`
	Content string   `xml:"content,attr"`
}

type EmotionMsgContent struct {
	Msg     xml.Name                  `xml:"msg"`
	Emoji   *EmotionMsgContentEmoji   `xml:"emoji"`
	Gameext *EmotionMsgContentGameext `xml:"gameext"`
}
