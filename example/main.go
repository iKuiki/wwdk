package main

import (
	"github.com/yinhui87/wechat-web"
	"log"
)

func main() {
	wx := wxweb.WechatWeb{}
	err := wx.Login()
	if err != nil {
		log.Printf("WxWeb Login error: %s\n", err.Error())
	}
	wx.StartServe()
}
