package conf

import (
	"os"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Config 简化的配置结构
type Config struct {
	DatabasePath string `mapstructure:"database.path"`
	ServerAddr   string `mapstructure:"server.addr"`
}

// ProvideViper 返回已加载的 viper 实例
func ProvideViper() (*viper.Viper, error) {
	v := viper.New()
	if cfg := os.Getenv("CONFIG"); cfg != "" {
		v.SetConfigFile(cfg)
	} else {
		v.SetConfigFile("configs/config.yaml")
	}
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	// ignore error: if file missing, still return v
	_ = v.ReadInConfig()
	return v, nil
}

// ProvideConfig 从 viper 中构建 Config
func ProvideConfig(v *viper.Viper) (*Config, error) {
	c := &Config{}
	if v == nil {
		c.DatabasePath = "./demo_logs.db"
		c.ServerAddr = ":8080"
		return c, nil
	}
	// support dot keys or nested maps via mapstructure
	if v.IsSet("database.path") {
		c.DatabasePath = v.GetString("database.path")
	} else if v.IsSet("database") {
		if sub := v.GetStringMap("database"); sub != nil {
			if p, ok := sub["path"].(string); ok {
				c.DatabasePath = p
			}
		}
	}
	if v.IsSet("server.addr") {
		c.ServerAddr = v.GetString("server.addr")
	} else if v.IsSet("server") {
		if sub := v.GetStringMap("server"); sub != nil {
			if p, ok := sub["addr"].(string); ok {
				c.ServerAddr = p
			}
		}
	}
	if env := os.Getenv("SERVER_ADDR"); env != "" {
		c.ServerAddr = env
	}
	if env := os.Getenv("DATABASE_PATH"); env != "" {
		c.DatabasePath = env
	}
	return c, nil
}

var ProviderSet = wire.NewSet(ProvideViper, ProvideConfig)
