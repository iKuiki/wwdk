package tool

import (
	"strings"
)

// AnalysisWxWindowRespond 	解析返回值中由js代码构成的参数
func AnalysisWxWindowRespond(respond string) (ret map[string]string) {
	ret = make(map[string]string)
	arr := strings.Split(respond, ";")
	for _, a := range arr {
		index := strings.Index(a, "=")
		if index > 0 && len(a) > index+1 {
			k := strings.TrimSpace(a[:index])
			v := strings.TrimSpace(a[index+1:])
			v = strings.TrimPrefix(v, `"`)
			v = strings.TrimSuffix(v, `"`)
			ret[k] = v
		}
	}
	return ret
}
