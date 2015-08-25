package main

import (
	"fmt"
	"github.com/everfore/exc"
	"github.com/shaalx/goutils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Info struct {
	dir  string
	ok   bool
	info string
}

func (i Info) String() string {
	if i.ok {
		return fmt.Sprintf("[SUCCESS]: %s", i.dir)
	}
	length := len(i.info)
	if length > 150 {
		length = 150
	}
	return fmt.Sprintf("[FAIL]: %s \n%s [more]...", i.dir, i.info[:length])
}

func NewInfo(dir string, ok bool, info string) *Info {
	return &Info{
		dir:  dir,
		ok:   ok,
		info: info,
	}
}

var (
	command     *exc.CMD
	installInfo chan *Info
)

func init() {
	installInfo = make(chan *Info, 10)
	command = exc.NewCMD("go install")

}

func searchDir(dir string) {
	file, err := os.Open(dir)
	if exc.Checkerr(err) {
		return
	}
	subdirs, err := file.Readdir(-1)
	if exc.Checkerr(err) {
		return
	}
	excuted := false
	for _, it := range subdirs {
		if strings.EqualFold(it.Name(), ".git") {
			continue
		}
		if it.IsDir() {
			go searchDir(filepath.Join(dir, it.Name()))
		}
		if strings.Contains(it.Name(), ".go") && !excuted {
			b, err := command.Cd(dir).Do()
			if nil != err {
				installInfo <- NewInfo(dir, false, goutils.ToString(b))
			} else {
				installInfo <- NewInfo(dir, true, "")
			}
			excuted = true
			command.Cd("..")
		}
	}
}

func logging() {
	var info *Info
	now := 0
	after := 0
	defer func() {
		fmt.Printf("install: %d.", now)
	}()
	ticker := time.NewTicker(12e8)
	for {
		select {
		case info = <-installInfo:
			fmt.Println(info.String())
			now++
		case <-ticker.C:
			after++
			if now < after {
				return
			}
			after = now
		}
	}
}

func main() {
	wd, err := os.Getwd()
	if exc.Checkerr(err) {
		os.Exit(-1)
	}
	go searchDir(wd)
	time.Sleep(2e8)
	logging()
}