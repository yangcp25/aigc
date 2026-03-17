package repo

import (
	"context"
	"time"

	"aigc/internal/model"
	"github.com/google/wire"
	// 引入你自己的 infra DB 组件
	"gitlab.hudonggz.cn/yangchunping/go-infra/db"
)

// Repos 聚合了所有的业务 Repo
type Repos struct {
	LogRepo LogRepo
}

// LogRepo 日志仓库接口
type LogRepo interface {
	Query(ctx context.Context, limit int) ([]*model.Log, error)
	Seed(ctx context.Context) error
}

// sqliteLogRepo 基于 GORM 封装的实现
type sqliteLogRepo struct {
	// 使用你 infra 库里的自定义 DB 结构体
	db *db.DB
}

// NewSqliteLogRepo 告诉 Wire 需要 *db.DB
func NewSqliteLogRepo(database *db.DB) LogRepo {
	return &sqliteLogRepo{db: database}
}

// Seed 插入演示数据
func (r *sqliteLogRepo) Seed(ctx context.Context) error {
	// 组装对象数组
	logs := []*model.Log{
		{Message: "system started", Level: "info", CreatedAt: time.Now()},
		{Message: "failed to connect", Level: "error", CreatedAt: time.Now()},
	}

	// GORM 一行代码搞定批量插入
	return r.db.DB(ctx).Create(&logs).Error
}

// Query 查询日志
func (r *sqliteLogRepo) Query(ctx context.Context, limit int) ([]*model.Log, error) {
	var out []*model.Log
	// GORM 链式调用：排序 -> 限制条数 -> 查出结果并自动映射到切片
	// 彻底告别 rows.Next() 和 Scan！
	err := r.db.DB(ctx).Debug().Order("created_at DESC").Limit(10).Find(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

// ProviderSet 暴露给 Wire
var ProviderSet = wire.NewSet(
	NewSqliteLogRepo,
	wire.Struct(new(Repos), "*"),
)
