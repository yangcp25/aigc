//go:generate wire
//go:build wireinject
// +build wireinject

package main

import (
	"aigc/internal/api"
	"aigc/internal/conf"
	"aigc/internal/db"
	"aigc/internal/repo"
	"aigc/internal/router"
	"aigc/internal/service"
	"github.com/google/wire"
	"gitlab.hudonggz.cn/yangchunping/go-infra/httpsrv"
)

// ProvideHTTPServer 适配器：它不仅提取配置，还负责向 Wire 索要所有的 Handler，并完成路由组装
func ProvideHTTPServer(c *conf.Config, handlers *api.Handlers) *httpsrv.Server {
	// 1. 初始化纯净的 Server
	srv := httpsrv.New(c.ServerAddr)

	// 2. 调用路由包，把拿到的 handlers 整包挂载到 srv 上
	router.RegisterRoutes(srv, handlers)

	// 3. 返回组装好的终极 Server
	return srv
}

func InitApp() (*httpsrv.Server, error) {
	panic(wire.Build(
		// 0. 配置与 DB
		conf.ProviderSet,
		db.ProviderSet,

		// 1. 底座 (使用适配器)
		ProvideHTTPServer,

		// 2. 业务层
		repo.ProviderSet,
		service.ProviderSet,
		api.ProviderSet,

		// 注意：这里删掉了 router.ProviderSet，因为它已经在 ProvideHTTPServer 里被直接调用了
	))
}
