package sub

import (
	"testing"
)

// func TestListDir(t *testing.T) {
// 	var ret []string
// 	ret = ListDir("./", "/")
// 	t.Log(ret)
// }

func TestWalkDir(t *testing.T) {
	ret, _ := WalkDir("../", ".go")
	t.Log(ret)
}

// func TestCMD(t *testing.T) {
// 	// CMD("E:/ItemForGo/src/github.com/shaalx/echo/oauth2")
// 	CMD2()
// }