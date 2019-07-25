# 群管理

- [群管理](#%e7%be%a4%e7%ae%a1%e7%90%86)
  - [修改群名](#%e4%bf%ae%e6%94%b9%e7%be%a4%e5%90%8d)

## 修改群名

| Key         | Value                                                                        |
| ----------- | ---------------------------------------------------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxupdatechatroom?fun=modtopic> |
| Method      | POST                                                                         |
| Cookie      | Need                                                                         |

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
