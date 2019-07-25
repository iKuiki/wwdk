# 联系人管理

- [联系人管理](#%e8%81%94%e7%b3%bb%e4%ba%ba%e7%ae%a1%e7%90%86)
  - [获取联系人列表](#%e8%8e%b7%e5%8f%96%e8%81%94%e7%b3%bb%e4%ba%ba%e5%88%97%e8%a1%a8)
  - [获取群成员列表](#%e8%8e%b7%e5%8f%96%e7%be%a4%e6%88%90%e5%91%98%e5%88%97%e8%a1%a8)
  - [修改用户备注](#%e4%bf%ae%e6%94%b9%e7%94%a8%e6%88%b7%e5%a4%87%e6%b3%a8)

## 获取联系人列表

注:*此接口获取到的联系人还不是完整联系人，要与之前init时获取到的联系人合并才是完整的联系人列表*
注:*此接口获取到的联系人中，群聊不包含群成员，需要在调用getBatchContact接口获取成员*

| Key         | Value                                                       | Remark                             |
| ----------- | ----------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxgetcontact> |                                    |
| Method      | Post                                                        |                                    |
| Param       | r                                                           | 13位时间戳取反                     |
| Param       | seq                                                         | 填0                                |
| Param       | skey                                                        | 获取登陆参数时获取到的             |
| PARAM       | pass_ticket                                                 | 部分Domain需要传，保险起见可以都传 |

**Response:**

返回值为一个json信息，里面包含绝大部分联系人的信息，于登陆时init获取到的联系人合并，就是完整的联系人列表

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    },
    "MemberCount": 1,
    "MemberList": [
        {
            "Uin": 0,
            "UserName": "weixin",
            "NickName": "微信团队",
            "HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgeticon?xxx=xxx&xxx=xxx",
            "ContactFlag": 3,
            "MemberCount": 0,
            "MemberList": [],
            "RemarkName": "",
            "HideInputBarFlag": 0,
            "Sex": 0,
            "Signature": "微信团队官方帐号",
            "VerifyFlag": 56,
            "OwnerUin": 0,
            "PYInitial": "WXTD",
            "PYQuanPin": "weixintuandui",
            "RemarkPYInitial": "",
            "RemarkPYQuanPin": "",
            "StarFriend": 0,
            "AppAccountFlag": 0,
            "Statues": 0,
            "AttrStatus": 4,
            "Province": "",
            "City": "",
            "Alias": "",
            "SnsFlag": 0,
            "UniFriend": 0,
            "DisplayName": "",
            "ChatRoomId": 0,
            "KeyWord": "wei",
            "EncryChatRoomId": "",
            "IsOwner": 0
        }],
    "Seq": 0
}
```

---

## 获取群成员列表

此方法可以获取指定群的群成员列表，可以一次获取多个群，所以直接将所有群的username一起打包请求，就可以一次拿到所有群的群成员列表了
群username列表要用json打包在请求的body中

注:*群判定方法为username开头为@@*

| Key         | Value                                                            | Remark                                        |
| ----------- | ---------------------------------------------------------------- | --------------------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxbatchgetcontact> |                                               |
| Method      | Post                                                             |                                               |
| Param       | type                                                             | 填ex                                          |
| Param       | r                                                                | 13位时间戳取反                                |
| Param       | pass_ticket                                                      | 获取登陆参数时获取到的,好像只有v2版本需要传？ |

**Request Body:**

```json
{
    "BaseRequest": {
        "Uin": 1880000000,
        "Sid": "xxxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxx",
        "DeviceID": "e60000000000000000"
    },
    "Count": 1,
    "List": [
        {
            "UserName": "@@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
            "EncryChatRoomId": ""
        }
    ]
}
```

**Response:**

返回为一个json数据，内容为请求的群的详细数据，如果请求的不是群，则会直接返回对应的原始联系人

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    },
    "Count": 1,
    "ContactList": [
        {
            "Uin": 0,
            "UserName": "@@xxxxxxxxxxxxxxxxxxxxxx",
            "NickName": "xxxxxxxx",
            "HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgetheadimg?xxx=xxx&xxx=xxx",
            "ContactFlag": 2,
            "MemberCount": 1,
            "MemberList": [
                {
                    "Uin": 0,
                    "UserName": "@xxxxxxxxxxxxxxxxxxxxxx",
                    "NickName": "xxxxxx",
                    "AttrStatus": 4000000,
                    "PYInitial": "",
                    "PYQuanPin": "",
                    "RemarkPYInitial": "",
                    "RemarkPYQuanPin": "",
                    "MemberStatus": 0,
                    "DisplayName": "",
                    "KeyWord": ""
                }
            ],
            "RemarkName": "",
            "HideInputBarFlag": 0,
            "Sex": 0,
            "Signature": "",
            "VerifyFlag": 0,
            "OwnerUin": 0,
            "PYInitial": "xxxxxxx",
            "PYQuanPin": "xxxxxxxxxxxxxxxxxxxxxx",
            "RemarkPYInitial": "",
            "RemarkPYQuanPin": "",
            "StarFriend": 0,
            "AppAccountFlag": 0,
            "Statues": 1,
            "AttrStatus": 0,
            "Province": "",
            "City": "",
            "Alias": "",
            "SnsFlag": 0,
            "UniFriend": 0,
            "DisplayName": "",
            "ChatRoomId": 0,
            "KeyWord": "",
            "EncryChatRoomId": "@xxxxxxxxxxxxxxxxxxxxxx",
            "IsOwner": 0
        }
    ]
}
```

---

## 修改用户备注

| Key         | Value                                                  |
| ----------- | ------------------------------------------------------ |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxoplog> |
| Method      | POST                                                   |
| Cookie      | Need                                                   |

**Body (json):**

``` json
{
    "UserName": "@xxxxxxxxxxxxxxxxxxxxxx",
    "CmdId": 2, // 固定值取2
    "RemarkName": "new remark Name",
    "BaseRequest": {
        "Uin": 200000000,
        "Sid": "xxxxxxxxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e960000000000000"
    }
}
```
