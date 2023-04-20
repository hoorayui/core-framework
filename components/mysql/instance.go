package mysql

import (
	"encoding/json"
	"framework/types"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Instance struct {
	config types.DBConfig
	client *gorm.DB
}

var instance *Instance

// GetName 组件名称
func (i *Instance) GetName() string {
	return "mysql"
}

// Init 初始化实例
func (i *Instance) Init(config interface{}) error {
	instance = i
	bytes, _ := json.Marshal(config)
	json.Unmarshal(bytes, &i.config)
	i.Validate()
	var dialector gorm.Dialector
	if i.config.DBDriver == "mysql" {
		dialector = mysql.Open(i.config.DBDSN[0].String(i.config.DBDriver))
	} else if i.config.DBDriver == "postgres" {
		dialector = postgres.Open(i.config.DBDSN[0].String(i.config.DBDriver))
	} else {
		logrus.Fatalf("不支持[%s]数据库", i.config.DBDriver)
	}
	rawDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 mysqlLogger(),
	})
	rawDB.Exec("set time_zone=\"+08:00\";")
	if err != nil {
		// 数据库未创建
		if strings.Contains(err.Error(), "Unknown database") {
			logrus.Fatalf("数据库[%s] 未创建", i.config.DBDSN[0].DBDatabase)
		}
		// 数据库连接失败
		logrus.Fatalf("Connect database failed, err: %v", err)
		return err
	}
	sqlDB, err := rawDB.DB()
	if nil != err {
		logrus.Fatalf("获取数据库实例失败，%s", err.Error())
	}
	sqlDB.SetMaxIdleConns(i.config.DBMaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(i.config.DBConnectTimeoutInSeconds))
	sqlDB.SetMaxOpenConns(i.config.DBMaxOpenConn)
	return nil
}

// Validate 验证配置
func (i *Instance) Validate() error {
	// TODO
	return nil
}

// GetInstance 获取实例
func GetInstance() *Instance {
	return instance
}
func (i *Instance) Client() *gorm.DB {
	return i.client
}

// Close 关闭
func (i *Instance) Close() {
	return
}
