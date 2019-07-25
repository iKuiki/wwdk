# 接收

接受消息主要分为2个步骤

- 轮询微信状态同步接口，获得当前微信状态
- 如果当前微信状态为有新消息，则调用获取新消息接口取得新消息

- [接收](#%e6%8e%a5%e6%94%b6)
  - [同步状态](#%e5%90%8c%e6%ad%a5%e7%8a%b6%e6%80%81)
  - [拉取新消息](#%e6%8b%89%e5%8f%96%e6%96%b0%e6%b6%88%e6%81%af)
    - [Request](#request)
    - [Response](#response)

---

## 同步状态

注：**同步状态接口中有个synckey，这个synckey是在登陆时的初始化(webwxinit)操作中获取的，并且每次调用getMessage接口都会刷新这个synckey**

| Key         | Value                                                         | Remark                                                |
| ----------- | ------------------------------------------------------------- | ----------------------------------------------------- |
| Request URL | <https://webpush.{{apiDomain}}/cgi-bin/mmwebwx-bin/synccheck> |                                                       |
| Method      | Get                                                           |                                                       |
| Param       | r                                                             | 13位unix时间戳                                        |
| Param       | skey                                                          | 登陆凭据skey                                          |
| Param       | sid                                                           | 登陆凭据sid                                           |
| Param       | uin                                                           | 登陆凭据uin                                           |
| Param       | deviceid                                                      |                                                       |
| Param       | synckey                                                       | 同步key，由"\|"分割为1_xxx\|2_xxx\|3_xxx\|4_xxx的格式 |
| Param       | _                                                             | 13位unix时间戳                                        |

**Response:**

返回为一串类似js赋值代码的语句，从中解析出retcode为0即可使用
selector代表状态，为2时为有新消息（还需再推敲

*Example:*

``` javascript
window.synccheck={retcode:"0",selector:"2"} // selector的值的含义还需要再推敲
```

| Key      | Value               | Remark                                                |
| -------- | ------------------- | ----------------------------------------------------- |
| retcode  | 0<br/>1101          | 正常<br/>已退出登陆                                   |
| selector | 0<br/>2<br/>4<br/>7 | 正常<br/>有新消息<br/>联系人有更新<br/>手机点击联系人 |

## 拉取新消息

当synccheck检测到有新消息时（可以是信息或者联系人变动等），就需要拉取新消息来展示给用户
拉取新消息的接口中会返回新的SyncKey，务必记得更新

### Request

| Key         | Value                                                 | Remark                             |
| ----------- | ----------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxsync> |                                    |
| Method      | Post                                                  |                                    |
| Param       | sid                                                   | 登陆凭据sid                        |
| Param       | skey                                                  | 登陆凭据skey                       |
| Param       | pass_ticket                                           | 部分Domain需要传，保险起见可以都传 |

**Body:**

```json
{
    "BaseRequest": {
        "Uin":"210000000",
        "Sid":"QQxxxxxxxxxxxxxx",
        "Skey":"@crypt_a6xxxxxx_6xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID":"e980000000000000"
    },
    "SyncKey": {
        "Count": 4,
        "List": [{
            "Key": 1,
            "Val": 600000000
        }, {
            "Key": 2,
            "Val": 600000000
        }, {
            "Key": 3,
            "Val": 600000000
        }, {
            "Key": 1000,
            "Val": 600000000
        }]
    },
    "rr": -600000000 // 13位unix时间戳取反
}
```

### Response

返回为一个json对象，其中包含新消息、联系人变化、联系人删除、SyncKey这些关键信息
其中SyncKey必须更新到本地，在syncCheck接口中必须使用最新的SyncKey

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    },
    "AddMsgCount": 1,
    "AddMsgList": [
        {
            "MsgId": "8800000000000000000",
            // @@开头的UserName是群聊
            "FromUserName": "@@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
            "ToUserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
            "MsgType": 1,
            // 如果是群消息，则此处会先是发言的群成员的UserName然后紧跟冒号与html的换行，接下来才是内容
            "Content": "@xxxxxxxxxxxxxxxxxxxxx:<br/>Test Message",
            "Status": 3,
            "ImgStatus": 1,
            "CreateTime": 1560000000,
            "VoiceLength": 0,
            "PlayLength": 0,
            "FileName": "",
            "FileSize": "",
            "MediaId": "",
            "Url": "",
            "AppMsgType": 0,
            "StatusNotifyCode": 0,
            "StatusNotifyUserName": "",
            "RecommendInfo": {
                "UserName": "",
                "NickName": "",
                "QQNum": 0,
                "Province": "",
                "City": "",
                "Content": "",
                "Signature": "",
                "Alias": "",
                "Scene": 0,
                "VerifyFlag": 0,
                "AttrStatus": 0,
                "Sex": 0,
                "Ticket": "",
                "OpCode": 0
            },
            "ForwardFlag": 0,
            "AppInfo": {
                "AppID": "",
                "Type": 0
            },
            "HasProductId": 0,
            "Ticket": "",
            "ImgHeight": 0,
            "ImgWidth": 0,
            "SubMsgType": 0,
            "NewMsgId": 8800000000000000000,
            "OriContent": "",
            "EncryFileName": ""
        }
    ],
    "ModContactCount": 1,
    "ModContactList": [ // 如果有联系人改动，则会出现在此
        {
            "UserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
            "NickName": "xxxxxxx",
            "Sex": 0,
            "HeadImgUpdateFlag": 1,
            "ContactType": 0,
            "Alias": "",
            "ChatRoomOwner": "",
            "HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgeticon?xxx=xxx&xxx=xxx",
            "ContactFlag": 3,
            "MemberCount": 0,
            "MemberList": [],
            "HideInputBarFlag": 0,
            "Signature": "",
            "VerifyFlag": 0,
            "RemarkName": "",
            "Statues": 0,
            "AttrStatus": 233765,
            "Province": "",
            "City": "",
            "SnsFlag": 1,
            "KeyWord": ""
        }
    ],
    "DelContactCount": 1,
    "DelContactList": [ // 如果有删除联系人，则会在此列表更新
        {
            "UserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
            "ContactFlag": 0
        }
    ],
    "ModChatRoomMemberCount": 0,
    "ModChatRoomMemberList": [], // 为见此字段有变化过，群成员变化会直接在ModContactList处反应
    "Profile": {
        "BitFlag": 0,
        "UserName": {
            "Buff": ""
        },
        "NickName": {
            "Buff": ""
        },
        "BindUin": 0,
        "BindEmail": {
            "Buff": ""
        },
        "BindMobile": {
            "Buff": ""
        },
        "Status": 0,
        "Sex": 0,
        "PersonalCard": 0,
        "Alias": "",
        "HeadImgUpdateFlag": 0,
        "HeadImgUrl": "",
        "Signature": ""
    },
    "ContinueFlag": 0,
    "SyncKey": {
        "Count": 2,
        "List": [
            {
                "Key": 1,
                "Val": 600000000
            },
            {
                "Key": 2,
                "Val": 600000000
            }
        ]
    },
    "SKey": "",
    "SyncCheckKey": {
        "Count": 2,
        "List": [
            {
                "Key": 1,
                "Val": 600000000
            },
            {
                "Key": 2,
                "Val": 600000000
            }
        ]
    }
}
```
