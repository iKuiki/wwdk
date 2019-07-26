# 微信信息类型解析

| TypeId | Description                    |
|:-------|:-------------------------------|
| 1      | 文字消息                           |
| 3      | 图片消息                           |
| 34     | 音频消息                           |
| 42     | 名片                             |
| 43     | 小视频消息                          |
| 47     | 动画表情                           |
| 49     | app消息类型，已知有转账、开始共享实时位置、合并转发聊天记录、收到文件 |
| 51     | 手机客户端切换聊天对象                    |
| 10000  | 拓展消息类型，已知有红包、停止共享实时位置、AA收款     |
| 10002  | 撤回消息                           |

## Type 3 图片消息

| SpecialField | UseFor                  |
|:-------------|:------------------------|
| ImgStatus    | 2                       |
| Status       | 3                       |

## Type 34 音频消息

| SpecialField | UseFor                  |
|:-------------|:------------------------|
| ImgStatus    | 1                       |
| Status       | 3                       |

## Type 43 小视频消息

| SpecialField | UseFor                |
|:-------------|:----------------------|
| ImgStatus    | 1                     |
| Status       | 3                     |

## Type47 动画表情

| SpecialField | UseFor             |
|:-------------|:-------------------|
| Content      | html转义的xml，记录了表情地址 |
| ImgStatus    | 2                  |
| Status       | 3                  |

## Type 49 程序消息

### AppMsgType 6 收到文件

| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 6      |
| FileName     | 记录文件名  |
| Status       | 3      |
| ImgStatus    | 1      |

### AppMsgType 2000 收到转账

| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 2000   |
| FileName     | ???    |
| Status       | 3      |
| ImgStatus    | 1      |

### AppMsgType 17 开始共享实时位置

| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 17     |
| FileName     | ???    |
| Status       | 3      |
| ImgStatus    | 1      |

### AppMsgType 0 合并转发聊天记录

| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 0      |
| Status       | 3      |
| ImgStatus    | 1      |

## Type 10000 扩展消息类型

### AppMsgType 0 停止共享实时位置

| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 0      |
| Content      | ???    |
| Status       | 4      |

### AppMsgType 0 红包

| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 0      |
| Content      | ???    |
| Status       | 4      |

### AppMsgType 0 AA收款

| SpecialField | UseFor |
|:-------------|:-------|
| AppMsgType   | 0      |
| Content      | ???    |
| Status       | 4      |
| ImgStatus    | 1      |

## Type 10002 撤回消息

| SpecialField | UseFor |
|:-------------|:-------|
| Status       | 4      |
| ImgStatus    | 1      |
