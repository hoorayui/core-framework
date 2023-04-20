package util

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var loaders = map[string]func([]byte, interface{}) error{
	".json": LoadFromJsonBytes,
	".toml": LoadFromTomlBytes,
	".yaml": LoadFromYamlBytes,
	".yml":  LoadFromYamlBytes,
}

func LoadConfig(file string, v interface{}) error {
	content, err := os.ReadFile(file)
	if nil != err {
		return err
	}
	loader := loaders[strings.ToLower(path.Ext(file))]
	return loader(content, v)
}

func LoadFromJsonBytes(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

func LoadFromYamlBytes(b []byte, v interface{}) error {
	return yaml.Unmarshal(b, v)
}

func LoadFromTomlBytes(b []byte, v interface{}) error {
	return toml.Unmarshal(b, v)
}

func GenerageSampleYaml(v interface{}) {
	yb, _ := yaml.Marshal(v)
	println(string(yb))
}

func MustLoadConfig(file string, v interface{}) {
	if err := LoadConfig(file, v); err != nil {
		logrus.Errorf("加载配置文件[%s]失败", file)
		return
	}
	logrus.Infof("配置文件[%s]加载成功", file)
}
