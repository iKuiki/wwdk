package appmsg

import (
	"encoding/xml"
)

type VoiceMsgContentVoice struct {
	VoiceMsg     xml.Name `xml:"voicemsg"`
	EndFlag      string   `xml:"endflag,attr"`
	Length       string   `xml:"length,attr"`
	VoiceLength  string   `xml:"voicelength,attr"`
	ClientMsgId  string   `xml:"clientmsgid,attr"`
	FromUserName string   `xml:"fromusername,attr"`
	DownCount    string   `xml:"downcount,attr"`
	CancelFlag   string   `xml:"cancelflag,attr"`
	VoiceFormat  string   `xml:"voiceformat,attr"`
	ForwardFlag  string   `xml:"forwardflag,attr"`
	BufId        string   `xml:"bufid,attr"`
}

type VoiceMsgContent struct {
	Msg      xml.Name              `xml:"msg"`
	VoiceMsg *VoiceMsgContentVoice `xml:"voicemsg"`
}
