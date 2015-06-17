package bkg

import (
	"testing"
)

func TestListDir(t *testing.T) {
	ret:= ListDir("$GOPATH/","")
	t.Log(ret)
}

func TestWalkDir(t *testing.T) {
	ret,_:= WalkDir("../",".go")
	t.Log(ret)
}
