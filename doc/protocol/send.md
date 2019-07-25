
# 发送

- [发送](#%e5%8f%91%e9%80%81)
  - [消息已读](#%e6%b6%88%e6%81%af%e5%b7%b2%e8%af%bb)
  - [发送消息](#%e5%8f%91%e9%80%81%e6%b6%88%e6%81%af)
  - [撤回消息](#%e6%92%a4%e5%9b%9e%e6%b6%88%e6%81%af)
  - [转发被撤回的图片消息](#%e8%bd%ac%e5%8f%91%e8%a2%ab%e6%92%a4%e5%9b%9e%e7%9a%84%e5%9b%be%e7%89%87%e6%b6%88%e6%81%af)

## 消息已读

将指定联系人的消息设为已读

| Key         | Value                                                         | Remark                             |
| ----------- | ------------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxstatusnotify> |                                    |
| Method      | POST                                                          |                                    |
| PARAM       | pass_ticket                                                   | 部分Domain需要传，保险起见可以都传 |

**Body (json):**

``` json
{"BaseRequest":
    {
        "Uin":"210000000",
        "Sid":"QQxxxxxxxxxxxxxx",
        "Skey":"@crypt_a6xxxxxx_6xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID":"e980000000000000"
    },
    "Code":1,
    "FromUserName":"@ae820581771d028ec2540c22b57ee02289811811caa08ec8d88e7cdb0f04502e", // 自己的用户id
    "ToUserName":"@4f460c580a3798ed6ed571593f694f72", // 目标用户id
    "ClientMsgId":1508393667441 // 直接使用13位时间戳填充即可
}
```

---

## 发送消息

| Key         | Value                                                    | Remark                             |
| ----------- | -------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxsendmsg> |                                    |
| Method      | POST                                                     |                                    |
| Cookie      | Need                                                     |                                    |
| Param       | pass_ticket                                              | 部分Domain需要传，保险起见可以都传 |

**Body (Json):**

``` json
{
    "BaseRequest": {
        "Uin": 200000000,
        "Sid": "xxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e680000000000000"
    },
    "Msg": {
        "Type": 1,
        "Content": "hello world", // 消息内容
        "FromUserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 自己的username
        "ToUserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 对方的username
        "LocalID": "14900000000000000",
        "ClientMsgId": "14900000000000000"
    },
    "Scene": 0
}
```

---

## 撤回消息

| Key         | Value                                                      | Remark |
| ----------- | ---------------------------------------------------------- | ------ |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxrevokemsg> |        |
| Method      | POST                                                       |        |

Body (json):

``` json
{
    "BaseRequest": {
        "Uin": 200000000,
        "Sid": "xxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e680000000000000"
    },
    "SvrMsgId": "5910000000000000000", // 发送消息时服务器返回的服务器端消息ID
    "ToUserName": "@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 目标用户
    "ClientMsgId": "14900000000000000" // 当时的本地消息ID
}
```

## 转发被撤回的图片消息

要转发被撤回的图片消息，只需将撤回的图片消息的Content中的aeskey(cdnthumbaeskey)、cdnthumburl(cdnmidimgurl)、md5复制到发送图片的对应字段中即可
