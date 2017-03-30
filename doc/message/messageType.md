# 微信信息类型解析

| TypeId | Description |
|--------|---------|
| 1 | 文字消息 |
| 3 | 图片消息 |
| 34 | 音频消息 |
| 42 | 名片 |
| 43 | 小视频消息 |
| 47 | 动画表情 |
| 49 | 链接消息类型，已知有转账、开始共享实时位置、合并转发聊天记录 |
| 10000 | 拓展消息类型，已知有红包、停止共享实时位置、AA收款 |
| 10002 | 撤回消息 |

### Type 3 Detail
| SpecialField | UseFor |
|--------------|--------|
| ImgHeight | 120 |
| ImgStatus | 2 |
| ImgWidth | 67 |
| Status | 3 |

### Type 49 Detail

#### 收到文件
| SpecialField | UseFor |
|--------------|--------|
| AppMsgType| 6 |
| FileName | 记录文件名 |
| FileSize | 文件大小 |
| MediaId | - |
| Status | 3 |
| ImgStatus | 1 |

#### 收到转账
| SpecialField | UseFor |
|--------------|--------|
| AppMsgType | 2000 |
| FileName | ??? |
| Status | 3 |
| ImgStatus | 1 |

#### 开始共享实时位置
| SpecialField | UseFor |
|--------------|--------|
| AppMsgType | 17 |
| FileName | ??? |
| Status | 3 |
| ImgStatus | 1 |

#### 合并转发聊天记录
| SpecialField | UseFor |
|--------------|--------|
| AppMsgType | 0 |
| Status | 3 |
| ImgStatus | 1 |


### Type 10000 Detail

#### 停止共享实时位置
| SpecialField | UseFor |
|--------------|--------|
| AppMsgType | 0 |
| Content | ??? |
| Status | 4 |

#### 红包
| SpecialField | UseFor |
|--------------|--------|
| AppMsgType | 0 |
| Content | ??? |
| Status | 4 |

#### AA收款
| SpecialField | UseFor |
|--------------|--------|
| AppMsgType | 0 |
| Content | ??? |
| Status | 4 |
| ImgStatus | 1 |

### Type 10002 Detail
#### 撤回消息
| SpecialField | UseFor |
|--------------|--------|
| Status | 4 |
| ImgStatus | 1 |
