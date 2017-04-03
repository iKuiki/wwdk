package appmsg

import (
	"encoding/xml"
)

type VideoMsgContentVideo struct {
	VideoMsg       xml.Name `xml:"videomsg"`
	AesKey         string   `xml:"aeskey,attr"`
	CdnThumBaesKey string   `xml:"cdnthumbaeskey,attr"`
	CdnVideoUrl    string   `xml:"cdnvideourl,attr"`
	CdnThumbUrl    string   `xml:"cdnthumburl,attr"`
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

type VideoMsgContent struct {
	Msg      xml.Name              `xml:"msg"`
	VideoMsg *VideoMsgContentVideo `xml:"videomsg"`
}
