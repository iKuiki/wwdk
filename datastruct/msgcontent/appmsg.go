package msgcontent

import (
	"encoding/xml"
)

// AppMsgContentAppAttach app消息的content的AppAttach节点
type AppMsgContentAppAttach struct {
	XMLName  xml.Name `xml:"appattach"`
	TotalLen int64    `xml:"totallen"`
	AttachID string   `xml:"attachid"`
	FileExt  string   `xml:"fileext"`
}

// AppMsgContent app消息的content
type AppMsgContent struct {
	XMLName   xml.Name                `xml:"appmsg"`
	AppID     string                  `xml:"appid,attr"`
	SDKVer    string                  `xml:"sdkver,attr"`
	Title     string                  `xml:"title"`
	Des       string                  `xml:"des"`
	Action    string                  `xml:"action"`
	Type      int                     `xml:"type"`
	Content   string                  `xml:"content"`
	URL       string                  `xml:"url"`
	LowURL    string                  `xml:"lowurl"`
	AppAttach *AppMsgContentAppAttach `xml:"appattach"`
	ExtInfo   string                  `xml:"extinfo"`
}
