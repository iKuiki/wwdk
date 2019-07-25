# 登陆web微信

登陆流程大致为以下步骤：

- 获取uuid
- 根据uuid生成二维码
- 轮询微信服务器检测用户是否扫码
- 手机扫码
- 获取到用户已扫码，得到用户头像（base64编码的字符串）
- 用户同意登陆
- 登陆成功

- [登陆web微信](#%e7%99%bb%e9%99%86web%e5%be%ae%e4%bf%a1)
  - [获取uuid](#%e8%8e%b7%e5%8f%96uuid)
  - [轮询用户扫码](#%e8%bd%ae%e8%af%a2%e7%94%a8%e6%88%b7%e6%89%ab%e7%a0%81)
  - [获取登陆参数](#%e8%8e%b7%e5%8f%96%e7%99%bb%e9%99%86%e5%8f%82%e6%95%b0)
  - [初始化](#%e5%88%9d%e5%a7%8b%e5%8c%96)

---

## 获取uuid

| Key         | Value                                 | Remark               |
| ----------- | ------------------------------------- | -------------------- |
| Request URL | <https://login.{{apiDomain}}/jslogin> |                      |
| Method      | Get                                   |                      |
| Param       | appid                                 | 填wx782c26e4c19acffb |
| Param       | fun                                   | 填new                |
| Param       | lang                                  | zh_CN或en_US         |
| Param       | _                                     | 13位unix时间戳       |

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

## 轮询用户扫码

| Key         | Value                                                   | Remark                     |
| ----------- | ------------------------------------------------------- | -------------------------- |
| Request URL | <https://login.{{apiDomain}}/cgi-bin/mmwebwx-bin/login> |                            |
| Method      | Get                                                     |                            |
| Param       | loginicon                                               | 填true                     |
| Param       | uuid                                                    | 之前获取的uuid             |
| Param       | tip                                                     | 填0                        |
| Param       | r                                                       | 13位时间戳取反(貌似可省略) |
| Param       | _                                                       | 13位unix时间戳             |

**Response:**

| Key                 | Value                       | Remark                                                     |
| ------------------- | --------------------------- | ---------------------------------------------------------- |
| window.code         | 200<br/>201<br/>400<br/>408 | 确认登陆<br/>已扫码<br/>登陆超时(二维码失效)<br/>等待登陆  |
| window.userAvatar   | data:img/jpeg;base64        | base64编码的用户头像，仅当code=201时才有                   |
| window.redirect_uri | redirect_uri                | 下一跳地址，获取到这个地址以后访问这个地址获取下一步的信息 |

*Example:*

``` javascript
window.code = 408;

// Or

window.code=201;window.userAvatar='data:img/jpeg;base64,iVBORw...'

// Or

window.code = 200;window.redirect_uri = "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=Aeuxxxxxxxxxxxxxxxxxxxx@qrticket_0&uuid=Yaxxx-xxxx==&lang=zh_&scan=1560000000";
```

注：**若用户取消登陆，返回仍为408，旧的二维码仍可重复使用，用户重新扫旧的二维码后会再次返回201**

注：*当登陆超时（二维码失效）后，重新调用获取uuid的方法即可重新拿到二维码*

---

## 获取登陆参数

此处其实就是在上一步获取到redirect_uri后在其后拼接```&fun=new```即可

注：**有部分domain版本在此处还需要添加一个```&version=v2```的参数，需要考证通用性的问题**

| Key         | Value                                                         | Remark                 |
| ----------- | ------------------------------------------------------------- | ---------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxnewloginpage> |                        |
| Method      | Get                                                           |                        |
| Param       | ticket                                                        |                        |
| Param       | uuid                                                          |                        |
| Param       | scan                                                          | 扫描成功后返回的时间戳 |
| Param       | fun                                                           | 填new                  |
| Param       | version                                                       | 填v2                   |

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

## 初始化

| Key         | Value                                                 | Remark                             |
| ----------- | ----------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxinit> |                                    |
| Method      | Post                                                  |                                    |
| Param       | r                                                     | 13位时间戳取反                     |
| Param       | pass_ticket                                           | 部分Domain需要传，保险起见可以都传 |

*Body (json):*

``` json
{"BaseRequest":
    {
        "Uin":"210000000",
        "Sid":"QQxxxxxxxxxxxxxx",
        "Skey":"@crypt_a6xxxxxx_6xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID":"e980000000000000"
    }
}
```

**Response:**

返回为一个json对象，内包括用户信息、联系人(此列表可能不全，如果不全，之后用获取联系人的接口获取联系人列表合并后即为完整的)、同步信息等
*该json中主要需要用到的数据为User信息、联系人列表与synckey，synckey用于后续同步状态时使用*

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    },
    "Count": 1,
    "ContactList": [{
        "Uin": 0,
        "UserName": "filehelper",
        "NickName": "文件传输助手",
        "HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgeticon?xxx=xxx&xxx=xxx",
        "ContactFlag": 3,
        "MemberCount": 0,
        "MemberList": [],
        "RemarkName": "",
        "HideInputBarFlag": 0,
        "Sex": 0,
        "Signature": "",
        "VerifyFlag": 0,
        "OwnerUin": 0,
        "PYInitial": "WJCSZS",
        "PYQuanPin": "wenjianchuanshuzhushou",
        "RemarkPYInitial": "",
        "RemarkPYQuanPin": "",
        "StarFriend": 0,
        "AppAccountFlag": 0,
        "Statues": 0,
        "AttrStatus": 0,
        "Province": "",
        "City": "",
        "Alias": "",
        "SnsFlag": 0,
        "UniFriend": 0,
        "DisplayName": "",
        "ChatRoomId": 0,
        "KeyWord": "fil",
        "EncryChatRoomId": "",
        "IsOwner": 0
    }],
    "SyncKey": {
        "Count": 4,
        "List": [{
            "Key": 1,
            "Val": 100000000
        }, {
            "Key": 2,
            "Val": 100000000
        }, {
            "Key": 3,
            "Val": 100000000
        }, {
            "Key": 1000,
            "Val": 100000000
        }]
    },
    "User": {
        "Uin": 1882092380,
        "UserName": "@xxxxxxxxx",
        "NickName": "xxxxxx",
        "HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgeticon?xxx=xxx&xxx=xxx",
        "RemarkName": "",
        "PYInitial": "",
        "PYQuanPin": "",
        "RemarkPYInitial": "",
        "RemarkPYQuanPin": "",
        "HideInputBarFlag": 0,
        "StarFriend": 0,
        "Sex": 1,
        "Signature": "",
        "AppAccountFlag": 0,
        "VerifyFlag": 0,
        "ContactFlag": 0,
        "WebWxPluginSwitch": 2,
        "HeadImgFlag": 1,
        "SnsFlag": 177
    },
    "ChatSet": "filehelper,@xxxx,@@xxxxx,",
    "SKey": "@crypt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "ClientVersion": 385877285,
    "SystemTime": 1560000000,
    "GrayScale": 1,
    "InviteStartCount": 40,
    "MPSubscribeMsgCount": 1,
    "MPSubscribeMsgList": [{
        "UserName": "@xxxx",
        "MPArticleCount": 1,
        "MPArticleList": [{
            "Title": "xxxx",
            "Digest": "xxxx",
            "Cover": "https://xxx.xxx.xxx/xxx/xxx/640?wxtype=jpeg&wxfrom=0",
            "Url": "http://mp.weixin.qq.com/xxxxxxxxxxx"
        }],
        "Time": 1560000000,
        "NickName": "xxx"
    }],
    "ClickReportInterval": 600000
}
```

到此登陆就成功了
