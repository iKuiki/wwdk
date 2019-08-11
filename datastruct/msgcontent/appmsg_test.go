package msgcontent_test

import (
	"bytes"
	"encoding/xml"
	"github.com/ikuiki/wwdk/datastruct/msgcontent"
	"io"
	"os"
	"testing"
)

// TestFileMsg 测试文件类型的content反序列化、序列化
func TestFileMsg(t *testing.T) {
	f, err := os.Open("fileMsgExample.xml")
	if err != nil {
		t.Fatalf("open example file fail: %v", err)
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	if err != nil {
		t.Fatalf("read example file fail: %v", err)
	}
	var appMsgContent msgcontent.AppMsgContent
	err = xml.Unmarshal(buf.Bytes(), &appMsgContent)
	if err != nil {
		t.Fatalf("unmarshal appMsgContent fail: %v", err)
	}
	shal, err := xml.Marshal(appMsgContent)
	if err != nil {
		t.Fatalf("marshal appMsgContent fail: %v", err)
	}
	if string(shal) != string(buf.Bytes()) {
		t.Fatalf("readed data \n%s\n diff with marshaled data \n%s\n",
			string(buf.Bytes()), string(shal))
	}
}
