package storer_test

import (
	"github.com/ikuiki/wechat-web/storer"
	"os"
	"testing"
)

func TestFileStorer(t *testing.T) {
	testStorerFilePath := "testFileStorer.txt"
	fStorer := storer.MustNewFileStorer(testStorerFilePath)
	dataA := "aaa"
	dataB := "bbb"
	// 先写入一次a
	err := fStorer.WriterString(dataA)
	if err != nil {
		t.Fatalf("storer.WriterString(dataA) error: %v", err)
	}
	// 测试读取a两次，验证读取一次后文件游标是否归零
	readData, err := fStorer.ReadString()
	if err != nil {
		t.Fatalf("storer.ReadString() error: %v", err)
	}
	if readData != dataA {
		t.Fatalf("data readed at file [%s] diff with dataA [%s]", readData, dataA)
	}
	// 第二次读取
	readData, err = fStorer.ReadString()
	if err != nil {
		t.Fatalf("storer.ReadString() error: %v", err)
	}
	if readData != dataA {
		t.Fatalf("data readed at file [%s] diff with dataA [%s]", readData, dataA)
	}
	// 再写入一次a
	err = fStorer.WriterString(dataA)
	if err != nil {
		t.Fatalf("storer.WriterString(dataA) error: %v", err)
	}
	// 然后写入一次B
	err = fStorer.WriterString(dataB)
	if err != nil {
		t.Fatalf("storer.WriterString(dataB) error: %v", err)
	}
	// 测试读取,验证结果是否为b
	readData, err = fStorer.ReadString()
	if err != nil {
		t.Fatalf("storer.ReadString() error: %v", err)
	}
	if readData != dataB {
		t.Fatalf("data readed at file [%s] diff with dataB [%s]", readData, dataB)
	}
	os.Remove(testStorerFilePath)
	t.Log("test pass, remove test tmp file")
}
