package main

import (
	"log"
	"os"
	"os/exec"
	// "path/filepath"
	"strings"
)

var (
	command = "go install"
)

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
	subDirs, err := baseDir.Readdir(-1)
	if checkErr(err) {
		return targetPaths
	}
	for _, v := range subDirs {
		fileInfo := v.(os.FileInfo)
		if fileInfo.IsDir() {
			if strings.EqualFold(fileInfo.Name(), ".git") {
				continue
			}
			targetPaths = append(targetPaths, fileInfo.Name())
		}
	}
	return targetPaths
}

func executeCmd(basePath, targetPath, command string) {
	err := os.Chdir(basePath)
	if checkErr(err) {
		return
	}
	// absTargetPath, err := filepath.Abs(targetPath)
	// if checkErr(err) {
	// 	log.Printf("AbsError (%s): %s", targetPath, err)
	// 	return
	// }
	err = os.Chdir(targetPath)
	if checkErr(err) {
		return
	}

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
	result, err := cmd.Output()
	if err != nil {
		log.Printf("CmdRunError(dir=%s, cmd=%s, agrs=%v): %s", targetPath, realCmd, args, err)
		return
	}
	log.Printf("Output(dir=%s, cmd=%s, agrs=%v): %v", targetPath, realCmd, args, string(result))
}

func main() {
	basePath, err := os.Getwd()
	if checkErr(err) {
		return
	}
	log.Printf("Base Path: %s", basePath)

	targetPaths := targetPath(basePath)

	for _, v := range targetPaths {
		executeCmd(basePath, v, command)
	}
	log.Println("The command(s) execution has been finished.")
}
