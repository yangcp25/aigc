package main

import (
	// 如果你没有引入 log，可以先用内置的 fmt 或者直接删掉打印
	"gitlab.hudonggz.cn/yangchunping/go-infra/log"
)

func main() {
	// 1. 🌟 初始化全局基建日志
	log.Init(log.Config{Level: "info", Format: "console"})
	defer log.Sync()
	
	// 重点在这里：把里面的参数清空，变成 ()
	app, err := InitApp()
	if err != nil {
		panic(err)
	}

	log.Info("🚀 AIGC 后端服务启动准备完毕...")
	app.Start()
}
