
# 发送

- [发送](#%e5%8f%91%e9%80%81)
  - [消息已读](#%e6%b6%88%e6%81%af%e5%b7%b2%e8%af%bb)
  - [发送消息](#%e5%8f%91%e9%80%81%e6%b6%88%e6%81%af)
  - [撤回消息](#%e6%92%a4%e5%9b%9e%e6%b6%88%e6%81%af)
  - [转发被撤回的图片消息](#%e8%bd%ac%e5%8f%91%e8%a2%ab%e6%92%a4%e5%9b%9e%e7%9a%84%e5%9b%be%e7%89%87%e6%b6%88%e6%81%af)
  - [上传文件](#%e4%b8%8a%e4%bc%a0%e6%96%87%e4%bb%b6)
  - [发送图片消息](#%e5%8f%91%e9%80%81%e5%9b%be%e7%89%87%e6%b6%88%e6%81%af)
  - [发送视频消息](#%e5%8f%91%e9%80%81%e8%a7%86%e9%a2%91%e6%b6%88%e6%81%af)
  - [发送动图消息](#%e5%8f%91%e9%80%81%e5%8a%a8%e5%9b%be%e6%b6%88%e6%81%af)
  - [发送文件](#%e5%8f%91%e9%80%81%e6%96%87%e4%bb%b6)

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
    "Code":1, // 此处的Code含义还需要判断，见过的Code包括1、3，目前发现的规律是，如果Code为3，则FromUserName=ToUserName，即自己发给自己
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
        "LocalID": "14900000000000000", // 直接传13位Unix时间戳即可
        "ClientMsgId": "14900000000000000" // 直接传13位Unix时间戳即可
    },
    "Scene": 0
}
```

**Response:**

返回是一个json对象，里面有MsgID，撤回消息时需要用到MsgID

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    },
    "MsgID": "2900000000000000000",
    "LocalID": "15600000000000000"
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

**Response:**

没啥意义的Json

``` json
{
    "BaseResponse": {
        "Ret": 0, // 已知如果在已经不存在的群发送消息，会返回1201
        "ErrMsg": ""
    },
    "Introduction": "你可以撤回2分钟内发送的消息（部分旧版本微信不支持这个功能）。",
    "SysWording": ""
}
```

---

## 转发被撤回的图片消息

要转发被撤回的图片消息，只需将撤回的图片消息的Content中的aeskey(cdnthumbaeskey)、cdnthumburl(cdnmidimgurl)、md5复制到发送图片的对应字段中即可

## 上传文件

不管是发图片还是视频还是传文件，都要先将文件上传
对于上传文件，Web微信使用统一的上传接口。上传接口因为要上传文件，所以使用了multipart表单的请求方式，并且大部分请求参数都在multipart表单中。请求如下

| Key         | Value                                                             | Remark                                                                                     |
| ----------- | ----------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| Request URL | <https://file.{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxuploadmedia> |                                                                                            |
| Method      | POST                                                              |                                                                                            |
| Header      | Content-Type                                                      | 形如这样的，由multipart组件生成multipart/form-data; boundary=----WebKitFormBoundarxxxxxxxx |
| Param       | f                                                                 | 填json                                                                                     |
| multipart   | id                                                                | WU_FILE_?，其中?为自增数字，每上传一个文件自增1                                            |
| multipart   | name                                                              | 文件名                                                                                     |
| multipart   | type                                                              | 文件的mine类型，详情见下表                                                                 |
| multipart   | lastModifiedDate                                                  | 最后编辑时间,格式为Mon Jan 02 2006 15:04:05 GMT+0700 (MST)                                 |
| multipart   | size                                                              | 文件大小                                                                                   |
| multipart   | mediatype                                                         | 文件类型，详见下表                                                                         |
| multipart   | uploadmediarequest                                                | 文件上传请求，json封装的BaseRequest等信息,下详                                             |
| multipart   | webwx_data_ticket                                                 | cookie中的数据ticket                                                                       |
| multipart   | pass_ticket                                                       | 部分Domain需要传，保险起见可以都传                                                         |
| multipart   | filename                                                          | 文件名以及要上传的文件本体                                                                 |

上面的type参数还有mediatype参数需要根据上传文件类型的不同而有不同，尝试后如下表

| 文件类型\对应字段 | type                | mediatype |
| ----------------- | ------------------- | --------- |
| gif               | image/gif           | doc       |
| png               | image/png           | pic       |
| jpg               | image/jpeg          | pic       |
| mp3               | audio/mp3           | doc       |
| aac               | audio/vnd.dlna.adts | doc       |
| mp4               | video/mp4           | video     |
| zip               | application/zip     | doc       |

uploadmediarequest字段需要将一下json对象作为字符串传入

``` json
{
    "UploadType": 2, // 固定填2
    "BaseRequest": {
        "Uin": 200000000,
        "Sid": "xxxxxxxxxxxxxxxx",
        "Skey": "@crypt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID": "e680000000000000"
    },
    "ClientMediaId": 1560000000000, // 直接传入13位Unix时间戳
    "TotalLen": 1024, // 文件大小
    "StartPos": 0,
    "DataLen": 1024, // 文件大小
    "MediaType": 4, // 固定填4
    "FromUserName": "@xxxx", // 发送人的UserName（当然填自己拉
    "ToUserName": "@@xxxxxxxxxx", // 接收人的UserName
    "FileMd5": "xxxxxxxxx" // 文件的md5值
}
```

**Response:**

返回为一个json对象，其中有用消息为MediaId

``` json
{
    "BaseResponse": {
        "Ret": 0,
        "ErrMsg": ""
    },
    "MediaId": "@crypt_xxxxxxxxxxxxxxxxxxxxx", // 最有用的就是这个MediaId了
    "StartPos": 1024, // 与大小有关，可能是压缩后的大小
    "CDNThumbImgHeight": 128, // 缩略图分辨率
    "CDNThumbImgWidth": 128, // 缩略图分辨率
    "EncryFileName": "xxxx%2Emp4" // urlEncode后的文件名
}
```

---

## 发送图片消息

发送图片消息需要先调用上传文件的接口将文件上传后获取MediaID，获取到MediaID后调用如下接口

| Key         | Value                                                       | Remark                             |
| ----------- | ----------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxsendmsgimg> |                                    |
| Method      | POST                                                        |                                    |
| Param       | fun                                                         | 填async                            |
| Param       | f                                                           | 填json                             |
| Param       | pass_ticket                                                 | 部分Domain需要传，保险起见可以都传 |

**Body (Json):**

``` json
{
    "BaseRequest": {
        "Uin":"210000000",
        "Sid":"QQxxxxxxxxxxxxxx",
        "Skey":"@crypt_a6xxxxxx_6xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID":"e980000000000000"
    },
    "Msg": {
        "Type": 3, // 图片消息的MsgType为3
        "MediaId": "@crypt_xxxxxxxxxxx", // 填通过上传文件接口获取到的MediaId
        "Content": "", // 图片消息没有内容
        "FromUserName": "@xxxx", // 发件人UserName，填自己的
        "ToUserName": "@@xxxxxxxx", // 收件人的UserName
        "LocalID": "15600000000000000", // 官方用的17位Unix时间戳，不过用14位应该也可以
        "ClientMsgId": "15600000000000000" // 官方用的17位Unix时间戳，不过用14位应该也可以
    },
    "Scene": 0 // 填0
}
```

发送图片消息后的返回与发送文字消息后的返回一毛一样，我就不重复了

---

## 发送视频消息

发送视频消息需要先调用上传文件的接口将文件上传后获取MediaID，获取到MediaID后调用如下接口

| Key         | Value                                                         | Remark                             |
| ----------- | ------------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxsendvideomsg> |                                    |
| Method      | POST                                                          |                                    |
| Param       | fun                                                           | 填async                            |
| Param       | f                                                             | 填json                             |
| Param       | pass_ticket                                                   | 部分Domain需要传，保险起见可以都传 |

**Body (Json):**

``` json
{
    "BaseRequest": {
        "Uin":"210000000",
        "Sid":"QQxxxxxxxxxxxxxx",
        "Skey":"@crypt_a6xxxxxx_6xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID":"e980000000000000"
    },
    "Msg": {
        "Type": 43, // 视频消息的MsgType为43
        "MediaId": "@crypt_xxxxxxxxxxx", // 填通过上传文件接口获取到的MediaId
        "Content": "", // 视频消息没有内容
        "FromUserName": "@xxxx", // 发件人UserName，填自己的
        "ToUserName": "@@xxxxxxxx", // 收件人的UserName
        "LocalID": "15600000000000000", // 官方用的17位Unix时间戳，不过用14位应该也可以
        "ClientMsgId": "15600000000000000" // 官方用的17位Unix时间戳，不过用14位应该也可以
    },
    "Scene": 0 // 填0
}
```

发送视频消息后的返回与发送文字消息后的返回也一毛一样，我就也不重复了

---

## 发送动图消息

发送动图消息需要先调用上传文件的接口将文件上传后获取MediaID，获取到MediaID后调用如下接口

| Key         | Value                                                         | Remark                             |
| ----------- | ------------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxsendemoticon> |                                    |
| Method      | POST                                                          |                                    |
| Param       | fun                                                           | 填sys                              |
| Param       | pass_ticket                                                   | 部分Domain需要传，保险起见可以都传 |

**Body (Json):**

``` json
{
    "BaseRequest": {
        "Uin":"210000000",
        "Sid":"QQxxxxxxxxxxxxxx",
        "Skey":"@crypt_a6xxxxxx_6xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID":"e980000000000000"
    },
    "Msg": {
        "Type": 47, // 动图消息的MsgType为47
        "EmojiFlag": 2, // 固定填2
        "MediaId": "@crypt_xxxxxxxxxxx", // 填通过上传文件接口获取到的MediaId
        "FromUserName": "@xxxx",
        "ToUserName": "@@xxxxxxxx",
        "LocalID": "15600000000000000", // 官方用的17位Unix时间戳，不过用14位应该也可以
        "ClientMsgId": "15600000000000000" // 官方用的17位Unix时间戳，不过用14位应该也可以
    },
    "Scene": 0 // 填0
}
```

发送动图消息后的返回与发送文字消息后的返回也一毛一样，我就也不重复了

---

## 发送文件

发送文件需要先调用上传文件的接口将文件上传后获取MediaID，获取到MediaID后调用如下接口

| Key         | Value                                                       | Remark                             |
| ----------- | ----------------------------------------------------------- | ---------------------------------- |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxsendappmsg> |                                    |
| Method      | POST                                                        |                                    |
| Param       | fun                                                         | 填async                            |
| Param       | f                                                           | 填json                             |
| Param       | pass_ticket                                                 | 部分Domain需要传，保险起见可以都传 |

**Body (Json):**

``` json
{
    "BaseRequest": {
        "Uin":"210000000",
        "Sid":"QQxxxxxxxxxxxxxx",
        "Skey":"@crypt_a6xxxxxx_6xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "DeviceID":"e980000000000000"
    },
    "Msg": {
        "Type": 6, // 文件的MsgType为6
        // Content内为一个xml对象，包含了appid(固定)、上传的文件名、文件大小、上传后的MediaId、文件后缀这几个信息
        "Content": "<appmsg appid='wxeb7ec651dd0aefa9' sdkver=''><title>fileName.txt</title><des></des><action></action><type>6</type><content></content><url></url><lowurl></lowurl><appattach><totallen>1024</totallen><attachid>@crypt_xxxxxxxxxxxx</attachid><fileext>txt</fileext></appattach><extinfo></extinfo></appmsg>",
        "FromUserName": "@xxxxxx", // 发件人UserName，填自己的
        "ToUserName": "@@xxxxxxxxx", // 收件人的UserName
        "LocalID": "15600000000000000", // 官方用的17位Unix时间戳，不过用14位应该也可以
        "ClientMsgId": "15600000000000000" // 官方用的17位Unix时间戳，不过用14位应该也可以
    },
    "Scene": 0 // 填0
}
```

上面Msg中的Content是一个这样的xml对象

``` xml
<appmsg appid='wxeb7ec651dd0aefa9' sdkver=''> <!--appid为固定的，就填wxeb7ec651dd0aefa9即可-->
    <title>fileName.txt</title> <!--文件名-->
    <des></des>
    <action></action>
    <type>6</type> <!--固定填写6-->
    <content></content>
    <url></url>
    <lowurl></lowurl>
    <appattach>
        <totallen>1024</totallen> <!--文件大小-->
        <attachid>@crypt_xxxxxxxxxxxx</attachid> <!--调用上传文件接口获得的MediaId-->
        <fileext>txt</fileext> <!--文件后缀-->
    </appattach>
    <extinfo></extinfo>
</appmsg>
```

发送文件后的返回与发送文字消息后的返回也一毛一样，我就也不重复了
