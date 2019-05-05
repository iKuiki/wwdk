# 联系人类型解析

| 字段             | 用途                                                |
| :--------------- | :-------------------------------------------------- |
| Alias            | 微信号                                              |
| AppAccountFlag   | 暂未知                                              |
| AttrStatus       | 暂未知<br/>(好友为一个很大的数字，公众号、群聊为0)  |
| ChatRoomID       | 暂未知<br/>(全为0)                                  |
| City             | 用户资料的所在城市                                  |
| ContactFlag      | 联系人标记（另附表）                                |
| DisplayName      | 暂未知<br/>(全为空)                                 |
| EncryChatRoomID  | 暂未知<br/>(全为空)                                 |
| HeadImgURL       | 用户头像地址                                        |
| HideInputBarFlag | 暂未知<br/>(全为0)                                  |
| IsOwner          | 是否是该群的群主<br/>是群主为1，其余为0，只对群有效 |
| KeyWord          | 微信自动归纳的关键字                                |
| MemberCount      | 暂未知<br/>(全为0)                                  |
| MemberList       | 暂未知<br/>(全为0)                                  |
| NickName         | 用户昵称                                            |
| OwnerUin         | 暂未知<br/>(全为0)                                  |
| PYInitial        | 用户昵称的拼音首字母                                |
| PYQuanPin        | 用户昵称全拼                                        |
| Province         | 用户资料的所在省份                                  |
| RemarkName       | 用户备注名称                                        |
| RemarkPYInitial  | 备注名称的拼音首字母                                |
| RemarkPYQuanPin  | 备注名称的全拼                                      |
| Sex              | 性别<br/>1男2女                                     |
| Signature        | 用户的个性签名                                      |
| SnsFlag          | 暂未知<br/>(好友部分为17部分为49)                   |
| StarFriend       | 是否星标好友<br/>1为星标，0为普通                   |
| Statues          | 暂未知                                              |
| Uin              | 暂未知<br/>(全为0)                                  |
| UniFriend        | 暂未知<br/>(全为0)                                  |
| UserName         | 临时用户识别码，消息都通过这个识别码发送和接收      |
| VerifyFlag       | 用户类别？<br/>公众号为24，微信团队为56             |

PS: 用户头像为uri，需要带上主机名wx2.qq.com并以https请求，并且要在url末拼接上skey作为查询参数，以及需要提交cookie

#### ContactFlag联系人标记解析
ContactFlag Value的值为对应选项相加的和

| 选项                                                            | 值    |
| :-------------------------------------------------------------- | :---- |
| 普通好友<br/>(基础值，应该是1和2相加，大家都有，不知是什么含义) | 3     |
| 星标好友                                                        | 64    |
| 不让他看我的朋友圈                                              | 256   |
| 免打扰                                                          | 512   |
| 置顶好友                                                        | 2048  |
| 不看他的朋友圈                                                  | 65536 |

例如一个置顶好友，其ContactFlag就为3+2048=2051; 一个置顶的星标好友，其ContactFlag就为3+64+2048=2115
拿到ContactFlag之后要还原其选了什么选项，只需将ContactFlag与上述值按位与即可
