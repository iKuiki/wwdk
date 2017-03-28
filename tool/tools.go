package tool

import (
	"strconv"
	"time"
)

func GetWxTimeStamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)[:13]
}
