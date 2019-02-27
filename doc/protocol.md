# 微信网页版通信协议

- [微信网页版通信协议](#%E5%BE%AE%E4%BF%A1%E7%BD%91%E9%A1%B5%E7%89%88%E9%80%9A%E4%BF%A1%E5%8D%8F%E8%AE%AE)
  - [通用](#%E9%80%9A%E7%94%A8)
  - [登陆](#%E7%99%BB%E9%99%86)
    - [获取uuid](#%E8%8E%B7%E5%8F%96uuid)
    - [轮询用户扫码](#%E8%BD%AE%E8%AF%A2%E7%94%A8%E6%88%B7%E6%89%AB%E7%A0%81)
    - [获取登陆参数](#%E8%8E%B7%E5%8F%96%E7%99%BB%E9%99%86%E5%8F%82%E6%95%B0)
    - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
  - [接收](#%E6%8E%A5%E6%94%B6)
    - [同步状态](#%E5%90%8C%E6%AD%A5%E7%8A%B6%E6%80%81)
  - [发送](#%E5%8F%91%E9%80%81)
    - [消息已读](#%E6%B6%88%E6%81%AF%E5%B7%B2%E8%AF%BB)
    - [发送消息](#%E5%8F%91%E9%80%81%E6%B6%88%E6%81%AF)
      - [Body (Json):](#body-json)
    - [撤回消息](#%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF)
    - [转发被撤回的图片消息](#%E8%BD%AC%E5%8F%91%E8%A2%AB%E6%92%A4%E5%9B%9E%E7%9A%84%E5%9B%BE%E7%89%87%E6%B6%88%E6%81%AF)
    - [修改用户备注](#%E4%BF%AE%E6%94%B9%E7%94%A8%E6%88%B7%E5%A4%87%E6%B3%A8)
    - [修改群名](#%E4%BF%AE%E6%94%B9%E7%BE%A4%E5%90%8D)

## 通用

*注：微信网页版API的返回包括一种特别的格式：看起来像js代码，每个字段作为一行js代码，以分号结尾，每句以等号分割左边为key右边为code*
例：
```
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

**登陆成功后，请求中凡需要Cookie的，都需要包含以下Cookie字段**

| 需要的Cookie      |
| ----------------- |
| webwxuvid         |
| webwx_auth_ticket |
| wxuin             |
| wxsid             |
| webwx_data_ticket |

---

## 登陆

登陆流程大致为以下步骤：
- 获取uuid
- 根据uuid生成二维码
- 轮询微信服务器检测用户是否扫码
- 手机扫码
- 获取到用户已扫码，得到用户头像（base64编码的字符串）
- 用户同意登陆
- 登陆成功

---

### 获取uuid

| Key         | Value                           | Remark               |
| ----------- | ------------------------------- | -------------------- |
| Request URL | https://login.wx.qq.com/jslogin |                      |
| Method      | Get                             |                      |
| Cookie      | No                              |                      |
| Param       | appid                           | 填wx782c26e4c19acffb |
| Param       | fun                             | 填new                |
| Param       | lang                            | zh_CN或en_US         |
| Param       | _                               | 13位unix时间戳       |

**Response:**

| Key                 | Value | Remark         |
| ------------------- | ----- | -------------- |
| window.QRLogin.code | 200   |                |
| window.QRLogin.uuid | xxx   | 当前会话的uuid |

*Example:*
``` javascript
window.QRLogin.code = 200; window.QRLogin.uuid = "gfNHoe0rgA==";
```

---

### 轮询用户扫码

| Key         | Value                                              | Remark                     |
| ----------- | -------------------------------------------------- | -------------------------- |
| Request URL | https://login.wx2.qq.com/cgi-bin/mmwebwx-bin/login |                            |
| Method      | Get                                                |                            |
| Param       | loginicon                                          | 填true                     |
| Param       | uuid                                               | 之前获取的uuid             |
| Param       | tip                                                | 1-未扫描 0-已扫描          |
| Param       | r                                                  | 13位时间戳取反(貌似可省略) |
| Param       | _                                                  | 13位unix时间戳             |

**Response:**

| Key               | Value                       | Remark                                                    |
| ----------------- | --------------------------- | --------------------------------------------------------- |
| window.code       | 200<br/>201<br/>400<br/>408 | 确认登陆<br/>已扫码<br/>登陆超时(二维码失效)<br/>等待登陆 |
| window.userAvatar | data:img/jpeg;base64        | base64编码的用户头像，仅当code=200时才有                  |

*Example:*
``` javascript
window.code=408;window.userAvatar='data:img/jpeg;base64,iVBORw...'
```

**若用户取消登陆，返回仍为408，旧的二维码仍可重复使用，用户重新扫旧的二维码后会再次返回201**

*注：当登陆超时（二维码失效）后，重新调用获取uuid的方法即可重新拿到二维码*

---

### 获取登陆参数

| Key         | Value                                                    | Remark                 |
| ----------- | -------------------------------------------------------- | ---------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage |                        |
| Method      | Get                                                      |                        |
| Param       | ticket                                                   |                        |
| Param       | uuid                                                     |                        |
| Param       | scan                                                     | 扫描成功后返回的时间戳 |
| Param       | fun                                                      | 填new                  |
| Param       | version                                                  | 填v2                   |

**Response:**

| Key               | Type    | Remark         |
| ----------------- | ------- | -------------- |
| wxsid             | Cookie  |                |
| wxuin             | Cookie  |                |
| webwxuvid         | Cookie  |                |
| webwx_auth_ticket | Cookie  |                |
| webwx_data_ticket | Cookie  |                |
| skey              | BodyXml |                |
| wxsid             | BodyXml | same as cookie |
| wxuin             | BodyXml | same as cookie |
| pass_ticket       | BodyXml |                |

---

### 初始化

| Key         | Value                                            | Remark                 |
| ----------- | ------------------------------------------------ | ---------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit |                        |
| Method      | Post                                             |                        |
| Cookie      | Need                                             |                        |
| Param       | r                                                | 13位时间戳取反         |
| Param       | pass_ticket                                      | 获取登陆参数时获取到的 |

*Body (json):*
``` json
{"BaseRequest":
    {
        "Uin":"216547950",
        "Sid":"QQ9iwKokvmPs7c/7",
        "Skey":"@crypt_a6a25b34_68efb91dcbec1fe990bf33d8fe770034",
        "DeviceID":"e987736822175688"
    }
}
```

**Response:**
返回为一个json对象，内包括用户信息、联系人(此列表不全，之后用获取联系人的接口获取完整联系人列表)、同步信息等
*该json中主要需要用到的数据为User信息与synckey，synckey用于后续同步状态时使用，而联系人可以之后通过获取联系人接口获取完整列表*

到此登陆就成功了

## 接收

接受消息主要分为2个步骤
- 轮询微信状态同步接口，获得当前微信状态
- 如果当前微信状态为有新消息，则调用获取新消息接口取得新消息

**同步状态接口中有个synckey，这个synckey是在登陆时的初始化(webwxinit)操作中获取的，并且每次调用getMessage接口都会刷新这个synckey**

---

### 同步状态

| Key         | Value                                                    | Remark                                                |
| ----------- | -------------------------------------------------------- | ----------------------------------------------------- |
| Request URL | https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck |                                                       |
| Method      | Get                                                      |                                                       |
| Cookie      | Need                                                     |                                                       |
| Param       | r                                                        | 13位unix时间戳                                        |
| Param       | skey                                                     | 登陆凭据skey                                          |
| Param       | sid                                                      | 登陆凭据sid                                           |
| Param       | uin                                                      | 登陆凭据uin                                           |
| Param       | deviceid                                                 |                                                       |
| Param       | synckey                                                  | 同步key，由"\|"分割为1_xxx\|2_xxx\|3_xxx\|4_xxx的格式 |
| Param       | _                                                        | 13位unix时间戳                                        |

**Response:**

| Key      | Value               | Remark                                                |
| -------- | ------------------- | ----------------------------------------------------- |
| retcode  | 0<br/>1101          | 正常<br/>已退出登陆                                   |
| selector | 0<br/>2<br/>4<br/>7 | 正常<br/>有新消息<br/>联系人有更新<br/>手机点击联系人 |

*Example:*
``` javascript
window.synccheck={retcode:"0",selector:"2"}
```


## 发送

### 消息已读

将指定联系人的消息设为已读

| Key         | Value                                                    |
| ----------- | -------------------------------------------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify |
| Method      | POST                                                     |
| Cookie      | Need                                                     |
| PARAM       | pass_ticket                                              |

**Body (json):**

``` json
{"BaseRequest":
    {
        "Uin":216547950,
        "Sid":"sUCEDniMkkAj6Wxh",
        "Skey":"@crypt_a6a25b34_525e3d8eb6358a022dffa7ad6f3806b2",
        "DeviceID":"e575465866227354"
    },
    "Code":1,
    "FromUserName":"@ae820581771d028ec2540c22b57ee02289811811caa08ec8d88e7cdb0f04502e", // 自己的用户id
    "ToUserName":"@4f460c580a3798ed6ed571593f694f72", // 目标用户id
    "ClientMsgId":1508393667441 // 直接使用13位时间戳填充即可
}
```

---

### 发送消息

| Key         | Value                                               |
| ----------- | --------------------------------------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg |
| Method      | POST                                                |
| Cookie      | Need                                                |
| Param       | pass_ticket                                         |

#### Body (Json):

``` json
{
    "BaseRequest": {
        "Uin": 216547950,
        "Sid": "Hga/ND66ty7ptu4f",
        "Skey": "@crypt_a6a25b34_2df15eec5697a324849770ce822e2b67",
        "DeviceID": "e680668306876822"
    },
    "Msg": {
        "Type": 1,
        "Content": "hhh",
        "FromUserName": "@be6d9f4847c79706435ca6bd55aa2f673848851278cc5b0001c49720ee9c3e04",
        "ToUserName": "@77d902e96d228e4eb17ee4f02c6e12ce",
        "LocalID": "14908900889660665",
        "ClientMsgId": "14908900889660665"
    },
    "Scene": 0
}
```

---

### 撤回消息

| Key         | Value                                                 |
| ----------- | ----------------------------------------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxrevokemsg |
| Method      | POST                                                  |
| Cookie      | Need                                                  |

**Body (json):**

``` json
{
    "BaseRequest": {
        "Uin": 216547950,
        "Sid": "1R+tNd8DPJkULlu+",
        "Skey": "@crypt_a6a25b34_cb08394069b5dba8d090c48ea84849a8",
        "DeviceID": "e935523984078190"
    },
    "SvrMsgId": "5918499768689813400",
    "ToUserName": "@8d0cb0307ce18d0c8c51dd788060bf56",
    "ClientMsgId": "14908762377750838"
}
```

### 转发被撤回的图片消息

要转发被撤回的图片消息，只需将撤回的图片消息的Content中的aeskey(cdnthumbaeskey)、cdnthumburl(cdnmidimgurl)、md5复制到发送图片的对应字段中即可

---

### 修改用户备注

| Key         | Value                                                 |
| ----------- | ----------------------------------------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxoplog |
| Method      | POST                                                  |
| Cookie      | Need                                                  |

**Body (json):**
``` json
{
    "UserName": "@f5fc0106e419ed58baafd50a9d4b4f4869d417411cd834ffe4f43ac62bbc38a6",
    "CmdId": 2,
    "RemarkName": "3123",
    "BaseRequest": {
        "Uin": 216547950,
        "Sid": "2ZgAp8arXklWje6v",
        "Skey": "@crypt_a6a25b34_bdf4052fc37832fabf2c2dc5421f8fca",
        "DeviceID": "e968379522819118"
    }
}
```

### 修改群名

| Key         | Value                                                 |
| ----------- | ----------------------------------------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxupdatechatroom?fun=modtopic |
| Method      | POST                                                  |
| Cookie      | Need                                                  |

**Body (json):**
``` json
{
    "NewTopic": "KKK1",
    "ChatRoomName": "@@5c6af4a5215187e41484226127ddfd19646559ec4586f0c2a7f6e048d5e8cb98",
    "BaseRequest": {
        "Uin": 216547950,
        "Sid": "2ZgAp8arXklWje6v",
        "Skey": "@crypt_a6a25b34_bdf4052fc37832fabf2c2dc5421f8fca",
        "DeviceID": "e786480153373391"
    }
}
```
