package config

import (
	"encoding/json"
	"log"
	"os"
)

type FunctionList struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	SubItems []SubItems `json:"subItems"`
}

type SubItems struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type FunctionConfig struct {
	List []FunctionList `json:"functionList"`
}

// 解析function.json文件
var functionList = &FunctionConfig{
	List: []FunctionList{},
}

func InitFunction(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("初始化function.json文件错误")
	}

	var fl []FunctionList
	err = json.Unmarshal(file, &fl)
	if err != nil {
		log.Fatal("读取function.json，解析成对应的结构体值错误")
	}

	functionList.List = fl
}

func GetFunctionList() *FunctionConfig {
	return functionList
}
