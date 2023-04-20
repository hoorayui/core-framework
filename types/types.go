package types

import "fmt"

type Config struct {
	Config CfgConfig `yaml:"config" json:"config" toml:"config"`
	// flags
	Log LogConfig `yaml:"log" json:"log" toml:"log"`

	// server flags
	Server ServerConfig `yaml:"server" json:"server" toml:"server"`

	// db flags
	DB    DBConfig    `yaml:"mysql" json:"mysql" toml:"mysql"`
	Redis RedisConfig `yaml:"redis" json:"redis" toml:"redis"`
}
type CfgConfig struct {
	Path string `yaml:"path" json:"path" toml:"path"`
}
type DBConfig struct {
	InitDatabase              bool       `yaml:"init_database" json:"init_database" toml:"init_database"`
	DBDSN                     []MySQLDSN `yaml:"db_dsn" json:"db_dsn" toml:"db_dsn"`
	MysqlDriverParams         string     `yaml:"mysql_driver_params" json:"mysql_driver_params" toml:"mysql_driver_params"`
	DBDriver                  string     `yaml:"db_driver" json:"db_driver" toml:"db_driver"`
	DBConnectTimeoutInSeconds int        `yaml:"db_connect_timeout_in_seconds" json:"db_connect_timeout_in_seconds" toml:"db_connect_timeout_in_seconds"`
	DBConnectionMaxLifetime   int        `yaml:"db_connection_max_lifetime" json:"db_connection_max_lifetime" toml:"db_connection_max_lifetime"`
	DBMaxOpenConn             int        `yaml:"db_max_open_conn" json:"db_max_open_conn" toml:"db_max_open_conn"`
	DBMaxIdleConn             int        `yaml:"db_max_idle_conn" json:"db_max_idle_conn" toml:"db_max_idle_conn"`
	DebugMode                 bool       `yaml:"debug_mode" json:"debug_mode" toml:"debug_mode"`
}
type MySQLDSN struct {
	DBHost     string `yaml:"db_host" json:"db_host" toml:"db_host"`
	DBPort     int    `yaml:"db_port" json:"db_port" toml:"db_port"`
	DBUser     string `yaml:"db_user" json:"db_user" toml:"db_user"`
	DBPassword string `yaml:"db_password" json:"db_password" toml:"db_password"`
	DBDatabase string `yaml:"db_database" json:"db_database" toml:"db_database"`
}

func (m MySQLDSN) String(driver string) string {
	if driver == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			m.DBUser, m.DBPassword, m.DBHost, m.DBPort, m.DBDatabase)
	} else if driver == "postgres" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
			m.DBHost, m.DBPort, m.DBUser, m.DBDatabase, m.DBPassword)
	} else {
		return ""
	}
}

type LogConfig struct {
	LogLevel         string `yaml:"log_level" json:"log_level" toml:"log_level"`
	LogDir           string `yaml:"log_dir" json:"log_dir" toml:"log_dir"`
	LogFileName      string `yaml:"log_file_name" json:"log_file_name" toml:"log_file_name"`
	ErrorLogFileName string `yaml:"log_error_file_name" json:"log_error_file_name" toml:"log_error_file_name"`
	LogFormat        string `yaml:"log_format" json:"log_format" toml:"log_format"`
	LogReserveDays   int    `yaml:"log_reserve_days" json:"log_reserve_days" toml:"log_reserve_days"`
	LogMaxSize       int    `yaml:"log_max_size" json:"log_max_size" toml:"log_max_size"`
}

type ServerConfig struct {
	ListenPort        int    `yaml:"listen_port" json:"listen_port" toml:"listen_port"`
	GatewayListenPort int    `yaml:"gateway_listen_port" json:"gateway_listen_port" toml:"gateway_listen_port"`
	GrpcEndpoint      string `yaml:"grpc_endpoint" json:"grpc_endpoint" toml:"grpc_endpoint"`
	RunMode           string `yaml:"run_mode" json:"run_mode" toml:"run_mode"`
	DebugMode         bool   `yaml:"debug_mode" json:"debug_mode" toml:"debug_mode"`
	BaseUrl           string `yaml:"base_url" json:"base_url" toml:"base_url"`
}

type RedisConfig struct {
	Addr     string `yaml:"redis_addr" json:"redis_addr" toml:"redis_addr"`
	Password string `yaml:"redis_password" json:"redis_password" toml:"redis_password"`
	DB       int    `yaml:"redis_db" json:"redis_db" toml:"redis_db"`
}
