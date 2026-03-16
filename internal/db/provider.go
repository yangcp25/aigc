package db

import (
	"database/sql"

	"aigc/internal/conf"
	"github.com/google/wire"
)

// ProvideSQLDB 从 Config 中读取 database.path 并打开 sqlite 连接
func ProvideSQLDB(c *conf.Config) (*sql.DB, error) {
	path := c.DatabasePath
	if path == "" {
		path = "./demo_logs.db"
	}
	// 使用已由其他依赖注册的 sqlite driver（避免重复注册引发 panic）
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

var ProviderSet = wire.NewSet(ProvideSQLDB)
