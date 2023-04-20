package test

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func DisplayObject(obj interface{}) {
	js, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", js)
}

// IsInTests 判断是否在单元测试中
func IsInTests() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}
