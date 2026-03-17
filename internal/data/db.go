package data // 或者你对应的包名

import (
	"aigc/internal/model" // 引入你的 model 包
	"context"
	infradb "gitlab.hudonggz.cn/yangchunping/go-infra/db"
	"gitlab.hudonggz.cn/yangchunping/go-infra/log"
	"go.uber.org/zap"
)

func InitDB(dbPath string) (*infradb.DB, error) {
	// 1. 实例化本地 SQLite
	client, err := infradb.NewSqliteDB(dbPath)
	if err != nil {
		log.Error("实例化本地 SQLite 失败", zap.Error(err))
		return nil, err
	}

	// 2. 🔥 核心：一键自动挡迁移所有表结构 🔥
	// 这里你需要传入一个 context，因为你的封装要求 DB(ctx) 才能拿到 *gorm.DB
	ctx := context.Background()

	// model.GetAllModels()... 会把切片打散成一个个参数传进去
	err = client.DB(ctx).AutoMigrate(&model.Log{})
	if err != nil {
		log.Error("数据库表结构自动迁移失败", zap.Error(err))
		return nil, err
	}

	log.Info("✅ 数据库实例化及表结构迁移成功", zap.String("path", dbPath))
	return client, nil
}
