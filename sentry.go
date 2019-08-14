package wwdk

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
)

type extraData struct {
	Key   string
	Value interface{}
}

// captureException 捕获异常
// @param err 错误本身，最好是带stracktrace的err
// @param errType 错误类型，错误的名字，会作为sentry的标题
// @param level 错误严重等级
// @param extras 附加数据，如果有附加数据会放进event的extras里
func (wwdk *WechatWeb) captureException(err error, errType string, level sentry.Level, extras ...extraData) {
	event := sentry.NewEvent()
	event.Level = level
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
	if wwdk.user != nil {
		event.User.ID = wwdk.user.UserName
		event.User.Username = wwdk.user.NickName
		event.Extra["wwdk_login_at"] = wwdk.runInfo.LoginAt
	}
	event.Extra["wwdk_StartAt"] = wwdk.runInfo.StartAt
	event.Extra["wwdk_SyncCount"] = wwdk.runInfo.SyncCount
	event.Extra["wwdk_ContactModifyCount"] = wwdk.runInfo.ContactModifyCount
	event.Extra["wwdk_MessageCount"] = wwdk.runInfo.MessageCount
	event.Extra["wwdk_PanicCount"] = wwdk.runInfo.PanicCount
	for _, extra := range extras {
		event.Extra[extra.Key] = extra.Value
	}
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
