package wxweb

type Context struct {
	App     *WechatWeb
	hasStop bool
}

func (this *Context) Stop() {
	this.hasStop = true
}
