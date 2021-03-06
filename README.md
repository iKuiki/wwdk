# wwdk

[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat)](LICENSE)

wwdk是web微信对接的sdk，主要专注于实现微信机器人相关的逻辑，关注文字、音频、视频、图片的对接

功能列表

- [x] 获取好友列表
- [x] 批量获取联系人信息
- [ ] ~~添加好友~~ 网页版微信好像现在没有这个能力了
- [ ] 接受好友请求
- [x] 修改好友备注
- [ ] ~~删除好友~~ 网页版微信好像没有这个能力
- [ ] 好友拉群
- [ ] 群成员添加
- [ ] 群成员移除（如果为管理员
- [x] 修改群名称
- [x] 文字信息收发
- [x] 接收图片
- [ ] 发送图片
- [x] 接收音频
- [ ] ~~发送音频~~ 网页版没有发送语音的功能
- [x] 接收视频
- [ ] 发送视频
- [x] 接收动图
- [ ] 发送动图
- [ ] 接收文件
- [ ] 发送文件
- [ ] 接收名片
- [ ] ~~发送微信名片~~ 网页版好像已经不再提供这个功能
- [ ] ~~发送公众号名片~~ 网页版好像已经不再提供这个功能
- [ ] ~~发送图文消息~~ 网页版好像已经不再提供这个功能
- [ ] 识别红包消息
- [ ] 识别小程序消息

另外我对web微信协议的分析，写在这里：[web微信协议分析](doc/protocol/README.md)

TODO：

- [ ] 实现微信加好友请求的处理
- [ ] 实现更多信息类型的兼容
- [x] 完善doc下的微信web协议记录
- [ ] 提供一个错误追踪注入（例如Sentry之类的）来便于debug

## 目前存在的问题

- ~~部分微信能正常登陆网页版但是用此程序无法登陆~~根据多版本domain理论已经解决此问题(**不过那些因为账号原因被腾讯禁止使用网页版微信的，肯定还是登陆不了的**)

---

## 使用方法

网页版微信运行步骤图如下

``` asciiflow
     +-------------------+    +-------------------+    +-------------------+
     |                   |    |                   |    |                   |
     |  Start web wechat +--->+   Wait for scan   +--->+  Login success    |
     |                   |    |                   |    |                   |
     +-------------------+    +-------------------+    +------------+------+
                                                                    |
                                                                    |
                                                                    |
                                                                    v
                                                        +-----------+-----+
                                                        |                 |
                                                        |   Get contact   +-------+
                                                        |                 |       |
                                                        +---------+-------+       |
                                                                  |               |
                                                                  |               v
                                                                  |      +--------+-----+
                                                                  |      |              |
                +------------------------------------------+      |      |  User input  |
                |                                          |      |      |              |
                |                                          |      |      +---------+----+
                |                                          |      |                |
                |                                          v      v                v
     +----------+---------+     +----------------+     +---+------+-----+   +------+-------+
     |                    |     |                |     |                |   |              |
+----+ Process new events +<----+  Web wx sync   +<----+   Sync check   |   | Send message |
|    |                    |     |                |     |                |   |              |
|    +--------------------+     +------+---------+     +-------+--------+   +-----+--------+
|                                      ^                       ^                  |
|                                      |                       |                  |
|    +------------+                    |                       |                  |
+----+ Message    |                    |                       |                  |
|    +------------+                    |                       |                  |
|                                      |                       |                  |
|    +----------------+                |                       |                  |
+----+ Modify contact |                +-----------------------+                  |
|    +----------------+                                        |                  |
|                                                              |                  |
|    +----------------+                                        |                  |
+----+ Delete contact |                                        |                  |
     +----------------+                                        |                  |
                                                               |                  |
                                                               |                  v
     +---------------------------------------------------------+------------------+--------+
     |                                                                                     |
     |                                  Wechat Server                                      |
     |                                                                                     |
     +-------------------------------------------------------------------------------------+
```

从上面可以看出，web微信整体总共可以大概分为3部分，第一部分是登陆，第二部分是同步消息，第三部分是发送消息。
web微信打开后，首先我们需要通过用户扫码登陆，这一部分被封装为wwdk的Login方法，并且会将login进度通过channel传递出来。当从channel中检测到登陆成功后，就可以开始启动sync服务，以及为用户提供发送消息的服务了。sync服务也会给一个channel，新消息会从channel中返回。发送消息时直接调用wwdk的发送消息方法即可

基于wwdk使用微信，大致步骤为：

- 实例化一个WechatWeb对象
- 然后调用其登陆方法
- 通过登陆channel获取到登陆二维码并且展示给用户
- 用户扫码成功后登陆channel会返回登陆成功并关闭
- 调用StartServe方法，并读取syncChannel
- 根据syncChannel读取到的状态调用对应的处理方法

---

## 代码示范

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

## 详细example

一个详细的例子包含储存登陆信息（可用于程序停止后重新运行免登录，文件名loginInfo.txt)、通过终端显示二维码、收到消息后打印到终端、收到图片、视频、音频保存到运行目录并打印文件名到终端

代码详见：[example](https://github.com/iKuiki/wwdk/blob/master/example/main.go)
