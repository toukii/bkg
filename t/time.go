package t

import (
	"time"
)

func Now() string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if nil != err {
		return ""
	}
	tp := time.Unix(time.Now().Unix(), 0).In(loc)
	return tp.String()
}
