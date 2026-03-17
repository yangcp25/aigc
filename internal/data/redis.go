package data

import (
	"fmt"

	"aigc/internal/conf"
	"github.com/alicebob/miniredis/v2"
	"gitlab.hudonggz.cn/yangchunping/go-infra/cache"
	"gitlab.hudonggz.cn/yangchunping/go-infra/log"
	"go.uber.org/zap"
)

// ProvideRedisCache 实例化缓存 (根据环境智能切换)
func ProvideRedisCache(c *conf.Config) (cache.Cache, error) {
	var addrs string
	var password string

	// 1. 核心路由：根据环境判断用真 Redis 还是假 Redis
	if c.Env == "prod" || c.Env == "test" {
		// 线上/测试环境：连真实的 Redis (假设你的 Config 里配置了这些字段)
		// 注意：如果你配置里是 URL 和 Port，这里灵活调整拼接方式
		addrs = fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
		password = c.Redis.Password
		log.Info("🔗 连接外部真实 Redis", zap.String("addrs", addrs), zap.String("env", c.Env))

	} else {
		// 本地开发环境：起飞！内嵌 MiniRedis
		mr, err := miniredis.Run()
		if err != nil {
			log.Error("启动 MiniRedis 失败", zap.Error(err))
			return nil, err
		}
		addrs = mr.Addr()
		password = "" // MiniRedis 默认无密码
		log.Info("🚀 本地开发模式: MiniRedis 启动成功", zap.String("addr", addrs))
	}

	// 2. 组装配置项 (复用你的 infra 逻辑)
	ops := []cache.RedisOption{
		cache.WithAddrs(addrs),
		cache.WithAuth("", password),
		cache.WithPoolSize(20),
	}

	// 如果你的 Config 里有 DB 库号配置，也可以在这里动态 append
	// if c.RedisDB != 0 {
	// 	ops = append(ops, cache.WithDB(c.RedisDB))
	// }

	// 3. 实例化底层 Redis 客户端
	redisInstance, err := cache.NewRedis(ops...)
	if err != nil {
		log.Error("连接 Redis 失败", zap.Error(err))
		return nil, err
	}

	// 4. 包装成你的业务 Cache 接口
	cacheInstance := cache.NewRedisCache(redisInstance, "auto_token")
	log.Info("✅ Redis 缓存层装配完毕")

	return cacheInstance, nil
}
