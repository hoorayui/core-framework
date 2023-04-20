package event

import (
	"github.com/hoorayui/core-framework/components/log"
)

// registerHook 注册回调函数
func registerHook(e string, hook func(string, interface{}) error) {
	GetInstance().l.Lock()
	defer GetInstance().l.Unlock()
	if watcher, ok := GetInstance().watcher[e]; ok {
		watcher = append(watcher, hook)
		GetInstance().watcher[e] = watcher
	} else {
		watcher = []func(string, interface{}) error{hook}
		GetInstance().watcher[e] = watcher
	}
}
func trigger(e string) {
	// get listeners and run them
	watchers, ok := GetInstance().watcher[e]
	if !ok {
		log.Errorf("不支持的事件[%s]", e)
		return
	}
	for _, listener := range watchers {
		err := listener(e, nil)
		if err != nil {
			log.Errorf("事件[%s]回调执行失败:%s", e, err.Error())
		}
	}
}
