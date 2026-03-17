package conf

import (
	"os"
	"strings"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Conf 聚合了所有的配置
type Conf struct {
	Viper  *viper.Viper
	Config *Config
}

// ==============================================================================
// 1. 企业级分层配置结构体定义 (使用 mapstructure 标签映射 YAML 节点)
// ==============================================================================

type Config struct {
	Env        string     `mapstructure:"env"` // dev, test, prod
	Server     Server     `mapstructure:"server"`
	Log        Log        `mapstructure:"log"`
	Database   Database   `mapstructure:"database"`
	Redis      Redis      `mapstructure:"redis"`
	Kafka      Kafka      `mapstructure:"kafka"`
	RabbitMQ   RabbitMQ   `mapstructure:"rabbitmq"`
	ClickHouse ClickHouse `mapstructure:"clickhouse"`
	Elastic    Elastic    `mapstructure:"elastic"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Server struct {
	Addr string `mapstructure:"addr"`
}

type Log struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
}

type Database struct {
	Path string `mapstructure:"path"` // 给 SQLite 用的
	DSN  string `mapstructure:"dsn"`  // 给 MySQL 用的
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"` // 支持多个 Broker
	GroupID string   `mapstructure:"group_id"`
}

type RabbitMQ struct {
	URL string `mapstructure:"url"` // amqp://guest:guest@localhost:5672/
}

type ClickHouse struct {
	Addr     string `mapstructure:"addr"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Elastic struct {
	Addresses []string `mapstructure:"addresses"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
}

type Prometheus struct {
	MetricsPath string `mapstructure:"metrics_path"` // 暴露的端点，默认 /metrics
	Port        int    `mapstructure:"port"`
}

// ==============================================================================
// 2. Viper 实例化与环境解析
// ==============================================================================

func ProvideViper() (*viper.Viper, error) {
	v := viper.New()

	// 1. 设置配置文件路径
	if cfg := os.Getenv("CONFIG"); cfg != "" {
		v.SetConfigFile(cfg)
	} else {
		v.SetConfigFile("configs/config.yaml")
	}
	v.SetConfigType("yaml")

	// 2. 极其核心：开启环境变量自动覆盖！
	// 比如你在服务器配置了环境变量 AIGC_REDIS_PASSWORD=123
	// Viper 会自动把这个值覆盖到 Redis.Password 字段上
	v.SetEnvPrefix("AIGC") // 你的项目环境变量前缀
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 忽略文件不存在的错误（方便本地纯用环境变量或默认值跑）
	_ = v.ReadInConfig()

	return v, nil
}

// ==============================================================================
// 3. 一把梭反序列化 (彻底告别手工 GetString)
// ==============================================================================

func ProvideConfig(v *viper.Viper) (*Config, error) {
	c := &Config{
		// 可以在这里设置一些硬核的默认值
		Env: "dev",
		Server: Server{
			Addr: ":8080",
		},
		Database: Database{
			Path: "./demo_logs.db",
		},
		Log: Log{
			Level: "info",
		},
	}

	if v != nil {
		// 🔥 核心黑魔法：直接把 YAML 的结构全部灌进 Go 的结构体里
		if err := v.Unmarshal(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// ProviderSet 暴露给 Wire
var ProviderSet = wire.NewSet(
	ProvideViper,
	ProvideConfig,
	wire.Struct(new(Conf), "*"),
)
