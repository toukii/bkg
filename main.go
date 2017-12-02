package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/everfore/exc"
	cr "github.com/fatih/color"
	"github.com/toukii/goutils"
	"github.com/toukii/jsnm"
	// "github.com/toukii/jsnm2"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type Info struct {
	dir  string
	ok   bool
	info string
}

func TrimGopath(pth string) string {
	return strings.TrimPrefix(pth, os.Getenv("GOPATH")+"/src/")
}

func (i Info) String() string {
	if i.ok {
		return cr.GreenString("[SUCCESS] ") + cr.CyanString(TrimGopath(i.dir))
	}
	failMsg := cr.RedString("[FAIL] ") + cr.YellowString(TrimGopath(i.dir))
	length := len(i.info)
	if length > 850 {
		length = 850
		return fmt.Sprintf("%s --> %s %s", failMsg, i.info[:length], cr.CyanString("[more]..."))
	}
	return fmt.Sprintf("%s --> %s", failMsg, i.info)
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
		if len(pkgs) > 0 {
			fmt.Println(cr.HiCyanString("missed import packages:\n"), cr.RedString(displayPkgs(pkgs)))
		}
		return pkgs
	}
	return nil
}

func displayPkgs(pkgs []string) string {
	newPkgs := make([]string, len(pkgs))
	for i, pkg := range pkgs {
		newPkgs[i] = fmt.Sprintf("\t%d. %s", i+1, pkg)
	}
	return strings.Join(newPkgs, "\n")
}

func imports(golist []byte) []string {
	arrs := jsnm.BytesFmt(golist).Get("Imports").Arr()
	pkgs := make([]string, 0, 5)
	for _, arr := range arrs {
		pkg := arr.RawData().String()
		if strings.HasPrefix(pkg, "github.com") || strings.HasPrefix(pkg, "gopkg.in") || strings.HasPrefix(pkg, "golang.org") {
			pkgs = append(pkgs, pkg)
		}
	}
	if len(pkgs) > 0 {
		fmt.Println(cr.HiCyanString("import packages:\n"), cr.HiGreenString(displayPkgs(pkgs)))
	}
	return pkgs
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
		if strings.HasPrefix(pkg, "gopkg.in") || strings.HasPrefix(pkg, "golang.org") {
			go func() {
				exc.NewCMD("go get " + pkg).Debug().DoTimeout(20e9)
				wg.Done()
			}()
			continue
		}
		if !strings.HasPrefix(pkg, "github.com") {
			cr.Red("pkg %s is not supported.\n", pkg)
			wg.Done()
			continue
		}
		go func(pkg string) {
			if strings.HasPrefix(pkg, "github.com/toukii") || strings.HasPrefix(pkg, "github.com/everfore") || strings.HasPrefix(pkg, "github.com/datc/") {
				exc.NewCMD("pull " + pkg).Debug().DoTimeout(20e9)
			} else {
				exc.NewCMD("pull -r " + pkg).Debug().DoTimeout(20e9)
			}
			wg.Done()
		}(pkg)
	}
	wg.Wait()
}

func searchDir(dir string) {
	file, err := os.Open(dir)
	if goutils.CheckErr(err) {
		return
	}

	bs, err := exc.NewCMD("go list -json " + TrimGopath(dir)).Do()

	if err == nil && len(bs) > 0 {
		imports(bs)
		pull(errPkgs(bs))
	}
	subdirs, err := file.Readdir(-1)
	if goutils.CheckErr(err) {
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
	success := 0
	after := 0
	defer func() {
		fmt.Printf("scan %s go pkgs, go install success %s.\n", cr.RedString("%d", now), cr.HiGreenString("%d", success))
	}()
	ticker := time.NewTicker(1e8)
	for {
		select {
		case info = <-installInfo:
			fmt.Println(info.String())
			if info.ok {
				success++
			}
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
	if goutils.CheckErr(err) {
		return
	}
	searchDir(filepath.Join(wd, *dir))
	logging()
}
