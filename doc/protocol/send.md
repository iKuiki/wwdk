
# 发送

- [发送](#%e5%8f%91%e9%80%81)
  - [消息已读](#%e6%b6%88%e6%81%af%e5%b7%b2%e8%af%bb)
  - [发送消息](#%e5%8f%91%e9%80%81%e6%b6%88%e6%81%af)
  - [撤回消息](#%e6%92%a4%e5%9b%9e%e6%b6%88%e6%81%af)
  - [转发被撤回的图片消息](#%e8%bd%ac%e5%8f%91%e8%a2%ab%e6%92%a4%e5%9b%9e%e7%9a%84%e5%9b%be%e7%89%87%e6%b6%88%e6%81%af)

## 消息已读

将指定联系人的消息设为已读

| Key         | Value                                                      |
| ----------- | ---------------------------------------------------------- |
| Request URL | <https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify> |
| Method      | POST                                                       |
| Cookie      | Need                                                       |
| PARAM       | pass_ticket                                                |

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

## 发送消息

| Key         | Value                                                 |
| ----------- | ----------------------------------------------------- |
| Request URL | <https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg> |
| Method      | POST                                                  |
| Cookie      | Need                                                  |
| Param       | pass_ticket                                           |

**Body (Json):**

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

## 撤回消息

| Key         | Value                                                   |
| ----------- | ------------------------------------------------------- |
| Request URL | <https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxrevokemsg> |
| Method      | POST                                                    |
| Cookie      | Need                                                    |

Body (json):

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

## 转发被撤回的图片消息

要转发被撤回的图片消息，只需将撤回的图片消息的Content中的aeskey(cdnthumbaeskey)、cdnthumburl(cdnmidimgurl)、md5复制到发送图片的对应字段中即可
