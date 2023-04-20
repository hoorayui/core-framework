package event

import (
	"encoding/json"
	"framework/types"
	"sync"
)

type Instance struct {
	config  types.RedisConfig
	watcher map[string][]func(string, interface{}) error
    l sync.RWMutex
}

var instance *Instance

// GetName 组件名称
func (i *Instance) GetName() string {
	return "event"
}

// Init 初始化实例
func (i *Instance) Init(config interface{}) error {
	instance = i
    i.l = sync.RWMutex{}
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
