package bkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ListDir(dir, suf string) []string {
	dirs, err := ioutil.ReadDir(suf + dir)
	if nil != err {
		return nil
	}
	var ret []string
	for _, it := range dirs {
		if it.IsDir() {
			sub_dirs := ListDir(it.Name(), suf+"/")
			ret = append(ret, sub_dirs...)
			fmt.Println(sub_dirs)
		} else {
			fmt.Println(it.Name())
			ret = append(ret, it.Name())
		}
		fmt.Println(it.Name())
	}
	return ret
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (dirs []string, err error) {
	dirs = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}

		if !fi.IsDir() {
			return nil
		}
		fmt.Println(fi.Name())
		if strings.EqualFold(fi.Name(), ".git") {
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			dirs = append(dirs, filename)
		}

		return nil
	})

	return dirs, err
}

func CMD() {
	output, err := exec.Command("go", "install").Output()
	// err := cmd.Run()
	if nil != err {
		fmt.Println(err)
	}
	fmt.Println(string(output))
}
