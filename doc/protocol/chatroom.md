# 群管理

- [群管理](#%e7%be%a4%e7%ae%a1%e7%90%86)
  - [好友拉群](#%e5%a5%bd%e5%8f%8b%e6%8b%89%e7%be%a4)
  - [添加新成员](#%e6%b7%bb%e5%8a%a0%e6%96%b0%e6%88%90%e5%91%98)
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

---

## 添加新成员

在已有的群中添加新成员

| Key         | Value                                                           | Remark                             |
| ----------- | --------------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxupdatechatroom> |                                    |
| Method      | POST                                                            |                                    |
| Param       | fun                                                             | 填addmember                        |
| Param       | pass_ticket                                                     | 部分Domain需要传，保险起见可以都传 |

**Body:**

``` json
{
    "AddMemberList": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 要添加的新成员的UserName
    "ChatRoomName": "@@xxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 要添加的目标群的UserName
    "BaseRequest": {
        "Uin": 200000000,
        "Sid": "xxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e78000000000000"
    }
}
```

**Response:**

返回为json对象，只要判断未出错即可，无其他重要内容(新成员因为本来就是自己邀请的，所以UserName肯定知道啦)
在返回无错误以后就应当把成员添加进对应聊天室的Menber中，虽然推送也会通知联系人变更，不过那个太慢了，据我最近一次测试无新消息的话延迟有3分钟

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    },
    "MemberCount": 1,
    "MemberList": [{
        "Uin": 0,
        "UserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 被添加的联系人的UserName
        "NickName": "xxxx", // 被添加的联系人的昵称
        "AttrStatus": 0,
        "PYInitial": "",
        "PYQuanPin": "",
        "RemarkPYInitial": "",
        "RemarkPYQuanPin": "",
        "MemberStatus": 0,
        "DisplayName": "",
        "KeyWord": ""
    }]
}
```

---

## 修改群名

| Key         | Value                                                           | Remark                             |
| ----------- | --------------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxupdatechatroom> |                                    |
| Method      | POST                                                            |                                    |
| Param       | fun                                                             | 填modtopic                         |
| Param       | pass_ticket                                                     | 部分Domain需要传，保险起见可以都传 |

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

**Response:**

返回为json对象，只要判断未出错即可，无其他重要内容

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    },
    "MemberCount": 0,
    "MemberList": []
}
```
