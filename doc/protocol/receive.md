# 接收信息

当webwxsync接口接收到新信息时，媒体类的信息需要二次处理

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
