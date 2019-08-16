package wwdk

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ikuiki/wwdk/api"

	"github.com/pkg/errors"

	"github.com/ikuiki/wwdk/datastruct"
)

// SyncStatus 同步状态
type SyncStatus int32

const (
	// SyncStatusModifyContact 同步状态：有联系人变更
	SyncStatusModifyContact SyncStatus = 1
	// SyncStatusNewMessage 同步状态：有新信息
	SyncStatusNewMessage SyncStatus = 2
	// SyncStatusPanic 致命错误，sync进程退出
	SyncStatusPanic SyncStatus = -1
	// SyncStatusErrorOccurred 非致命性错误发生，具体错误请参考Msg
	SyncStatusErrorOccurred SyncStatus = -2
)

// SyncChannelItem 同步管道通信要素
type SyncChannelItem struct {
	Code    SyncStatus          // 同步状态
	Contact *datastruct.Contact // 联系人（如果同步状态是有联系人变更则有
	Message *datastruct.Message // 新信息（如果同步状态是有新信息则有
	Err     error               // 错误（如有发生
	// Msg     string              // 其他附带信息
}

// StartServe 启动消息同步服务
func (wxwb *WechatWeb) StartServe(syncChannel chan<- SyncChannelItem) {
	wxwb.syncChannel = syncChannel
	go func() {
		// 方法结束时关闭channel
		defer close(syncChannel)
		getMessage := func() {
			modContacts, delContacts, addMessage, body, err := wxwb.api.WebwxSync()
			if err != nil {
				wxwb.captureException(err, "WebwxSync fatal", sentry.LevelError, extraData{"body", string(body)})
				wxwb.logger.Infof("WebwxSync error: %s\n", err.Error())
				return
			}
			// 处理新增联系人
			wxwb.contactManager.SetFriend(modContacts...)
			// 并发送给channel
			for _, contact := range modContacts {
				wxwb.runInfo.ContactModifyCount++
				wxwb.logger.Infof("Modify contact: %s\n", contact.NickName)
				syncChannel <- SyncChannelItem{
					Code:    SyncStatusModifyContact,
					Contact: &contact,
				}
			}
			// 处理删除的联系人
			for _, delContact := range delContacts {
				wxwb.contactManager.DelFriend(delContact.UserName)
			}
			// 新消息
			for _, msg := range addMessage {
				if msg.MsgType == datastruct.RevokeMsg {
					wxwb.runInfo.MessageRevokeCount++
				} else {
					wxwb.runInfo.MessageCount++
				}
				err = wxwb.messageProcesser(&msg, syncChannel)
				if err != nil {
					wxwb.captureException(err, "MessageProcesser error", sentry.LevelError, extraData{"msg", msg})
					wxwb.logger.Infof("MessageProcesser error: %+v\n", err)
					syncChannel <- SyncChannelItem{
						Code: SyncStatusErrorOccurred,
						Err:  err,
					}
					continue
				}
			}
		}
		for {
			isBreaked := func() (isBreaked bool) {
				defer func() {
					if r := recover(); r != nil {
						var eErr error
						if err, ok := r.(error); ok {
							eErr = err
						}
						wxwb.captureException(eErr, "Sync loop panic", sentry.LevelError, extraData{"panicItem", r})
						wxwb.logger.Infof("Recovered in Sync loop: %v\n", r)
						wxwb.runInfo.PanicCount++
						syncChannel <- SyncChannelItem{
							Code: SyncStatusErrorOccurred,
							Err:  errors.Errorf("recovered panic: %v", r),
						}
					}
				}()
				_, selector, body, err := wxwb.api.SyncCheck()
				if err != nil {
					if err == api.ErrLogout {
						wxwb.logger.Info("User has logout web wechat, exit...\n")
						// 关闭登录消息通知管道
						wxwb.api.CloseLoginModifyNotifyChan()
						// 清空登录记录
						wxwb.loginStorer.Truncate()
						// 对外发送通知
						syncChannel <- SyncChannelItem{
							Code: SyncStatusPanic,
							Err:  errors.New("Err1101: user has logout"),
						}
						return true
					}
					wxwb.captureException(err, "SyncCheck fatal", sentry.LevelError, extraData{"body", string(body)})
					wxwb.logger.Infof("SyncCheck error: %s\n", err.Error())
					syncChannel <- SyncChannelItem{
						Code: SyncStatusErrorOccurred,
						Err:  err,
					}
					return false
				}
				// wxwb.logger.Infof("selector: %v\n", selector)
				switch selector {
				case "0":
					// wxwb.logger.Info("SyncCheck 0\n")
					// normal
					// wxwb.logger.Info("no new message\n")
				case "6":
					wxwb.logger.Info("selector is 6\n")
					getMessage()
				case "7":
					wxwb.logger.Info("selector is 7\n")
					getMessage()
				case "1":
					wxwb.logger.Info("selector is 1\n")
					getMessage()
				case "3":
					wxwb.logger.Info("selector is 3\n")
					getMessage()
				case "4":
					wxwb.logger.Info("selector is 4\n")
					getMessage()
				case "5":
					wxwb.logger.Info("selector is 5\n")
					getMessage()
				case "2":
					// wxwb.logger.Info("SyncCheck 2\n")
					getMessage()
				default:
					wxwb.captureException(nil, "SyncCheck Unknow selector", sentry.LevelWarning, extraData{"selector", selector})
					wxwb.logger.Infof("SyncCheck Unknow selector: %s\n", selector)
					syncChannel <- SyncChannelItem{
						Code: SyncStatusErrorOccurred,
						Err:  errors.Errorf("syncCheck unknow selector: %s", selector),
					}
				}
				wxwb.runInfo.SyncCount++
				time.Sleep(1000 * time.Millisecond)
				return false
			}()
			if isBreaked {
				break
			}
		}
	}()
}
