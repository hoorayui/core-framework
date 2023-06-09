package core

import (
	"fmt"
	"github.com/hoorayui/core-framework/components/config"
	"github.com/hoorayui/core-framework/components/event"
	"github.com/hoorayui/core-framework/components/log"
	"github.com/hoorayui/core-framework/components/mysql"
	"github.com/hoorayui/core-framework/components/redis"
	"github.com/hoorayui/core-framework/core/middleware"
	"github.com/hoorayui/core-framework/types"
	"github.com/hoorayui/core-framework/util"
	"github.com/hoorayui/core-framework/util/flag"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var version = "0.0.1"

type InterfaceCore interface {
	New(string) InterfaceCore              // 初始化组件
	SetConf(string)                        // 初始化组件
	Run()                                  // 运行web server
	InitComponents(...InterfaceComponents) // 初始化组件
	Stop()                                 // 关闭服务，等待业务处理完成
}

type core struct {
	configFile string
	WorkDir    string //TODO
	Debug      bool   // TODO
	components map[string]InterfaceComponents
	deferFuncs map[string]func()
}

// New 初始化应用组件并返回一个应用实例
func New(configFile string, components ...InterfaceComponents) *core {
	app := &core{
		components: map[string]InterfaceComponents{},
		deferFuncs: map[string]func(){},
	}
	flag.BackendVersion = version
	flag.ParseOrDie()

	app.WorkDir = util.GetAppRoot()
	app.configFile = flag.ConfigFile
	if app.configFile == "" {
		app.configFile = configFile
	}
	app.LoadComponents(&config.Instance{}, types.CfgConfig{
		Path: app.configFile,
	})
	app.LoadComponents(&log.Instance{}, config.GetConfig("log"))
	app.InitComponents(components...)
	return app
}
func Default(configFile string) *core {
	app := New(configFile, &mysql.Instance{},
		&event.Instance{},
		&mysql.Instance{},
		&redis.Instance{})
	return app
}

// Run start server
func (c *core) Run() {
	//rest 服务
	server := gin.New()
	pprof.Register(server)
	gin.DefaultWriter = log.GetMultiWriter()
	server.Use(middleware.GetAllBaseMiddleware()...)

	// 对外统一监听端口
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", config.GetInstance().Server.ListenPort))
	logrus.Infof("web server is starting,listening on port [%d]", config.GetInstance().Server.ListenPort)
	err := http.Serve(ln, server)
	if nil != err {
		log.Errorf("Server failed to run, err: %v", err)
		// TODO 传递退出错误
		return
	}
	// TODO 传递退出错误
	return
}

// TODO addDefer
// Stop the function that stop program with return value
func (c *core) Stop() {
	for _, close := range c.deferFuncs {
		close()
	}
	// it is required, to work around bug of go
	// 在程序退出前，需要先sleep一段时间，否则有可能日志打印不全
	time.Sleep(100 * time.Millisecond)
	//if hasError {
	//	logrus.Error("程序异常退出")
	//	os.Exit(-1)
	//}
	logrus.Info("程序运行结束")
	os.Exit(0)
}
