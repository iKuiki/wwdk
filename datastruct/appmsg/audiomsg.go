package appmsg

import (
	"encoding/xml"
)

// VoiceMsgContentVoice 音频消息的content的voice节点
type VoiceMsgContentVoice struct {
	VoiceMsg     xml.Name `xml:"voicemsg"`
	EndFlag      string   `xml:"endflag,attr"`
	Length       string   `xml:"length,attr"`
	VoiceLength  string   `xml:"voicelength,attr"`
	ClientMsgID  string   `xml:"clientmsgid,attr"`
	FromUserName string   `xml:"fromusername,attr"`
	DownCount    string   `xml:"downcount,attr"`
	CancelFlag   string   `xml:"cancelflag,attr"`
	VoiceFormat  string   `xml:"voiceformat,attr"`
	ForwardFlag  string   `xml:"forwardflag,attr"`
	BufID        string   `xml:"bufid,attr"`
}

// VoiceMsgContent 音频消息的content
type VoiceMsgContent struct {
	Msg      xml.Name              `xml:"msg"`
	VoiceMsg *VoiceMsgContentVoice `xml:"voicemsg"`
}
