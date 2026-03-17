package data

import (
	"context"

	"aigc/internal/conf"
	"aigc/internal/model" // 确保引入了你的 model 包
	infradb "gitlab.hudonggz.cn/yangchunping/go-infra/db"
	"gitlab.hudonggz.cn/yangchunping/go-infra/log"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

func NewDB(c *conf.Config) (*infradb.DB, error) {
	// 1. 获取数据库路径 (🚨 注意这里改成了 c.Database.Path)
	dbPath := c.Database.Path
	if dbPath == "" {
		dbPath = "./demo_logs.db" // 兜底
	}

	// 2. 🔥 核心：动态判断日志级别
	// 默认给 Info 级别（开发/测试环境，打印所有 SQL）
	gormLogLevel := logger.Info

	// 生产环境为了性能和日志容量，通常只打印慢查询和报错，设为 Warn
	if c.Env == "prod" {
		gormLogLevel = logger.Warn
	}

	// 3. 实例化本地 SQLite，并把动态计算好的 Option 塞进去
	client, err := infradb.NewSqliteDB(
		dbPath,
		infradb.WithLogLevel(gormLogLevel), // 动态注入日志级别！
	)
	if err != nil {
		log.Error("实例化本地 SQLite 失败", zap.Error(err))
		return nil, err
	}

	// 4. 表结构迁移 (🚨 换回大厂最爱的 registry 花名册模式)
	ctx := context.Background()

	// 注意：如果你还没写 GetAllModels()，可以暂时用回 &model.Log{}
	err = client.DB(ctx).AutoMigrate(model.GetAllModels()...)
	if err != nil {
		log.Error("数据库表结构自动迁移失败", zap.Error(err))
		return nil, err
	}

	log.Info("✅ 数据库实例化及表结构迁移成功", zap.String("path", dbPath), zap.String("env", c.Env))
	return client, nil
}
