package msgcontent

import (
	"encoding/xml"
)

// VideoMsgContentVideo 视频消息的Content中的video节点
type VideoMsgContentVideo struct {
	XMLName        xml.Name `xml:"videomsg"`
	AesKey         string   `xml:"aeskey,attr"`
	CdnThumBaesKey string   `xml:"cdnthumbaeskey,attr"`
	CdnVideoURL    string   `xml:"cdnvideourl,attr"`
	CdnThumbURL    string   `xml:"cdnthumburl,attr"`
	Length         string   `xml:"length,attr"`
	PlayLength     string   `xml:"playlength,attr"`
	CdnThumbLength string   `xml:"cdnthumblength,attr"`
	CdnThumbWidth  string   `xml:"cdnthumbwidth,attr"`
	CdnThumbHeight string   `xml:"cdnthumbheight,attr"`
	FromUserName   string   `xml:"fromusername,attr"`
	Md5            string   `xml:"md5,attr"`
	NewMd5         string   `xml:"newmd5,attr"`
	IsAd           string   `xml:"isad,attr"`
}

// VideoMsgContent 视频消息的Content
type VideoMsgContent struct {
	XMLName  xml.Name              `xml:"msg"`
	VideoMsg *VideoMsgContentVideo `xml:"videomsg"`
}
