package wwdk

// 此文件中存储的为WechatWeb读写登录凭据的方法

import (
	"encoding/json"

	"github.com/getsentry/sentry-go"
	"github.com/ikuiki/wwdk/api"
	"github.com/ikuiki/wwdk/contactmgr"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/runinfo"
	"github.com/pkg/errors"
)

// storeLoginInfo 用于储存的登录信息
type storeLoginInfo struct {
	APIMarshaled []byte
	RunInfo      runinfo.WechatRunInfo // 运行统计信息
	User         *datastruct.User      // 用户信息
	ContactList  []datastruct.Contact
	FriendList   []string
}

// 重置登录信息
func (wxwb *WechatWeb) resetLoginInfo() (err error) {
	// 此处如果panic，不应当阻止其传播
	// 如果发生panic，是影响到逻辑的panic
	if wxwb.loginStorer != nil {
		wxwb.loginStorer.Truncate()
	}
	wxwb.api = api.MustNewWechatwebAPI()
	// 重置runInfo
	wxwb.runInfo = runinfo.WechatRunInfo{
		StartAt: wxwb.runInfo.StartAt,
	}
	// 切记也要重置用户信息与联系人啊
	wxwb.contactManager, err = contactmgr.NewManager()
	return nil
}

// 往storer中写入信息
func (wxwb *WechatWeb) writeLoginInfo() (err error) {
	defer func() {
		if r := recover(); r != nil {
			var eErr error
			if err, ok := r.(error); ok {
				eErr = err
			}
			wxwb.captureException(eErr, "WriteLoginInfo panic", sentry.LevelError, extraData{"panicItem", r})
			wxwb.logger.Infof("Recovered in writeLoginInfo: %v\n", r)
			wxwb.runInfo.PanicCount++
			err = errors.Errorf("panic recovered: %+v", r)
		}
	}()
	if wxwb.loginStorer != nil {
		apiMarshaled, err := wxwb.api.Marshal()
		if err != nil {
			return errors.WithStack(err)
		}
		contactList, err := wxwb.contactManager.GetAllContacts()
		if err != nil {
			return err
		}
		friendList, err := wxwb.contactManager.GetFriendList()
		if err != nil {
			return err
		}
		var friendUserNames []string
		for _, f := range friendList {
			friendUserNames = append(friendUserNames, f.UserName)
		}
		storeInfo := storeLoginInfo{
			APIMarshaled: apiMarshaled,
			User:         wxwb.user,
			ContactList:  contactList,
			FriendList:   friendUserNames,
			RunInfo:      wxwb.runInfo,
		}
		data, err := json.Marshal(storeInfo)
		if err != nil {
			return errors.WithStack(err)
		}
		err = wxwb.loginStorer.Write(data)
		return errors.WithStack(err)
	}
	return nil
}

// 从storer中读取信息
// 返回是否成功读取到信息
func (wxwb *WechatWeb) readLoginInfo() (readed bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			var eErr error
			if err, ok := r.(error); ok {
				eErr = err
			}
			wxwb.captureException(eErr, "ReadLoginInfo panic", sentry.LevelError, extraData{"panicItem", r})
			wxwb.logger.Infof("Recovered in readLoginInfo: %v\n", r)
			wxwb.runInfo.PanicCount++
			err = errors.Errorf("panic recovered: %+v", r)
		}
	}()
	if wxwb.loginStorer != nil {
		data, err := wxwb.loginStorer.Read()
		if err != nil {
			return false, errors.WithStack(err)
		}
		if string(data) == "" {
			// 登陆信息为空，不继续接下来的流程了
			return false, nil
		}
		var storeInfo storeLoginInfo
		err = json.Unmarshal(data, &storeInfo)
		if err != nil {
			return false, errors.WithStack(err)
		}
		err = wxwb.api.Unmarshal(storeInfo.APIMarshaled)
		if err != nil {
			if err == api.ErrEmptyLoginInfo {
				// api恢复时其内部关键信息为空
				// 判定为未读取到登录信息
				return false, nil
			}
			return false, errors.WithStack(err)
		}
		if storeInfo.User == nil {
			// 只要userInfo中的User为空
			// 则判定为未读取到登陆信息
			return false, nil
		}
		// 认为读取到了登陆信息，则开始还原
		{
			// 先暂存StartAt，对StartAt不做覆盖
			started := wxwb.runInfo.StartAt
			wxwb.runInfo = storeInfo.RunInfo
			// 还原startat
			wxwb.runInfo.StartAt = started
		}
		wxwb.user = storeInfo.User
		// 还原ContactList
		wxwb.contactManager.SetContact(storeInfo.ContactList...)
		// 还原FriendList
		wxwb.contactManager.Contact2Friend(storeInfo.FriendList...)
		// 还原完成
		return true, nil
	}
	return false, nil
}
