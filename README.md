# wwdk

[![License](https://img.shields.io/badge/License-GPL_3.0-blue.svg?style=flat)](LICENSE)

wwdk是web微信对接的sdk

---

### 使用方法：

基于wwdk使用微信，大致步骤为：
- 实例化一个WechatWeb对象
- 然后调用其登陆方法，
- 通过登陆channel获取到登陆二维码并且展示给用户
- 用户扫码成功后登陆channel会返回登陆成功并关闭
- 调用StartServe方法，并读取syncChannel
- 根据syncChannel读取到的状态调用对应的处理方法

---

#### 代码示范

一个简单的例子,通过终端显示二维码、收到消息后打印到终端的实现如下（仅处理文字信息）

``` golang
func main() {
	// 实例化WechatWeb对象
	wx, err := wwdk.NewWechatWeb()
	if err != nil {
		panic("Get new wechatweb client error: " + err.Error())
	}
	// 创建登陆用channel用于回传登陆信息
	loginChan := make(chan wwdk.LoginChannelItem)
	wx.Login(loginChan)
	// 根据channel返回信息进行处理
	for item := range loginChan {
		switch item.Code {
		case wwdk.LoginStatusWaitForScan:
			// 返回了登陆二维码链接，输出到屏幕
			qrterminal.Generate(item.Msg, qrterminal.L, os.Stdout)
		case wwdk.LoginStatusErrorOccurred:
			// 登陆失败
			panic(fmt.Sprintf("WxWeb Login error: %+v", item.Err))
		}
	}
	// 创建同步channel
	syncChannel := make(chan wwdk.SyncChannelItem)
	// 将channel传入startServe方法，开始同步服务并且将新信息通过syncChannel传回
	wx.StartServe(syncChannel)
	// 处理syncChannel传回信息
	for item := range syncChannel {
		// 在子方法内执行逻辑
		switch item.Code {
		// 收到新信息
		case wwdk.SyncStatusNewMessage:
			// 根据收到的信息类型分别处理
			msg := item.Message
			switch msg.MsgType {
			case datastruct.TextMsg:
				// 处理文字信息
				processTextMessage(wx, msg)
			}
		case wwdk.SyncStatusPanic:
			// 发生致命错误，sync中断
			fmt.Printf("sync panic: %+v\n", err)
			break
		}
	}
}

func processTextMessage(app *wwdk.WechatWeb, msg *datastruct.Message) {
	from, err := app.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
		return
	}
	log.Printf("Recived a text msg from %s: %s", from.NickName, msg.Content)
}
```

---

#### 详细example

一个详细的例子包含储存登陆信息（可用于程序停止后重新运行免登录，文件名loginInfo.txt)、通过终端显示二维码、收到消息后打印到终端、收到图片、视频、音频保存到运行目录并打印文件名到终端

代码详见：[example](https://github.com/iKuiki/wwdk/blob/master/example/main.go)
