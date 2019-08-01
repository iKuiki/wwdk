package api_test

import (
	"fmt"
	"github.com/ikuiki/wwdk/api"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"github.com/mdp/qrterminal"
	"os"
	"strings"
)

var (
	client api.WechatwebAPI
	// wxInit获取到的user
	user datastruct.User
	// wxInit获取到的联系人列表
	contactMap map[string]datastruct.Contact = make(map[string]datastruct.Contact)
	// 修改联系人时传入的管道
	modContactChan chan datastruct.Contact = make(chan datastruct.Contact, 1000)
	// 删除联系人时传入的管道
	delContactChan chan datastruct.WebwxSyncRespondDelContactListItem = make(chan datastruct.WebwxSyncRespondDelContactListItem, 100)
	// 有新消息时传入的管道
	addMessageChan chan datastruct.Message = make(chan datastruct.Message, 10000)
)

func init() {
	// 获取一个client备用
	client = api.MustNewWechatwebAPI()
	login(client)
	sync(client, modContactChan,
		delContactChan,
		addMessageChan)
}

func checkErrorIsNil(err error) {
	if err != nil {
		panic(err)
	}
}

// 执行登陆逻辑
func login(client api.WechatwebAPI) {
	uuid, _, err := client.JsLogin()
	checkErrorIsNil(err)
	qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)
	_, _, redirectURL, _, err := client.Login(uuid, "1")
	checkErrorIsNil(err)
	for redirectURL == "" {
		var code string
		code, _, redirectURL, _, err = client.Login(uuid, "0")
		if !strings.HasPrefix(code, "2") {
			panic("code is " + code)
		}
		checkErrorIsNil(err)
	}
	_, err = client.WebwxNewLoginPage(redirectURL)
	checkErrorIsNil(err)
	userPtr, contactList, _, err := client.WebwxInit()
	checkErrorIsNil(err)
	user = *userPtr
	for _, contact := range contactList {
		contactMap[contact.UserName] = contact
	}
	contactList, _, err = client.GetContact()
	for _, contact := range contactList {
		contactMap[contact.UserName] = contact
	}
	fmt.Println("user " + user.NickName + " login success")
}

// 一个微型服务来维持sync与联系人列表
func sync(client api.WechatwebAPI,
	modContactChan chan datastruct.Contact,
	delContactChan chan datastruct.WebwxSyncRespondDelContactListItem,
	addMessageChan chan datastruct.Message,
) {
	go func() {
		for {
			_, selector, _, err := client.SyncCheck()
			checkErrorIsNil(err)
			if selector != "0" {
				modContacts, delContacts, addMessages, _, err := client.WebwxSync()
				checkErrorIsNil(err)
				go func() {
					for _, contact := range modContacts {
						contactMap[contact.UserName] = contact
						modContactChan <- contact
					}
				}()
				go func() {
					for _, item := range delContacts {
						delete(contactMap, item.UserName)
						delContactChan <- item
					}
				}()
				go func() {
					for _, msg := range addMessages {
						if _, ok := contactMap[msg.FromUserName]; !ok {
							contactList, _, err := client.BatchGetContact([]datastruct.BatchGetContactRequestListItem{
								datastruct.BatchGetContactRequestListItem{UserName: msg.FromUserName},
							})
							checkErrorIsNil(err)
							for _, contact := range contactList {
								contactMap[contact.UserName] = contact
							}
						}
						addMessageChan <- msg
					}
				}()
			}
		}
	}()
}

// 需要操作联系人时，为了不影响他人，需要先生成一个随机数，让对应联系人发过来确认
// fnName是当前测试名称
// isChatroom指定了是否是需要群
func getTestContact(fnName string, isChatroom bool) (contact datastruct.Contact, skip bool) {
	validCode := tool.GetRandomStringFromNum(4)
	target := "contact"
	if isChatroom {
		target = "chatroom"
	}
	fmt.Printf("need to select a %s to test\n", target)
	fmt.Printf("if you want to Run test function %s please send [run %s] via a %s\n",
		fnName,
		validCode,
		target,
	)
	fmt.Printf("if you want to Skip test function %s please send [skip %s]\n",
		fnName,
		validCode,
	)
	for {
		msg := <-addMessageChan
		if msg.GetContent() == "skip "+validCode {
			return contact, true
		} else if msg.GetContent() == "run "+validCode {
			contact = contactMap[msg.FromUserName]
			if isChatroom == contact.IsChatroom() {
				return contactMap[msg.FromUserName], false
			} else {
				fmt.Printf("please use a %s to send again", target)
			}
		} else {
			fmt.Printf("content [%s] mismatch\n", msg.GetContent())
		}
	}
}
