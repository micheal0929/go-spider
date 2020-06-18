package version

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	BuildTime    string //构建时间
	BuildVersion string //构建的Git Commit Id
	BuildBranch  string //构建所在的分支
	BuildName    string //构建者的名称
	BuildPath    string //构建路径
	GoVersion    string //构建所用的Go版本
)
var (
	ProgramName string
)
var aviliableVerb = map[string]bool{
	"-v":        true,
	"--version": true,
	"version":   true,
}

func init() {
	if len(os.Args) > 1 {
		if aviliableVerb[os.Args[1]] {
			fmt.Printf("Build Time    : %s \n", BuildTime)
			fmt.Printf("Build Version : %s \n", BuildVersion)
			fmt.Printf("Build Branch  : %s \n", BuildBranch)
			fmt.Printf("Build Name    : %s \n", BuildName)
			fmt.Printf("Build Path    : %s \n", BuildPath)
			fmt.Printf("Go Version    : %s \n", GoVersion)
			os.Exit(0)
		} else {
			ProgramName = filepath.Base(os.Args[0])
		}
	}
}
