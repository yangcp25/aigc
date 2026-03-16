package router

import (
	"aigc/internal/api"
	"gitlab.hudonggz.cn/yangchunping/go-infra/httpsrv"
)

// RegisterRoutes 核心组装逻辑：接收 Server 和所有的 Handler，完成挂载并返回装配好的 Server
func RegisterRoutes(srv *httpsrv.Server, handlers *api.Handlers) *httpsrv.Server {
	engine := srv.Router()

	// 可以在这里统一加全局跨域、JWT 鉴权等中间件
	// engine.Use(middleware.Cors())

	// API 分组
	v1 := engine.Group("/api/v1")
	{
		// 日志示例
		logs := v1.Group("/log")
		{
			logs.GET("/list", handlers.LogHandler.List)
		}

		// 如果以后有用户系统，接着往下加
		// user := v1.Group("/user")
		// {
		//     user.POST("/login", userHandler.Login)
		// }
	}

	// 返回组装完毕的 Server
	return srv
}
