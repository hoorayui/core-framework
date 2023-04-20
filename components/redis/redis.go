package redis

import (
	"encoding/json"
	"errors"
	goredis "github.com/go-redis/redis"
	"github.com/hoorayui/core-framework/types"
	"github.com/sirupsen/logrus"
)

type Instance struct {
	config types.RedisConfig
	client *goredis.Client
}

var instance *Instance

// GetName 组件名称
func (i *Instance) GetName() string {
	return "redis"
}

// Init 初始化实例
func (i *Instance) Init(config interface{}) error {
	instance = i
	bytes, _ := json.Marshal(config)
	json.Unmarshal(bytes, &i.config)
	i.Validate()
	i.client = goredis.NewClient(&goredis.Options{
		Addr:     i.config.Addr,
		Password: i.config.Password, // no password set
		DB:       i.config.DB,       // use default DB
	})
	if err := i.client.Ping().Err(); err != nil {
		return err
	}
	return nil
}

// Validate 验证配置
func (i *Instance) Validate() error {
	if i.config.Addr == "" {
		return errors.New("redis地址为空")
	}
	if i.config.Password == "" {
		logrus.Warn("redis密码为空，为了安全，请设置密码")
	}
	return nil
}

// GetInstance 获取实例
func GetInstance() *Instance {
	return instance
}
func (i *Instance) Client() *goredis.Client {
	return i.client
}

// Close 关闭
func (i *Instance) Close() {
	i.client.Close()
}
