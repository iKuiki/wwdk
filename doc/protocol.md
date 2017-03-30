# 微信网页版通信协议-2017.3.30
此协议分析于2017年3月30日

请求中凡需要Cookie的，都需要包含以下Cookie字段
| 需要的Cookie |
|-------------|
| webwxuvid |
| webwx_auth_ticket |
| wxuin |
| wxsid |
| webwx_data_ticket |

---

#### 发送消息
| Key | Value |
|-----|-------|
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg |
| Method | POST |
| Cookie | Need |
| Param | pass_ticket |
###### Body (Json):
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


#### 撤回消息
| Key | Value |
|-----|-------|
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxrevokemsg |
| Method | POST |
| Cookie | Need |
###### Body (Json):
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
