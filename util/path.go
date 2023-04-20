package util

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// 定义项目根目录
func GetAppRoot() string {
	dir := getBinAbPath()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		dir = getCallerAbPath()
	}
	return path.Dir(dir)
}

func getBinAbPath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func getCallerAbPath() string {
	var abPath string
	//TODO
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
