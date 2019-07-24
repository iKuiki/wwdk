# 联系人管理

- [联系人管理](#%e8%81%94%e7%b3%bb%e4%ba%ba%e7%ae%a1%e7%90%86)
  - [修改用户备注](#%e4%bf%ae%e6%94%b9%e7%94%a8%e6%88%b7%e5%a4%87%e6%b3%a8)

## 修改用户备注

| Key         | Value                                               |
| ----------- | --------------------------------------------------- |
| Request URL | <https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxoplog> |
| Method      | POST                                                |
| Cookie      | Need                                                |

**Body (json):**

``` json
{
    "UserName": "@f5fc0106e419ed58baafd50a9d4b4f4869d417411cd834ffe4f43ac62bbc38a6",
    "CmdId": 2,
    "RemarkName": "3123",
    "BaseRequest": {
        "Uin": 216547950,
        "Sid": "2ZgAp8arXklWje6v",
        "Skey": "@crypt_a6a25b34_bdf4052fc37832fabf2c2dc5421f8fca",
        "DeviceID": "e968379522819118"
    }
}
```
