package main

import (
	// "flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// command = flag.String("go", "", "install")
	command = &("")
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	// flag.Parse()
	if len(*command) == 0 {
		log.Println("The argument '--command' is NOT specified!")
		return
	}
	basePath, err := os.Getwd()
	if err != nil {
		log.Println("GetwdError:", err)
		return
	}
	log.Printf("Base Path: %s\\n", basePath)
	baseDir, err := os.Open(basePath)
	if err != nil {
		log.Printf("OpenError (%s): %s\\n", baseDir, err)
		return
	}
	subDirs, err := baseDir.Readdir(-1)
	if err != nil {
		log.Println("ReaddirError:", err)
		return
	}
	targetPaths := []string{}
	for _, v := range subDirs {
		fileInfo := v.(os.FileInfo)
		if fileInfo.IsDir() {
			log.Printf("Target: %s\\n", fileInfo.Name())
			targetPaths = append(targetPaths, fileInfo.Name())
		} else {
			log.Printf("Non-target: %s\\n", fileInfo.Name())
		}
	}
	for _, v := range targetPaths {
		err = os.Chdir(basePath)
		if err != nil {
			log.Printf("ChdirError (%s): %s\\n", baseDir, err)
			return
		}
		targetPath, err := filepath.Abs(v)
		if err != nil {
			log.Printf("AbsError (%s): %s\\n", v, err)
			return
		}
		log.Printf("Target Path: %s\\n", targetPath)
		err = os.Chdir(targetPath)
		if err != nil {
			log.Printf("ChdirError (%s): %s\\n", targetPath, err)
			return
		}
		cmdWithArgs := strings.Split(*command, " ")
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
			log.Printf("CmdRunError (cmd=%s, agrs=%v): %s\\n", realCmd, args, err)
			return
		}
		log.Printf("Output (dir=%s, cmd=%s, agrs=%v): \\n%v\\n", targetPath, realCmd, args, string(result))
	}
	log.Println("The command(s) execution has been finished.")
}

//该片段来自于http://outofmemory.cn
