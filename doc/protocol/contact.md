# 联系人管理

- [联系人管理](#%e8%81%94%e7%b3%bb%e4%ba%ba%e7%ae%a1%e7%90%86)
  - [获取联系人列表](#%e8%8e%b7%e5%8f%96%e8%81%94%e7%b3%bb%e4%ba%ba%e5%88%97%e8%a1%a8)
  - [批量获取联系人](#%e6%89%b9%e9%87%8f%e8%8e%b7%e5%8f%96%e8%81%94%e7%b3%bb%e4%ba%ba)
  - [修改用户备注](#%e4%bf%ae%e6%94%b9%e7%94%a8%e6%88%b7%e5%a4%87%e6%b3%a8)
  - [处理添加好友的消息](#%e5%a4%84%e7%90%86%e6%b7%bb%e5%8a%a0%e5%a5%bd%e5%8f%8b%e7%9a%84%e6%b6%88%e6%81%af)
    - [接到请求添加好友请求](#%e6%8e%a5%e5%88%b0%e8%af%b7%e6%b1%82%e6%b7%bb%e5%8a%a0%e5%a5%bd%e5%8f%8b%e8%af%b7%e6%b1%82)
    - [同意添加](#%e5%90%8c%e6%84%8f%e6%b7%bb%e5%8a%a0)

## 获取联系人列表

注:*此接口获取到的联系人还不是完整联系人，要与之前init时获取到的联系人合并才是完整的联系人列表*
注:*此接口获取到的联系人中，群聊不包含群成员，需要在调用getBatchContact接口获取成员*

| Key         | Value                                                       | Remark                             |
| ----------- | ----------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxgetcontact> |                                    |
| Method      | Post                                                        |                                    |
| Param       | r                                                           | 13位时间戳                         |
| Param       | seq                                                         | 填0                                |
| Param       | skey                                                        | 获取登陆参数时获取到的             |
| Param       | pass_ticket                                                 | 部分Domain需要传，保险起见可以都传 |

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

## 批量获取联系人

此方法用于批量获取联系人，它有2个用途：

- 获取群聊的成员列表
- 获取群成员的详细信息

在官方web微信中，这两种用法是分开的，所以建议分开调用
在获取群聊成员列表时，获取到的成员信息是简略信息，只有相当简陋的内容（几乎只有名字）
而获取群成员的详细信息是，要将目标成员的username作为UserName字段，群username作为EncryChatRoomId字段写入body的json的List中，这样返回的群成员详细信息就与好友的返回信息近似了

注:*群判定方法为username开头为@@*

| Key         | Value                                                            | Remark                                        |
| ----------- | ---------------------------------------------------------------- | --------------------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxbatchgetcontact> |                                               |
| Method      | Post                                                             |                                               |
| Param       | type                                                             | 填ex                                          |
| Param       | r                                                                | 13位时间戳                                    |
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
    "List": [ // 此List中填入要获取详细信息的目标对象，目标可以是群也可以是群成员
        {
            // 如果是获取群的成员列表，此处填群UserName
            // 如果是获取群成员的详情，此处填群成员的UserName
            "UserName": "@@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
            // 如果是获取群的成员列表，此处留空
            // 如果是获取群成员的详情，此处填群UserName
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

---

## 处理添加好友的消息

web微信可以处理添加好友请求，当有人添加好友时，可以直接通过web微信接口同意

### 接到请求添加好友请求

当收到好友请求时，webwxsync接口会接到一条由UserName为**fmessage**发来的消息，其消息本体如下(即AddMsgList中的对象，已隐去webwxsync外部无关的结构)

``` json
{
    "MsgId": "1510000000000000000",
    "FromUserName": "fmessage",
    "ToUserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxx",
    "MsgType": 37,
    "Content": "&lt;msg fromusername=\"wx_id\" encryptusername=\"v1_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@stranger\" fromnickname=\"xxxx\" content=\"我是xxxxxxxxxx\"  shortpy=\"xxxx\" imagestatus=\"3\" scene=\"14\" country=\"\" province=\"\" city=\"\" sign=\"\" percard=\"0\" sex=\"0\" alias=\"\" weibo=\"\" albumflag=\"0\" albumstyle=\"0\" albumbgimgid=\"\" snsflag=\"1\" snsbgimgid=\"\" snsbgobjectid=\"0\" mhash=\"xxxxxxxxxxxxxxxx\" mfullhash=\"xxxxxxxxxxxxxxxx\" bigheadimgurl=\"http://wx.qlogo.cn/mmhead/ver_1/xxxxxxxxx/0\" smallheadimgurl=\"http://wx.qlogo.cn/mmhead/ver_1/xxxxxxxxx/96\" ticket=\"v2_xxxxxxxxx@stranger\" opcode=\"2\" googlecontact=\"\" qrticket=\"\" chatroomusername=\"1000000@chatroom\" sourceusername=\"\" sourcenickname=\"\"&gt;&lt;brandlist count=\"0\" ver=\"680000000\"&gt;&lt;/brandlist&gt;&lt;/msg&gt;", // 进行html Unescape后，可以得到一个xml对象
    "Status": 3,
    "ImgStatus": 1,
    "CreateTime": 1564046056,
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
        "UserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxx",
        "NickName": "xxxx",
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "我是xxxx",
        "Signature": "",
        "Alias": "",
        "Scene": 14, // 未知参数
        "VerifyFlag": 0,
        "AttrStatus": 233765, // 未知参数
        "Sex": 0,
        "Ticket": "v2_xxxxxxxxx@stranger",
        "OpCode": 2
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
    "NewMsgId": 1512190743315786023,
    "OriContent": "",
    "EncryFileName": ""
}
```

其上Content中的字符串进行html Unescape后可得到如下的xml对象
注：*其实这个xml对api好像没啥用，里面比较有用的就是头像了，别的数据好像暂时用不到*

``` xml
<msg
    fromusername="wx_id"
    encryptusername="v1_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@stranger"
    fromnickname="xxxx"
    content="我是xxxxxxxxxx"
    shortpy="xxxx"
    imagestatus="3"
    scene="14"
    country=""
    province=""
    city=""
    sign=""
    percard="0"
    sex="0"
    alias=""
    weibo=""
    albumflag="0"
    albumstyle="0"
    albumbgimgid=""
    snsflag="1"
    snsbgimgid=""
    snsbgobjectid="0"
    mhash="xxxxxxxxxxxxxxxx"
    mfullhash="xxxxxxxxxxxxxxxx"
    bigheadimgurl="http://wx.qlogo.cn/mmhead/ver_1/xxxxxxxxx/0"
    smallheadimgurl="http://wx.qlogo.cn/mmhead/ver_1/xxxxxxxxx/96"
    ticket="v2_xxxxxxxxx@stranger"
    opcode="2"
    googlecontact=""
    qrticket=""
    chatroomusername="1000000@chatroom"
    sourceusername=""
    sourcenickname="">
    <brandlist count="0"
        ver="680000000"></brandlist>
</msg>
```

### 同意添加

接到以上请求后，调用如下接口即可接受添加好友请求

| Key         | Value                                                       | Remark     |
| ----------- | ----------------------------------------------------------- | ---------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxverifyuser> |            |
| Method      | Post                                                        |            |
| Param       | r                                                           | 13位时间戳 |

**Body:**

body是一个json对象,构建时需要用到刚刚收到得添加好友请求中得参数

``` json
{
    "BaseRequest": {
        "Uin": 200000000,
        "Sid": "xxxxxxxxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e960000000000000"
    },
    "Opcode": 3, // 固定值，照填
    "VerifyUserListSize": 1, // 网页版一次操作也只能同意一个请求，所以别作死在调用接口得时候一次同意多个
    "VerifyUserList": [{
        "Value": "@xxxxxxxxxxxxxxxxxxxxxxxxxx", // 刚刚获取到得对方得UserName
        "VerifyUserTicket": "v2_xxxxxxxxx@stranger" // 刚刚获取到得Msg中RecommendInfo里的Ticket
    }],
    "VerifyContent": "",
    "SceneListCount": 1, // 固定值，照填
    "SceneList": [33], // 固定值，照填
    "skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxx" // 登陆凭据Skey，和BaseRequest的Skey是同一个东西
}
```

**Response:**

发送完请求后，如无意外会接到如下Json返回

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    }
}
```
