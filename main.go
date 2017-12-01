package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/everfore/exc"
	"github.com/toukii/goutils"
	"github.com/toukii/jsnm"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
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
	if length > 850 {
		length = 850
		return fmt.Sprintf("[FAIL]: %s \n%s [more]...", i.dir, i.info[:length])
	}
	return fmt.Sprintf("[FAIL]: %s \n%s", i.dir, i.info)
}

func NewInfo(dir string, ok bool, info string) *Info {
	return &Info{
		dir:  dir,
		ok:   ok,
		info: info,
	}
}

func errPkgs(golist []byte) []string {
	js := jsnm.BytesFmt(golist)
	if arrs := js.Get("DepsErrors").Arr(); len(arrs) > 0 {
		pkgs := make([]string, 0, 5)
		pkgm := make(map[string]bool)
		pkgm[js.Get("ImportPath").RawData().String()] = true
		for _, arr := range arrs {
			depkgs := arr.Get("ImportStack").Arr()
			for _, pkg := range depkgs {
				pkgName := pkg.RawData().String()
				if _, ex := pkgm[pkgName]; ex {
					continue
				}
				pkgm[pkgName] = true
				pkgs = append(pkgs, pkgName)
			}
		}
		return pkgs
	}
	return nil
}

var (
	command     *exc.CMD
	installInfo chan *Info
	dir         = kingpin.Arg("./", "go build -R path").String()
)

func init() {
	installInfo = make(chan *Info, 50)
	command = exc.NewCMD("go install")
}

func pull(pkgs []string) {
	size := len(pkgs)
	if size <= 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(size)
	for _, pkg := range pkgs {
		if !strings.HasPrefix(pkg, "github.com") {
			fmt.Printf("pkg %s is not supported.\n", pkg)
			wg.Done()
			continue
		}
		go func(pkg string) {
			exc.NewCMD("pull -r " + pkg[11:]).Debug().Execute()
			wg.Done()
		}(pkg)
	}
	wg.Wait()
}

func searchDir(dir string) {
	bs, _ := exc.NewCMD("go list -json").Do()
	pull(errPkgs(bs))
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
			/*go*/ searchDir(filepath.Join(dir, it.Name()))
		}
		if strings.HasSuffix(it.Name(), ".go") && !excuted {
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
		fmt.Printf("install: %d.\n", now)
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
	_ = kingpin.Parse()
	wd, err := os.Getwd()
	if exc.Checkerr(err) {
		os.Exit(-1)
	}
	searchDir(filepath.Join(wd, *dir))
	logging()
}
