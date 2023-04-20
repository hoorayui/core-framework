package config

import (
	"os"

	"github.com/hoorayui/core-framework/types"
	"github.com/hoorayui/core-framework/util"

	"gopkg.in/yaml.v2"

	"log"

	"github.com/spf13/viper"
)

var cfg = &Config{
	Options: &Options{},
}

// Config 自定义配置
type Config struct {
	*types.Config
	Options *Options
}

// Options 读取yam.yml配置文件
type Options struct {
	TokenExpireDuration      int
	RetrieveLogRetentionTime int
	AccessKey                string
	CasbinFileName           string
	UploadDir                string
	DocumentLink             string
	InitAccount              []string
}

// NewOption 读取iam.yml文件，生成options需要的结果
func NewOption(path string) (*Options, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Options{
		TokenExpireDuration:      viper.GetInt("token_expire_duration"),
		RetrieveLogRetentionTime: viper.GetInt("retrieve_log_retention_time"),
		AccessKey:                viper.GetString("access_key"),
		CasbinFileName:           viper.GetString("casbin_file_name"),
		UploadDir:                viper.GetString("upload_dir"),
		DocumentLink:             viper.GetString("document_link"),
		InitAccount:              viper.GetStringSlice("init_account"),
	}, nil
}

func InitConfig(path string) {
	fs, err := os.Stat(path)
	if err != nil || fs.IsDir() {
		log.Fatalf("配置文件路径[%s]不正确", path)
	}
	util.MustLoadConfig(path, cfg)
	options, err := NewOption(path)
	if err != nil {
		log.Fatalf("解析配置yaml文件失败，错误:[%s]", err.Error())
	}
	cfg.Options = options
}

// GenerateSampleYaml : 生成配置示例文件
func GenerateSampleYaml() {
	cfg := types.Config{}
	cfg.DB.DBDSN = []types.MySQLDSN{
		{DBHost: "localhost"},
	}
	yb, _ := yaml.Marshal(cfg)
	println(string(yb))
}
