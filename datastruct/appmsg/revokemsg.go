package appmsg

import (
	"encoding/xml"
)

type RevokeMsgContentRevoke struct {
	RevokeMsg  xml.Name `xml:"revokemsg"`
	Session    string   `xml:"session"`
	OldMsgId   string   `xml:"oldmsgid"`
	MsgId      string   `xml:"msgid"`
	ReplaceMsg string   `xml:"replacemsg"`
}

type RevokeMsgContent struct {
	SysMsg    xml.Name                `xml:"sysmsg"`
	RevokeMsg *RevokeMsgContentRevoke `xml:"revokemsg"`
}
