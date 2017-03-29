package wxweb

type wechatCookie struct {
	Skey       string
	Wxsid      string
	Wxuin      string
	Uvid       string
	DataTicket string
	AuthTicket string
}

type WechatWeb struct {
	cookie wechatCookie
}
