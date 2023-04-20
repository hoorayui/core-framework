package flag

import (
	"flag"
	"fmt"
	"framework/components/config"
	"framework/util"
	"net/http"
	"os"
	"runtime"
)

var (
	APPNAME        = "GoMetal"
	ShowVersion    bool
	Ping           bool
	ConfigFile     string
	LogDir         string
	BackendVersion string
	//	编译注入
	Version   string
	BuildTime string
	OsArch    string
)

func init() {
	flag.BoolVar(&ShowVersion, "v", false, "show version info")
	flag.BoolVar(&Ping, "ping", false, "check server health")
	flag.StringVar(&ConfigFile, "f", "", "set config file")
	flag.StringVar(&LogDir, "log-dir", util.GetAppRoot()+"/..", "set log file directory")
}

// ParseOrDie : parse flags
func ParseOrDie() {
	if !flag.Parsed() {
		flag.Parse()
	}
	if ShowVersion {
		fmt.Printf(APPNAME+` %s, Compiler: %s, %s, Copyright (C) 2022 Pintechs Inc.`,
			"v0.1",
			runtime.Compiler,
			runtime.Version())
		fmt.Printf("\nVersion: %s\nBuilt: %s\nOS/Arch: %s\n", Version, BuildTime, OsArch)
		os.Exit(0)
	}
	if Ping {
		_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", config.GetInstance().Server.ListenPort))
		if nil != err {
			println("未在运行")
			os.Exit(1)
		}
		println("正常运行中")
		os.Exit(0)
	}
}
