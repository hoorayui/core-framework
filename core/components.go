package core

import (
	"log"

	"github.com/hoorayui/core-framework/components/config"
	"github.com/sirupsen/logrus"
)

// InterfaceComponents 组件接口
type InterfaceComponents interface {
	GetName() string // 获取组件名
	//GetInstance() *InterfaceComponents // 获取实例
	Init(interface{}) error // 初始化
	Close()                 // 关闭组件
}

func (c *core) InitComponents(components ...InterfaceComponents) {
	c.components = map[string]InterfaceComponents{}
	c.deferFuncs = map[string]func(){}
	for _, v := range components {
		// 初始化
		if err := v.Init(config.GetConfig(v.GetName())); err != nil {
			log.Fatalf("组件[%s]初始化失败:错误详情：%s", v.GetName(), err.Error())
		}
		// 注册组件
		c.components[v.GetName()] = v
		// 关闭组件回调
		c.deferFuncs[v.GetName()] = v.Close
		logrus.Infof("组件[%s]加载成功", v.GetName())
	}
	logrus.Info("组件加载完成")
}
func (c *core) LoadComponents(component InterfaceComponents, config interface{}) {
	component.Init(config)
	c.components[component.GetName()] = component
	c.deferFuncs[component.GetName()] = component.Close
}
