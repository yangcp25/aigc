package repo

import (
	"aigc/internal/data"
	"context"
	"time"

	"aigc/internal/model"
)

// LogRepo 日志仓库接口
type LogRepo interface {
	Query(ctx context.Context, limit int) ([]*model.Log, error)
	Seed(ctx context.Context) error
}

// sqliteLogRepo 基于 GORM 封装的实现
type sqliteLogRepo struct {
	// 以前是 db *infradb.DB，现在直接换成大管家 Data
	data *data.Data
}

// NewSqliteLogRepo 构造函数：向 Wire 索要 *data.Data
func NewSqliteLogRepo(d *data.Data) LogRepo {
	return &sqliteLogRepo{data: d}
}

// Seed 插入演示数据
func (r *sqliteLogRepo) Seed(ctx context.Context) error {
	// 组装对象数组
	logs := []*model.Log{
		{Message: "system started", Level: "info", CreatedAt: time.Now()},
		{Message: "failed to connect", Level: "error", CreatedAt: time.Now()},
	}

	// GORM 一行代码搞定批量插入
	return r.data.DB.DB(ctx).Create(&logs).Error
}

// Query 查询日志
func (r *sqliteLogRepo) Query(ctx context.Context, limit int) ([]*model.Log, error) {
	var out []*model.Log
	// GORM 链式调用：排序 -> 限制条数 -> 查出结果并自动映射到切片
	// 彻底告别 rows.Next() 和 Scan！
	err := r.data.DB.DB(ctx).Debug().Order("created_at DESC").Limit(10).Find(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
