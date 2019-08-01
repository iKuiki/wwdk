package wwdk

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"runtime"
	"strings"
)

// captureException 捕获异常
// @param err 错误本身，最好是带stracktrace的err
// @param errType 错误类型，错误的名字，会作为sentry的标题
// @param operation 当前的操作，例如：wwdk.NewWechatWeb、wwdk.api.JsLogin
// @param level 错误严重等级
func (wwdk *WechatWeb) captureException(err error, errType, operation string, level sentry.Level) {
	event := sentry.NewEvent()
	event.Level = sentry.LevelWarning
	if err == nil {
		event.Message = fmt.Sprintf("Called %s with nil value", callerFunctionName())
		event.Exception = []sentry.Exception{{
			Value:      "",
			Type:       errType,
			Stacktrace: sentry.NewStacktrace(),
		}}
	} else {
		stacktrace := sentry.ExtractStacktrace(err)
		if stacktrace == nil {
			stacktrace = sentry.NewStacktrace()
		}
		event.Exception = []sentry.Exception{{
			Value:      err.Error(),
			Type:       errType,
			Stacktrace: stacktrace,
		}}
	}
	event.Tags["wwdk_operation"] = operation
	if wwdk.userInfo.user != nil {
		event.User.ID = wwdk.userInfo.user.UserName
		event.User.Username = wwdk.userInfo.user.NickName
		event.Extra["wwdk_login_at"] = wwdk.runInfo.LoginAt
	}
	event.Extra["wwdk_operation"] = operation
	event.Extra["wwdk_StartAt"] = wwdk.runInfo.StartAt
	event.Extra["wwdk_SyncCount"] = wwdk.runInfo.SyncCount
	event.Extra["wwdk_ContactModifyCount"] = wwdk.runInfo.ContactModifyCount
	event.Extra["wwdk_MessageCount"] = wwdk.runInfo.MessageCount
	event.Extra["wwdk_PanicCount"] = wwdk.runInfo.PanicCount
	wwdk.sentryHub.CaptureEvent(event)
}

func callerFunctionName() string {
	pcs := make([]uintptr, 1)
	runtime.Callers(3, pcs)
	callersFrames := runtime.CallersFrames(pcs)
	callerFrame, _ := callersFrames.Next()
	_, function := deconstructFunctionName(callerFrame.Function)
	return function
}

// Transform `runtime/debug.*T·ptrmethod` into `{ module: runtime/debug, function: *T.ptrmethod }`
func deconstructFunctionName(name string) (module string, function string) {
	if idx := strings.LastIndex(name, "."); idx != -1 {
		module = name[:idx]
		function = name[idx+1:]
	}
	function = strings.Replace(function, "·", ".", -1)
	return module, function
}
