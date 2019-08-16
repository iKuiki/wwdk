package wwdk_test

import (
	"fmt"
	"os"

	"github.com/ikuiki/wwdk"
	"github.com/mdp/qrterminal"
)

var webwx *wwdk.WechatWeb

func init() {
	var err error
	webwx, err = wwdk.NewWechatWeb()
	assertErrIsNil(err)
	loginChan := make(chan wwdk.LoginChannelItem)
	webwx.Login(loginChan)
	for item := range loginChan {
		switch item.Code {
		case wwdk.LoginStatusWaitForScan:
			qrterminal.Generate(item.Msg, qrterminal.L, os.Stdout)
		case wwdk.LoginStatusScanedWaitForLogin:
			fmt.Println("scaned")
		case wwdk.LoginStatusScanedFinish:
			fmt.Println("accepted")
		case wwdk.LoginStatusGotCookie:
			fmt.Println("got cookie")
		case wwdk.LoginStatusInitFinish:
			fmt.Println("init finish")
		case wwdk.LoginStatusGotContact:
			fmt.Println("got contact")
		case wwdk.LoginStatusBatchGotContact:
			fmt.Println("got batch contact")
			break
		case wwdk.LoginStatusErrorOccurred:
			panic(fmt.Sprintf("WxWeb Login error: %+v", item.Err))
		default:
			fmt.Printf("unknown code: %+v", item)
		}
	}
	user, err := webwx.GetUser()
	assertErrIsNil(err)
	fmt.Println(user.NickName, " Login successful")
}

func assertErrIsNil(err error) {
	if err != nil {
		panic(err)
	}
}
