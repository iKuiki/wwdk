package tool

import (
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// GetWxTimeStamp 获取web微信使用的13位时间戳
func GetWxTimeStamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)[:13]
}

// GetRandomStringFromNum 获取指定长度的随机数字
func GetRandomStringFromNum(length int) string {
	bytes := []byte("0123456789")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// WriteToFile 将io流写入文件
func WriteToFile(filename string, data io.ReadCloser) (n int, err error) {
	f, err := os.Create(filename)
	if err != nil {
		return 0, errors.New("create " + filename + " error: " + err.Error())
	}
	d, err := ioutil.ReadAll(data)
	if err != nil {
		return 0, errors.New("Read io.ReadCloser error: " + err.Error())
	}
	n, err = f.Write(d)
	if err != nil {
		return 0, errors.New("Write to file error: " + err.Error())
	}
	return n, nil
}
