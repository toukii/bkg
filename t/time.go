package t

import (
	"time"
)

func LocNow(location string) *time.Time {
	loc, err := time.LoadLocation(location) //"Asia/Shanghai"
	if nil != err {
		return nil
	}
	tp := time.Unix(time.Now().Unix(), 0).In(loc)
	return &tp
}
