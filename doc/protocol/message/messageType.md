# 微信信息类型解析

| TypeId | Description                    |
|:-------|:-------------------------------|
| 1      | 文字消息                           |
| 3      | 图片消息                           |
| 34     | 音频消息                           |
| 42     | 名片                             |
| 43     | 小视频消息                          |
| 47     | 动画表情                           |
| 49     | 链接消息类型，已知有转账、开始共享实时位置、合并转发聊天记录 |
| 51     | 手机客户端切换聊天对象                    |
| 10000  | 拓展消息类型，已知有红包、停止共享实时位置、AA收款     |
| 10002  | 撤回消息                           |

### Type 3 Detail
| SpecialField | UseFor                  |
|:-------------|:------------------------|
| Content      | html转译的xml，记录了图片转发所需的信息 |
| ImgHeight    | 120                     |
| ImgStatus    | 2                       |
| ImgWidth     | 67                      |
| Status       | 3                       |

若需要获取图片本身，则需要访问以下地址获得：

| GetImgMessage | value                                                                      |
|:--------------|:---------------------------------------------------------------------------|
| Url           | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetmsgimg                      |
| Method        | Get                                                                        |
| Param         | MsgID(webwxsync接口拿到的)<br>skey(wxinit接口拿到的)<br>type:代表图片大小，已知的可选项为slave、big |

### Type 34 Detail
| SpecialField | UseFor                  |
|:-------------|:------------------------|
| Content      | html转译的xml，记录了语音转发所需的信息 |
| ImgHeight    | 0                       |
| ImgStatus    | 1                       |
| ImgWidth     | 0                       |
| Status       | 3                       |

若需要获取语音本身，则需要访问以下地址获得：

| GetVoiceMessage | value                                                |
|:----------------|:-----------------------------------------------------|
| Url             | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvoice |
| Method          | Get                                                  |
| Header          | Range: bytes=0-                                      |
| Param           | msgid(webwxsync接口拿到的)<br>skey(wxinit接口拿到的)           |

### Type 43 Detail
| SpecialField | UseFor                |
|:-------------|:----------------------|
| Content      | html转译的xml，记录了视频的相关信息 |
| ImgHeight    | 540                   |
| ImgStatus    | 1                     |
| ImgWidth     | 944                   |
| Status       | 3                     |

若需要获取视频本身，则需要访问以下地址获得：

| GetVideoMessage | value                                                |
|:----------------|:-----------------------------------------------------|
| Url             | https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvideo |
| Method          | Get                                                  |
| Header          | Range: bytes=0-                                      |
| Param           | msgid(webwxsync接口拿到的)<br>skey(wxinit接口拿到的)           |

### Type47 Detail
| SpecialField | UseFor             |
|:-------------|:-------------------|
| Content      | html转义的xml，记录了表情地址 |
| ImgHeight    | 100                |
| ImgWidth     | 100                |
| ImgStatus    | 2                  |
| Status       | 3                  |

### Type 49 Detail

#### 收到文件
| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 6      |
| FileName     | 记录文件名  |
| FileSize     | 文件大小   |
| MediaId      | -      |
| Status       | 3      |
| ImgStatus    | 1      |

#### 收到转账
| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 2000   |
| FileName     | ???    |
| Status       | 3      |
| ImgStatus    | 1      |

#### 开始共享实时位置
| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 17     |
| FileName     | ???    |
| Status       | 3      |
| ImgStatus    | 1      |

#### 合并转发聊天记录
| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 0      |
| Status       | 3      |
| ImgStatus    | 1      |


### Type 10000 Detail

#### 停止共享实时位置
| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 0      |
| Content      | ???    |
| Status       | 4      |

#### 红包
| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 0      |
| Content      | ???    |
| Status       | 4      |

#### AA收款
| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 0      |
| Content      | ???    |
| Status       | 4      |
| ImgStatus    | 1      |

### Type 10002 Detail
#### 撤回消息
| SpecialField | UseFor |
|:-------------|:-------|
| Status       | 4      |
| ImgStatus    | 1      |
