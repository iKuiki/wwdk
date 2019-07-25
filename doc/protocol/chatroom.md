# 群管理

- [群管理](#%e7%be%a4%e7%ae%a1%e7%90%86)
  - [好友拉群](#%e5%a5%bd%e5%8f%8b%e6%8b%89%e7%be%a4)
  - [修改群名](#%e4%bf%ae%e6%94%b9%e7%be%a4%e5%90%8d)

## 好友拉群

选定好友创建一个新的群聊

**Request:**

| Key         | Value                                                           | Remark                             |
| ----------- | --------------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxcreatechatroom> |                                    |
| Method      | POST                                                            |                                    |
| Param       | r                                                               |                                    |
| Param       | pass_ticket                                                     | 部分Domain需要传，保险起见可以都传 |

**Body:**

Body里MemberList为要添加进入群聊的联系人的UserName

``` json
{
    "MemberCount": 2,
    "MemberList": [{
        "UserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    }, {
        "UserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    }],
    "Topic": "", // 此处可以填入群标题？不过网页版没有提供这个界面，不建议使用
    "BaseRequest": {
        "Uin": 1800000000,
        "Sid": "xxxxxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e920000000000"
    }
}
```

**Response:**

返回为json对象，内包含了新的群的UserName，务必记录

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": "Everything is OK"
    },
    "Topic": "",
    "PYInitial": "",
    "QuanPin": "",
    "MemberCount": 2,
    "MemberList": [{
        "Uin": 0,
        "UserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "NickName": "xxxx",
        "AttrStatus": 0,
        "PYInitial": "",
        "PYQuanPin": "",
        "RemarkPYInitial": "",
        "RemarkPYQuanPin": "",
        "MemberStatus": 0,
        "DisplayName": "",
        "KeyWord": ""
    }, {
        "Uin": 0,
        "UserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "NickName": "xxxx",
        "AttrStatus": 0,
        "PYInitial": "",
        "PYQuanPin": "",
        "RemarkPYInitial": "",
        "RemarkPYQuanPin": "",
        "MemberStatus": 0,
        "DisplayName": "",
        "KeyWord": ""
    }],
    "ChatRoomName": "@@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 这个是关键，这个是新的群的UserName，通过这个请求新的群的完整群成员
    "BlackList": ""
}
```

## 修改群名

| Key         | Value                                                                        | Remark |
| ----------- | ---------------------------------------------------------------------------- | ------ |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxupdatechatroom?fun=modtopic> |        |
| Method      | POST                                                                         |        |

**Body (json):**

``` json
{
    "NewTopic": "new topic name", // 新群名
    "ChatRoomName": "@@xxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 群的username
    "BaseRequest": {
        "Uin": 200000000,
        "Sid": "xxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e78000000000000"
    }
}
```
