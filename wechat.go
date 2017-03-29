package wxweb

import (
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
)

type wechatCookie struct {
	Skey       string
	Wxsid      string
	Wxuin      string
	Uvid       string
	DataTicket string
	AuthTicket string
}

type WechatWeb struct {
	cookie      wechatCookie
	userAgent   string
	deviceId    string
	contactList []*datastruct.Contact
	user        *datastruct.User
	syncKey     *datastruct.SyncKey
}

func NewWechatWeb() (wxweb WechatWeb) {
	return WechatWeb{
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		deviceId:  "e" + tool.GetRandomStringFromNum(15),
	}
}
