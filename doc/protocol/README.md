# 微信网页版通信协议

- [微信网页版通信协议](#%e5%be%ae%e4%bf%a1%e7%bd%91%e9%a1%b5%e7%89%88%e9%80%9a%e4%bf%a1%e5%8d%8f%e8%ae%ae)
  - [web微信的工作流程](#web%e5%be%ae%e4%bf%a1%e7%9a%84%e5%b7%a5%e4%bd%9c%e6%b5%81%e7%a8%8b)
  - [一种特别的返回格式](#%e4%b8%80%e7%a7%8d%e7%89%b9%e5%88%ab%e7%9a%84%e8%bf%94%e5%9b%9e%e6%a0%bc%e5%bc%8f)
  - [关于cookie](#%e5%85%b3%e4%ba%8ecookie)

---

## web微信的工作流程

web微信打开后，会先让用户登陆，当登陆完成后就进入同步阶段，主进程监听服务器的新消息，如果有消息则显示给用户，如果用户发送消息，发送消息时会通过事件函数将发送的消息提交给服务器，两个进程互相独立

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
