package components

import (
	"encoding/json"

	goredis "github.com/go-redis/redis"
	"github.com/hoorayui/core-framework/types"
)

type Instance struct {
	config types.RedisConfig
	client *goredis.Client
}

var instance *Instance

// GetName 组件名称
func (i *Instance) GetName() string {
	return "UNTITLED"
}

// Init 初始化实例
func (i *Instance) Init(config interface{}) error {
	instance = i
	bytes, _ := json.Marshal(config)
	json.Unmarshal(bytes, &i.config)
	i.Validate()
	//TODO
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
