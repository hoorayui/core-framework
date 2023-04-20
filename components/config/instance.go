package config

import (
	"encoding/json"

	"github.com/hoorayui/core-framework/types"
	"github.com/hoorayui/core-framework/util"
)

type Instance struct {
	config     types.CfgConfig
	cfg        types.Config
	configFile string
}

var instance *Instance

func (i *Instance) GetName() string {
	return "config"
}
func (i *Instance) Init(config interface{}) error {
	instance = i
	bytes, _ := json.Marshal(config)
	json.Unmarshal(bytes, &instance.config)
	i.Validate()
	util.MustLoadConfig(instance.config.Path, &instance.cfg)
	return nil
}

// Validate 验证配置
func (i *Instance) Validate() error {
	return nil
}

// GetInstance 获取实例
func GetInstance() types.Config {
	return instance.cfg
}

// Close 关闭
func (i *Instance) Close() {
}
func GetConfigMap() map[string]interface{} {
	var configMap map[string]interface{}
	bytes, _ := json.Marshal(instance.cfg)
	json.Unmarshal(bytes, &configMap)
	return configMap
}
func GetConfig(key string) interface{} {
	return GetConfigMap()[key]
}
