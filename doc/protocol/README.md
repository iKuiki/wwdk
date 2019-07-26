# 微信网页版通信协议

- [微信网页版通信协议](#%e5%be%ae%e4%bf%a1%e7%bd%91%e9%a1%b5%e7%89%88%e9%80%9a%e4%bf%a1%e5%8d%8f%e8%ae%ae)
  - [文档已经编写完成的功能](#%e6%96%87%e6%a1%a3%e5%b7%b2%e7%bb%8f%e7%bc%96%e5%86%99%e5%ae%8c%e6%88%90%e7%9a%84%e5%8a%9f%e8%83%bd)
  - [web微信的工作流程](#web%e5%be%ae%e4%bf%a1%e7%9a%84%e5%b7%a5%e4%bd%9c%e6%b5%81%e7%a8%8b)
  - [一种特别的返回格式](#%e4%b8%80%e7%a7%8d%e7%89%b9%e5%88%ab%e7%9a%84%e8%bf%94%e5%9b%9e%e6%a0%bc%e5%bc%8f)
  - [关于cookie](#%e5%85%b3%e4%ba%8ecookie)
  - [关于Domain（极其重要）](#%e5%85%b3%e4%ba%8edomain%e6%9e%81%e5%85%b6%e9%87%8d%e8%a6%81)

---

## 文档已经编写完成的功能

- [x] 获取好友列表
- [x] 批量获取联系人信息
- [ ] ~~添加好友~~ 网页版微信好像现在没有这个能力了
- [x] 接受好友请求
- [x] 修改好友备注
- [ ] ~~删除好友~~ 网页版微信好像没有这个能力
- [x] 好友拉群
- [x] 群成员添加
- [x] 群成员移除（如果为管理员
- [x] 修改群名称
- [x] 文字信息收发
- [x] 接收图片
- [ ] 发送图片
- [ ] 接收视频
- [ ] 发送视频
- [ ] 接收微信名片
- [ ] 发送微信名片
- [ ] 接收公众号名片
- [ ] 发送公众号名片
- [ ] 发送图文消息
- [x] 接收文件
- [ ] 发送文件

## web微信的工作流程

web微信打开后，会先让用户[登陆](login.md)，当登陆完成后就进入[同步](sync.md)阶段，主进程监听服务器的新消息，如果有消息则显示给用户，如果用户发送消息，[发送消息](send.md)时会通过事件函数将发送的消息提交给服务器，两个进程互相独立

``` asciiflow
                 +-----------+
                 |   Start   |
                 +-----+-----+
                       |
                       |
                       v
                 +-----+-----+
                 |   Login   |
                 +-----+-----+
                       |                               +----------------+
       +---------------+------------+                  |                |
       |                            |                  |                |
       |                            v                  |                |
+------+-------+              +-----+-----+            |                |
|  User Input  |         +--->+   Sync    | <----------+     Wechat     |
+------+-------+         |    +-----+-----+            |                |
       |                 |          |                  |     Server     |
       |                 |          |                  |                |
       |                 +----------+                  |                |
       |                                               |                |
       |                                               |                |
       |                                               |                |
       +---------------------------------------------->+                |
                                                       +----------------+

```

---

## 一种特别的返回格式

*注：微信网页版API的返回包括一种特别的格式：看起来像js代码，每个字段作为一行js代码，以分号结尾，每句以等号分割左边为key右边为code*
例：

``` js
window.QRLogin.code = 200; window.QRLogin.uuid = "gfNHoe0rgA==";
```

其中包含两个值：

| Key                 | Value        |
| ------------------- | ------------ |
| window.QRLogin.code | 200          |
| window.QRLogin.uuid | gfNHoe0rgA== |

解析方案：

1. 使用[github.com/robertkrimen/otto](https://github.com/robertkrimen/otto)解释后获取值（优点：可靠性高；缺点：因为要运行js，解析速度相对慢）
2. 自己通过匹配格式来解析值（优点：简单，速度快；缺点：若返回值复杂则可能解析错误）

---

## 关于cookie

登陆成功后，大部分请求都需要带cookie访问，*建议所有请求使用一个实现了cookie功能的client来执行请求，让client自动接受、发送cookie*
另外部分带json body的请求会需要在json中内置一个BaseRequest，格式一般如下：

``` json
    "BaseRequest": {
        "Uin": 1880000000,
        "Sid": "LAxxxxxxxxxxxxx",
        "Skey": "@crypt_7dxxxxxx_e84xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e330000000000000"
    },
```

---

## 关于Domain（极其重要）

经过研究发现，网页版微信会根据微信账号的不同，分配不同的domain来处理（可能是出于负载均衡的目的），而分配domain以后，之后的所有请求都必须用配套的domain来请求，否则就会出错。这个坑，我是对我2个微信号分析了很久才发现规律，一开始只是觉得很乱，然后偶尔一两个接口用错domain也没关系，但是有些关键接口用错domain则会直接导致失败。而这个domain的来源，在于用户扫码登陆后获取到的redirectURL中的domain。这个获取到的domain，我们在下文中就称它为apiDomain好了，与之配套的还有上传文件用的uploadDomain、接收推送信息的syncDomain，他们3个都需要相互匹配才能保证功能正常。

apiDomain、uploadDomain、syncDomain的相互关系为：```uploadDomain="file."+apiDomain; syncDomain="webpush."+syncDomain```

另外过程中还发现，根据domain的不同，某些接口的参数也会有细微差异，体现在字段pass_ticket上，部分domain有传此参数，而部分domain没用传此参数，出于通用考虑，可以统一都传此参数，目前暂未发现不良影响

另外，登陆前其实还会用到一个loginDomain，不过那个是登陆之前，还没被分配domain，所以用混是无所谓的。```loginDomain="login."+apiDomain```
