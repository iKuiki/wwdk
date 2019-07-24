# 群管理

- [群管理](#%e7%be%a4%e7%ae%a1%e7%90%86)
  - [修改群名](#%e4%bf%ae%e6%94%b9%e7%be%a4%e5%90%8d)

## 修改群名

| Key         | Value                                                                     |
| ----------- | ------------------------------------------------------------------------- |
| Request URL | <https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxupdatechatroom?fun=modtopic> |
| Method      | POST                                                                      |
| Cookie      | Need                                                                      |

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
