# 接收信息

当webwxsync接口接收到新信息时，媒体类的信息需要二次处理

- [接收信息](#%e6%8e%a5%e6%94%b6%e4%bf%a1%e6%81%af)
  - [接收图片](#%e6%8e%a5%e6%94%b6%e5%9b%be%e7%89%87)
  - [接收语音](#%e6%8e%a5%e6%94%b6%e8%af%ad%e9%9f%b3)
  - [接收视频](#%e6%8e%a5%e6%94%b6%e8%a7%86%e9%a2%91)
  - [接收动图](#%e6%8e%a5%e6%94%b6%e5%8a%a8%e5%9b%be)
  - [接收文件](#%e6%8e%a5%e6%94%b6%e6%96%87%e4%bb%b6)
  - [接收名片](#%e6%8e%a5%e6%94%b6%e5%90%8d%e7%89%87)
  - [接收定位](#%e6%8e%a5%e6%94%b6%e5%ae%9a%e4%bd%8d)
  - [接收撤回消息](#%e6%8e%a5%e6%94%b6%e6%92%a4%e5%9b%9e%e6%b6%88%e6%81%af)

---

## 接收图片

当接收到图片消息时，Msg会是如下模样(MsgType=3)
其中需要用到的消息只有MsgID，其余内容并未发现有用之处

``` json
{
    "MsgId": "4600000000000000000",
    "FromUserName": "@@xxxxxxxxxxxxxxxxx",
    "ToUserName": "@xxxxxxxxxxx",
    "MsgType": 3, // type = 3代表是图片消息
    // Content中如果是群聊,则有发言人的UserName
    // 需要小心一个陷阱,如果是群聊中自己发的,Content里是一个xml数据,前面是没用UserName的,所以不要尝试获取UserName,会异常
    "Content": "@xxxxxxxxxxx:<br/>@xxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "Status": 3,
    "ImgStatus": 2,
    "CreateTime": 1560000000,
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
        "UserName": "",
        "NickName": "",
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "",
        "Signature": "",
        "Alias": "",
        "Scene": 0,
        "VerifyFlag": 0,
        "AttrStatus": 0,
        "Sex": 0,
        "Ticket": "",
        "OpCode": 0
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
    "NewMsgId": 4600000000000000000,
    "OriContent": "",
    "EncryFileName": ""
}
```

要下载图片内容，则调用以下api
注：**type参数用于请求缩略图，请求原图时不要添加该参数**

| Key         | Value                                                      | Remark                                     |
| ----------- | ---------------------------------------------------------- | ------------------------------------------ |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxgetmsgimg> |                                            |
| Method      | Get                                                        |                                            |
| Param       | MsgID                                                      | 填消息的MsgID(注意参数名大小写)            |
| Param       | skey                                                       | 填登陆信息中的skey                         |
| Param       | type                                                       | 填slave为缩略图,不填为原图,如果是动图用big |

**Response:**

返回值的body可以直接保存为图片，文件类型可以参考Header中的Content-Type

---

## 接收语音

语音消息的Msg如下(MsgType=34)

``` json
{
    "MsgId": "6100000000000000000",
    "FromUserName": "@@xxxxxxxxxxxxxxxx",
    "ToUserName": "@xxxxxxxxxxx",
    "MsgType": 34, // 语音消息MsgType为34
    // Content中如果是群聊,则有发言人的UserName
    // 需要小心一个陷阱,如果是群聊中自己发的,Content里是一个xml数据,前面是没用UserName的,所以不要尝试获取UserName,会异常
    "Content": "@xxxxxxxxxxx:<br/>@xxxxxxxxxxxxxxxxxxxxxxx",
    "Status": 3,
    "ImgStatus": 1,
    "CreateTime": 1560000000,
    "VoiceLength": 3215, // 声音长度(毫秒数)
    "PlayLength": 0,
    "FileName": "",
    "FileSize": "",
    "MediaId": "",
    "Url": "",
    "AppMsgType": 0,
    "StatusNotifyCode": 0,
    "StatusNotifyUserName": "",
    "RecommendInfo": {
        "UserName": "",
        "NickName": "",
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "",
        "Signature": "",
        "Alias": "",
        "Scene": 0,
        "VerifyFlag": 0,
        "AttrStatus": 0,
        "Sex": 0,
        "Ticket": "",
        "OpCode": 0
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
    "NewMsgId": 6100000000000000000,
    "OriContent": "",
    "EncryFileName": ""
}
```

要下载声音可以通过如下接口下载

注：**type参数用于请求缩略图，请求原图时不要添加该参数**

| Key         | Value                                                     | Remark                   |
| ----------- | --------------------------------------------------------- | ------------------------ |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxgetvoice> |                          |
| Method      | Get                                                       |                          |
| Header      | Range: bytes=0-                                           | 很隐蔽，一不小心要漏掉了 |
| Param       | msgid                                                     | 填消息的MsgID            |
| Param       | skey                                                      | 填登陆信息中的skey       |

**Response:**

返回值的body可以直接保存为音频，文件类型可以参考Header中的Content-Type，目前所见都是mp3

---

## 接收视频

如果接到视频消息，Msg会长的像下面这样(MsgType=43)

``` json
{
    "MsgId": "6000000000000000000",
    "FromUserName": "@xxxxxxxxxxx",
    "ToUserName": "@@xxxxxxxxxxxxxxxxxxxxxxxxx",
    "MsgType": 43, // 43则为视频消息
    // Content中如果是群聊,则有发言人的UserName
    // 需要小心一个陷阱,如果是群聊中自己发的,Content里是一个xml数据,前面是没用UserName的,所以不要尝试获取UserName,会异常
    "Content": "@xxxx:<br/>@xxxxxxx",
    "Status": 3,
    "ImgStatus": 1,
    "CreateTime": 1560000000,
    "VoiceLength": 0,
    "PlayLength": 1,
    "FileName": "",
    "FileSize": "",
    "MediaId": "",
    "Url": "",
    "AppMsgType": 0,
    "StatusNotifyCode": 0,
    "StatusNotifyUserName": "",
    "RecommendInfo": {
        "UserName": "",
        "NickName": "",
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "",
        "Signature": "",
        "Alias": "",
        "Scene": 0,
        "VerifyFlag": 0,
        "AttrStatus": 0,
        "Sex": 0,
        "Ticket": "",
        "OpCode": 0
    },
    "ForwardFlag": 0,
    "AppInfo": {
        "AppID": "",
        "Type": 0
    },
    "HasProductId": 0,
    "Ticket": "",
    "ImgHeight": 400, // 宽
    "ImgWidth": 200, // 高
    "SubMsgType": 0,
    "NewMsgId": 6000000000000000000,
    "OriContent": "",
    "EncryFileName": ""
}
```

视频的缩略图可以调用查看图片消息的缩略图的接口获取，要下载视频可以通过如下接口下载

注：**type参数用于请求缩略图，请求原图时不要添加该参数**

| Key         | Value                                                     | Remark                   |
| ----------- | --------------------------------------------------------- | ------------------------ |
| Request URL | <https://{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxgetvideo> |                          |
| Method      | Get                                                       |                          |
| Header      | Range: bytes=0-                                           | 很隐蔽，一不小心要漏掉了 |
| Param       | msgid                                                     | 填消息的MsgID            |
| Param       | skey                                                      | 填登陆信息中的skey       |

**Response:**

返回值的body可以直接保存为视频，文件类型可以参考Header中的Content-Type，一般而言基本是mp4

---

## 接收动图

如果MsgType为47则为动态表情消息,消息本地如下

``` json
{
    "MsgId": "67000000000000000",
    "FromUserName": "@@xxxxxxxxxxxxxxxxxxxxxxxxx",
    "ToUserName": "@xxxxxxxxxxxxxxxxxxxxxxxxx",
    "MsgType": 47, // MsgType 47为动态表情消息
    // Content中无实质性消息
    // Content中如果是群聊,则有发言人的UserName
    // 需要小心一个陷阱,如果是群聊中自己发的,前面是没用UserName的,所以不要尝试获取UserName,会异常
    "Content": "@xxxxxxxxxxxxxxxxxxxxxxxxx:<br/>&lt;msg&gt;&lt;emoji fromusername=\"xxx\" tousername=\"xxx\" type=\"2\" idbuffer=\"media:0_0\" md5=\"xxxx\" len=\"1024\" productid=\"\" androidmd5=\"xxxx\" androidlen=\"1024\" s60v3md5=\"xxxx\" s60v3len=\"1024\" s60v5md5=\"xxxx\" s60v5len=\"1024\" cdnurl=\"http://emoji.qpic.cn/wx_emoji/xxxxxxxx/\" designerid=\"\" thumburl=\"\" encrypturl=\"http://emoji.qpic.cn/wx_emoji/xxxx/\" aeskey=\"xxxx\" externurl=\"http://emoji.qpic.cn/wx_emoji/xxxx/\" externmd5=\"xxxx\" width=\"80\" height=\"80\" tpurl=\"\" tpauthkey=\"\" attachedtext=\"\" attachedtextcolor=\"\" lensid=\"\"&gt;&lt;/emoji&gt;&lt;/msg&gt;",
    "Status": 3,
    "ImgStatus": 2,
    "CreateTime": 1560000000,
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
        "UserName": "",
        "NickName": "",
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "",
        "Signature": "",
        "Alias": "",
        "Scene": 0,
        "VerifyFlag": 0,
        "AttrStatus": 0,
        "Sex": 0,
        "Ticket": "",
        "OpCode": 0
    },
    "ForwardFlag": 0,
    "AppInfo": {
        "AppID": "",
        "Type": 0
    },
    "HasProductId": 0,
    "Ticket": "",
    "ImgHeight": 80,
    "ImgWidth": 80,
    "SubMsgType": 0,
    "NewMsgId": 67000000000000000,
    "OriContent": "",
    "EncryFileName": ""
}
```

获取动图本地可以调用保存图片的接口webwxgetmsgimg来保存即可,不过type参数在官方案例中是big

---

## 接收文件

当收到文件消息，Msg会长下面这样

``` json
{
    "MsgId": "8800000000000000000",
    "FromUserName": "@@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "ToUserName": "@xxxxxxxxxxxxxxxxx",
    "MsgType": 49,
    "Content": "@xxxxxxxxxxxxxxxxx:<br/>&lt;msg&gt;&lt;appmsg appid=\"wx6618f1cfc6c132f8\" sdkver=\"0\"&gt;&lt;title&gt;fileName.zip&lt;/title&gt;&lt;des&gt;&lt;/des&gt;&lt;action&gt;view&lt;/action&gt;&lt;type&gt;6&lt;/type&gt;&lt;showtype&gt;0&lt;/showtype&gt;&lt;content&gt;&lt;/content&gt;&lt;url&gt;&lt;/url&gt;&lt;dataurl&gt;&lt;/dataurl&gt;&lt;lowurl&gt;&lt;/lowurl&gt;&lt;lowdataurl&gt;&lt;/lowdataurl&gt;&lt;recorditem&gt;&lt;![CDATA[]]&gt;&lt;/recorditem&gt;&lt;thumburl&gt;&lt;/thumburl&gt;&lt;messageaction&gt;&lt;/messageaction&gt;&lt;extinfo&gt;&lt;/extinfo&gt;&lt;sourceusername&gt;&lt;/sourceusername&gt;&lt;sourcedisplayname&gt;&lt;/sourcedisplayname&gt;&lt;commenturl&gt;&lt;/commenturl&gt;&lt;appattach&gt;&lt;totallen&gt;1024&lt;/totallen&gt;&lt;attachid&gt;@cdn_qwerty_abcd_1&lt;/attachid&gt;&lt;emoticonmd5&gt;&lt;/emoticonmd5&gt;&lt;fileext&gt;zip&lt;/fileext&gt;&lt;cdnattachurl&gt;qwerty&lt;/cdnattachurl&gt;&lt;aeskey&gt;abcd&lt;/aeskey&gt;&lt;encryver&gt;1&lt;/encryver&gt;&lt;/appattach&gt;&lt;md5&gt;xxxx&lt;/md5&gt;&lt;/appmsg&gt;&lt;fromusername&gt;&lt;/fromusername&gt;&lt;scene&gt;0&lt;/scene&gt;&lt;appinfo&gt;&lt;version&gt;7&lt;/version&gt;&lt;appname&gt;微信电脑版&lt;/appname&gt;&lt;/appinfo&gt;&lt;commenturl&gt;&lt;/commenturl&gt;&lt;/msg&gt;&lt;br/&gt;<br/>",
    "Status": 3,
    "ImgStatus": 1,
    "CreateTime": 1560000000,
    "VoiceLength": 0,
    "PlayLength": 0,
    "FileName": "fileName.zip", // 文件名，重要，获取文件时需要用到
    "FileSize": "1024", // 文件大小
    "MediaId": "@crypt_xxxxxxxxxxxxxxxxxxxxxxxxxxx", // 媒体ID，重要，获取文件时需要用到
    "Url": "",
    "AppMsgType": 6,
    "StatusNotifyCode": 0,
    "StatusNotifyUserName": "",
    "RecommendInfo": {
        "UserName": "",
        "NickName": "",
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "",
        "Signature": "",
        "Alias": "",
        "Scene": 0,
        "VerifyFlag": 0,
        "AttrStatus": 0,
        "Sex": 0,
        "Ticket": "",
        "OpCode": 0
    },
    "ForwardFlag": 0,
    "AppInfo": {
        "AppID": "wx6618f1cfc6c132f8", // 这个应该是固定固定的，用于区分什么设备发送的
        "Type": 0
    },
    "HasProductId": 0,
    "Ticket": "",
    "ImgHeight": 0,
    "ImgWidth": 0,
    "SubMsgType": 0,
    "NewMsgId": 8800000000000000000,
    "OriContent": "",
    "EncryFileName": "fileName%2Ezip" // 文件名进行urlEncode后的产物
}
```

其中content中的内容，进行url Unescape后可以得到如下xml
目前来看，这个xml意义不大，下载文件的时候不需要用到其中的信息，里面只有发件客户端名可以用以显示用途

``` xml
<msg>
    <appmsg appid=\"wx6618f1cfc6c132f8\" sdkver=\"0\">
        <title>fileName.zip</title>
        <des></des>
        <action>view</action>
        <type>6</type>
        <showtype>0</showtype>
        <content></content>
        <url></url>
        <dataurl></dataurl>
        <lowurl></lowurl>
        <lowdataurl></lowdataurl>
        <recorditem>
            <![CDATA[]]>
        </recorditem>
        <thumburl></thumburl>
        <messageaction></messageaction>
        <extinfo></extinfo>
        <sourceusername></sourceusername>
        <sourcedisplayname></sourcedisplayname>
        <commenturl></commenturl>
        <appattach>
            <totallen>1024</totallen> <!--文件大小-->
            <attachid>@cdn_qwerty_abcd_1</attachid> <!--这里这个attachid是下面的cdnattachurl与aeskey拼接后的参数-->
            <emoticonmd5></emoticonmd5>
            <fileext>zip</fileext>
            <cdnattachurl>qwerty</cdnattachurl>
            <aeskey>abcd</aeskey>
            <encryver>1</encryver>
        </appattach>
        <md5>xxxx</md5>
    </appmsg>
    <fromusername></fromusername>
    <scene>0</scene>
    <appinfo>
        <version>7</version>
        <appname>微信电脑版</appname> <!--发送这个文件的客户端-->
    </appinfo>
    <commenturl></commenturl>
</msg>
<br/>
```

要将这个附件下载回来，需要调用以下地址，会直接返回下载文件内容

| Key         | Value                                                          | Remark                                 |
| ----------- | -------------------------------------------------------------- | -------------------------------------- |
| Request URL | <https://file.{{apiDomain}}/cgi-bin/mmwebwx-bin/webwxgetmedia> |                                        |
| Method      | Get                                                            |                                        |
| Param       | sender                                                         | 发送这条消息的联系人（或群）的UserName |
| Param       | mediaid                                                        | 上面消息json中的mediaID                |
| Param       | encryfilename                                                  | 上面消息json中的文件名                 |
| Param       | fromuser                                                       | 填自己登陆信息中的wxuin                |
| Param       | pass_ticket                                                    | 部分Domain需要传，保险起见可以都传     |
| Param       | webwx_data_ticket                                              | 填自己登陆信息中的DataTicket           |

**Response:**

返回值的body可以直接保存为文件，文件名在Header的Content-Disposition字段中

---

## 接收名片

接收到名片后，Msg如下（MsgType=42）

``` json
{
    "MsgId": "8500000000000000000",
    "FromUserName": "@@xxxxxxxxxxx",
    "ToUserName": "@xxxxxxxxx",
    "MsgType": 42, // 名片消息MsgType为42
    // Content中是一个xml对象，个人名片字段更多(但基本都是空值)，公众号名片少几个字段，但是都有填值
    "Content": "@xxxxxxxxxxx:<br/>&lt;?xml version=\"1.0\"?&gt;<br/>&lt;msg bigheadimgurl=\"http://wx.qlogo.cn/mmhead/ver_1/xxxx/0\" smallheadimgurl=\"http://wx.qlogo.cn/mmhead/ver_1/xxxx/132\" username=\"v1_xxxx@stranger\" nickname=\"xxxx\"  shortpy=\"\" alias=\"\" imagestatus=\"3\" scene=\"17\" province=\"\" city=\"\" sign=\"\" sex=\"0\" certflag=\"0\" certinfo=\"\" brandIconUrl=\"\" brandHomeUrl=\"\" brandSubscriptConfigUrl=\"\" brandFlags=\"0\" regionCode=\"\" antispamticket=\"v2_xxxx@stranger\" /&gt;<br/>",
    "Status": 3,
    "ImgStatus": 1,
    "CreateTime": 1560000000,
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
        "UserName": "@xxxxxxxxxxxxxxx", // 名片中对方的UserName
        "NickName": "xxx", // 对方的昵称
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "",
        "Signature": "",
        "Alias": "",
        "Scene": 17, // 固定值17？
        "VerifyFlag": 0, // 这个VerifyFlag可以用来判定公众号，公众号是24
        "AttrStatus": 32,
        "Sex": 0,
        "Ticket": "",
        "OpCode": 0
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
    "NewMsgId": 8500000000000000000,
    "OriContent": "",
    "EncryFileName": ""
}
```

现在收到名片后无法添加为联系人or关注，所以没有后续处理了

---

## 接收定位

定位的Msg本体如下，其与文字消息的MsgType相同

``` json
{
    "MsgId": "790000000000000",
    "FromUserName": "@xxxxxxxxxxxxxxxx",
    "ToUserName": "@xxxxxxxxx",
    "MsgType": 1,
    // 地名，然后后面接着获取地址地图截图的url
    "Content": "深圳市xxxx:<br/>/cgi-bin/mmwebwx-bin/webwxgetpubliclinkimg?url=xxx&msgid=790000000000000&pictype=location",
    "Status": 3,
    "ImgStatus": 1,
    "CreateTime": 1560000000,
    "VoiceLength": 0,
    "PlayLength": 0,
    "FileName": "",
    "FileSize": "",
    "MediaId": "",
    "Url": "http://apis.map.qq.com/uri/v1/geocoder?coord=22.000000,113.000000", // 跳转向定位的链接，内含腾讯系的坐标
    "AppMsgType": 0,
    "StatusNotifyCode": 0,
    "StatusNotifyUserName": "",
    "RecommendInfo": {
        "UserName": "",
        "NickName": "",
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "",
        "Signature": "",
        "Alias": "",
        "Scene": 0,
        "VerifyFlag": 0,
        "AttrStatus": 0,
        "Sex": 0,
        "Ticket": "",
        "OpCode": 0
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
    "SubMsgType": 48, // 可能这个是关键
    "NewMsgId": 790000000000000,
    // OriContent中包含一个xml，内含经纬度、地名等信息
    "OriContent": "<?xml version=\"1.0\"?>\n<msg>\n\t<location x=\"22.000000\" y=\"113.000000\" scale=\"16\" label=\"深圳市xxxx\" maptype=\"0\" poiname=\"[位置]\" poiid=\"\" />\n</msg>\n",
    "EncryFileName": ""
}
```

如果要获取地址地图截图，将Content中地名后面的uri拼接到{{apiDomain}}后面即可

---

## 接收撤回消息

当收到撤回消息是，Msg本体如下，MsgType=10002

``` json
{
    "MsgId": "55000000000000000",
    "FromUserName": "@xxxxxxxxxxxxxxxxxxx",
    "ToUserName": "@xxxxxxxxxxxxxxxxxxx",
    "MsgType": 10002,
    "Content": "&lt;sysmsg type=\"revokemsg\"&gt;<br/>\t&lt;revokemsg&gt;<br/>\t\t&lt;session&gt;wxid_xxxxxxx&lt;/session&gt;<br/>\t\t&lt;oldmsgid&gt;1660000000&lt;/oldmsgid&gt;<br/>\t\t&lt;msgid&gt;41000000000000000&lt;/msgid&gt;<br/>\t\t&lt;replacemsg&gt;&lt;![CDATA[\"xx\" 撤回了一条消息]]&gt;&lt;/replacemsg&gt;<br/>\t&lt;/revokemsg&gt;<br/>&lt;/sysmsg&gt;<br/>",
    "Status": 4,
    "ImgStatus": 1,
    "CreateTime": 1560000000,
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
        "UserName": "",
        "NickName": "",
        "QQNum": 0,
        "Province": "",
        "City": "",
        "Content": "",
        "Signature": "",
        "Alias": "",
        "Scene": 0,
        "VerifyFlag": 0,
        "AttrStatus": 0,
        "Sex": 0,
        "Ticket": "",
        "OpCode": 0
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
    "NewMsgId": 55000000000000000,
    "OriContent": "",
    "EncryFileName": ""
}
```

撤回消息中很重要的被撤回的MsgID在消息content中的xml里，将content进行Unescape后得进行html渲染（主要是\t和br换行）到如下xml对象

``` xml
<sysmsg type=\"revokemsg\">
    <revokemsg>
        <session>wxid_xxxxxxx</session>
        <oldmsgid>1660000000</oldmsgid>
        <msgid>41000000000000000</msgid><!--就这条最重要，这条标记着被撤回的消息的ID-->
        <replacemsg><![CDATA[\"xx\" 撤回了一条消息]]></replacemsg>
    </revokemsg>
</sysmsg>
```
