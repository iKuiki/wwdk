# 接收

接受消息主要分为2个步骤

- 轮询微信状态同步接口，获得当前微信状态
- 如果当前微信状态为有新消息，则调用获取新消息接口取得新消息

注：**同步状态接口中有个synckey，这个synckey是在登陆时的初始化(webwxinit)操作中获取的，并且每次调用getMessage接口都会刷新这个synckey**

---

## 同步状态

| Key         | Value                                                      | Remark                                                |
| ----------- | ---------------------------------------------------------- | ----------------------------------------------------- |
| Request URL | <https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck> |                                                       |
| Method      | Get                                                        |                                                       |
| Cookie      | Need                                                       |                                                       |
| Param       | r                                                          | 13位unix时间戳                                        |
| Param       | skey                                                       | 登陆凭据skey                                          |
| Param       | sid                                                        | 登陆凭据sid                                           |
| Param       | uin                                                        | 登陆凭据uin                                           |
| Param       | deviceid                                                   |                                                       |
| Param       | synckey                                                    | 同步key，由"\|"分割为1_xxx\|2_xxx\|3_xxx\|4_xxx的格式 |
| Param       | _                                                          | 13位unix时间戳                                        |

**Response:**

| Key      | Value               | Remark                                                |
| -------- | ------------------- | ----------------------------------------------------- |
| retcode  | 0<br/>1101          | 正常<br/>已退出登陆                                   |
| selector | 0<br/>2<br/>4<br/>7 | 正常<br/>有新消息<br/>联系人有更新<br/>手机点击联系人 |

*Example:*

``` javascript
window.synccheck={retcode:"0",selector:"2"}
```
