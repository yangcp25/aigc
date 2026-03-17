package data

import (
	"github.com/google/wire"
	"gitlab.hudonggz.cn/yangchunping/go-infra/cache"
	infradb "gitlab.hudonggz.cn/yangchunping/go-infra/db"
)

// Data 是整个基础设施的大管家（大礼包）
// 以后加了 Redis、Kafka，直接往这个结构体里加字段就行
type Data struct {
	DB    *infradb.DB
	Cache cache.Cache
	// Redis *redis.Client
	// Kafka *kafka.Writer
}

// ProviderSet 暴露给外层的 Wire 统筹组装
// ProviderSet 暴露给 Wire
var ProviderSet = wire.NewSet(
	NewDB,
	// NewRedis,  <-- 以后写了 Redis 的构造函数，直接扔这里
	ProvideRedisCache,

	// 核心黑科技：告诉 Wire 把 NewDB 和 NewRedis 生产出来的零件，自动塞进 Data 结构体里！
	wire.Struct(new(Data), "*"),
)
