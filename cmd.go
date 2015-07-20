package main

import (
	"log"
	"os"
	"os/exec"
	// "path/filepath"
	"strings"
	"time"
)

var (
	command    = "go install"
	ResultChan chan Result
)

type Result struct {
	dir string
	ok  bool
}

func init() {
	ResultChan = make(chan Result, 20)
}

func checkErr(err error) bool {
	if nil != err {
		log.Println(err)
		return true
	}
	return false
}

func targetPath(basePath string) []string {
	targetPaths := []string{basePath}
	baseDir, err := os.Open(basePath)
	if checkErr(err) {
		return targetPaths
	}

	// log.Print(baseDir.Name())

	subDirs, err := baseDir.Readdir(-1)
	if checkErr(err) {
		return targetPaths
	}

	err = os.Chdir(basePath)
	if checkErr(err) {
		return targetPaths
	}
	base, err := os.Getwd()
	if !checkErr(err) {
		targetPaths = append(targetPaths, base)
	}

	executing := false
	for _, v := range subDirs {
		fileInfo := v.(os.FileInfo)
		if fileInfo.IsDir() {
			continue
		}
		if strings.Contains(fileInfo.Name(), ".go") {
			executing = true
		}
	}
	if executing {
		// log.Printf("current Dir:%s", base)
		executed := executeCmdHere(command)
		baseP, _ := os.Getwd()
		result := Result{dir: baseP, ok: executed}
		ResultChan <- result
	}

	for _, v := range subDirs {
		// log.Printf("subDirs:%s", v.Name())
		fileInfo := v.(os.FileInfo)
		// log.Printf("subDirs:%s", fileInfo.Name())
		if fileInfo.IsDir() {
			if strings.EqualFold(fileInfo.Name(), ".git") {
				continue
			}
			// targetPaths = append(targetPaths, fileInfo.Name())
			subTargetPaths := targetPath(fileInfo.Name())
			targetPaths = append(targetPaths, subTargetPaths[1:]...)
		}
	}
	return targetPaths
}

func executeCmdHere(command string) bool {
	cmdWithArgs := strings.Split(command, " ")
	var cmd *exec.Cmd
	cmdLength := len(cmdWithArgs)
	realCmd := cmdWithArgs[0]
	var args []string
	if cmdLength > 1 {
		args = cmdWithArgs[1:cmdLength]
		cmd = exec.Command(realCmd, args...)
	} else {
		cmd = exec.Command(realCmd)
	}
	_, err := cmd.Output()
	if err != nil {
		log.Printf("CmdRunError(cmd=%s, agrs=%v): %s", realCmd, args, err)
		return false
	}
	// log.Printf("Output(cmd=%s, agrs=%v): %v", realCmd, args, string(result))
	return true
}

func logging() {
	for {
		select {
		case result := <-ResultChan:
			log.Printf("[LOG] %v\n", result)
		}
	}
}
func main() {
	basePath, err := os.Getwd()
	if checkErr(err) {
		return
	}
	log.Printf("Base Path: %s", basePath)
	go logging()
	targetPath(basePath)

	time.Sleep(1e8)
	close(ResultChan)
}
