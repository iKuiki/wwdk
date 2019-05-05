package tool

import (
	"math/rand"
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
