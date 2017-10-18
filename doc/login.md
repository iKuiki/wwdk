# 登陆流程

---

登陆流程大致为以下步骤：
- 获取uuid
- 根据uuid生成二维码
- 轮询微信服务器检测用户是否扫码
- 手机扫码
- 获取到用户已扫码，得到用户头像（base64编码的字符串）
- 用户同意登陆
- 登陆成功

---
*注：微信网页版API的返回包括一种特别的格式：看起来像js代码，每个字段作为一行js代码，以分号结尾，每句以等号分割左边为key右边为code*
例：
```
window.QRLogin.code = 200; window.QRLogin.uuid = "gfNHoe0rgA==";
```
其中包含两个值：

| Key                 | Value        |
| ------------------- | ------------ |
| window.QRLogin.code | 200          |
| window.QRLogin.uuid | gfNHoe0rgA== |

解析方案：
1. 使用[github.com/robertkrimen/otto](https://github.com/robertkrimen/otto)解释后获取值（优点：可靠性高；缺点：因为要运行js，解析速度相对慢）
2. 自己通过匹配格式来解析值（优点：简单，速度快；缺点：若返回值复杂则可能解析错误）

---
### API流程

| api path          | response                             | remark  |
| ----------------- | ------------------------------------ | ------- |
| jslogin           | code<br/>uuid                        |         |
| login             | code                                 |         |
| webwxnewloginpage | skey<br/>sid<br/>uin<br/>pass_ticket | xml格式 |

---

### API详情

#### 获取uuid

| Key         | Value                           | Remark               |
| ----------- | ------------------------------- | -------------------- |
| Request URL | https://login.wx.qq.com/jslogin |                      |
| Method      | Get                             |                      |
| Cookie      | No                              |                      |
| Param       | appid                           | 填wx782c26e4c19acffb |
| Param       | fun                             | 填new                |
| Param       | lang                            | zh_CN或en_US         |
| Param       | _                               | 13位unix时间戳       |

Response:

| Key                 | Value | Remark         |
| ------------------- | ----- | -------------- |
| window.QRLogin.code | 200   |                |
| window.QRLogin.uuid | xxx   | 当前会话的uuid |

Example:
```
window.QRLogin.code = 200; window.QRLogin.uuid = "gfNHoe0rgA==";
```

---

#### 轮询用户扫码

| Key         | Value                                              | Remark                     |
| ----------- | -------------------------------------------------- | -------------------------- |
| Request URL | https://login.wx2.qq.com/cgi-bin/mmwebwx-bin/login |                            |
| Method      | Get                                                |                            |
| Param       | loginicon                                          | 填true                     |
| Param       | uuid                                               | 之前获取的uuid             |
| Param       | tip                                                | 1-未扫描 0-已扫描          |
| Param       | r                                                  | 13位时间戳取反(貌似可省略) |
| Param       | _                                                  | 13位unix时间戳             |

Response:

| Key               | Value                       | Remark                                                    |
| ----------------- | --------------------------- | --------------------------------------------------------- |
| window.code       | 200<br/>201<br/>400<br/>408 | 确认登陆<br/>已扫码<br/>登陆超时(二维码失效)<br/>等待登陆 |
| window.userAvatar | data:img/jpeg;base64        | base64编码的用户头像，仅当code=200时才有                  |

Example:
```javascript
window.code=408;window.userAvatar='data:img/jpeg;base64,iVBORw...'
```

**若用户取消登陆，返回仍为408，旧的二维码仍可重复使用，用户重新扫旧的二维码后会再次返回201**

*注：当登陆超时（二维码失效）后，重新调用获取uuid的方法即可重新拿到二维码*

---

#### 获取登陆参数

| Key         | Value                                                    | Remark                 |
| ----------- | -------------------------------------------------------- | ---------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage |                        |
| Method      | Get                                                      |                        |
| Param       | ticket                                                   |                        |
| Param       | uuid                                                     |                        |
| Param       | scan                                                     | 扫描成功后返回的时间戳 |

Response:
| Key               | Type    | Remark         |
| ----------------- | ------- | -------------- |
| wxsid             | Cookie  |                |
| wxuin             | Cookie  |                |
| webwxuvid         | Cookie  |                |
| webwx_auth_ticket | Cookie  |                |
| webwx_data_ticket | Cookie  |                |
| skey              | BodyXml |                |
| wxsid             | BodyXml | same as cookie |
| wxuin             | BodyXml | same as cookie |
| pass_ticket       | BodyXml |                |

---

#### 初始化

| Key         | Value                                            | Remark                              |
| ----------- | ------------------------------------------------ | ----------------------------------- |
| Request URL | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit |                                     |
| Method      | Post                                             |                                     |
| Param       | r                                                | 13位时间戳取反                      |
| Param       | pass_ticket                                      | 获取登陆参数时获取到的              |
| Param       | BaseRequest JsonObject                           | json对象，其中唯一元素为BaseRequest |

Json Example:
```json
{"BaseRequest":
    {
        "Uin":"216547950",
        "Sid":"QQ9iwKokvmPs7c/7",
        "Skey":"@crypt_a6a25b34_68efb91dcbec1fe990bf33d8fe770034",
        "DeviceID":"e987736822175688"
    }
}
```

Response:
返回为一个json对象，内包括用户信息、联系人(此列表不全，之后用获取联系人的接口获取完整联系人列表)等
