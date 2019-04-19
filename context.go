package wwdk

// Context 一次处理流程的上下文
type Context struct {
	App     *WechatWeb
	hasStop bool
}

// Stop 终止当前处理流程
func (context *Context) Stop() {
	context.hasStop = true
}
