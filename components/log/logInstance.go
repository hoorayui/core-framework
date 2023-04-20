package log

import (
	"encoding/json"
	"time"

	"github.com/hoorayui/core-framework/types"
	"github.com/sirupsen/logrus"
)

type Instance struct {
	config types.LogConfig
	logger *logrus.Logger
}

var instance *Instance

func init() {
	instance = &Instance{
		logger: logrus.New(),
	}
}

// GetName 组件名称
func (i *Instance) GetName() string {
	return "log"
}

// Init 初始化实例
func (i *Instance) Init(config interface{}) error {
	instance = i
	bytes, _ := json.Marshal(config)
	json.Unmarshal(bytes, &i.config)
	i.Validate()

	i.logger = logrus.New()
	i.logger.SetReportCaller(true)
	i.logger.SetNoLock()
	i.logger.SetOutput(GetMultiWriter())
	i.logger.Formatter = &UTCFormatter{&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
		PrettyPrint:     false,
	}}

	i.logger.SetLevel(logrus.InfoLevel)

	return nil
}

// Validate 验证配置
func (i *Instance) Validate() error {
	return nil
}

// GetInstance 获取实例
func GetInstance() *Instance {
	return instance
}

// Close 关闭
func (i *Instance) Close() {
}
