package api

import "github.com/google/wire"

// Handlers 聚合了所有的业务 Handler (以后新增 CRUD，只管往这里加字段)
type Handlers struct {
	LogHandler *LogHandler
	// PetHandler *PetHandler
	// UserHandler *UserHandler
}

// ProviderSet 一次性暴露出所有的 Handler
var ProviderSet = wire.NewSet(
	NewLogHandler,
	// NewPetHandler,

	// 核心黑科技：告诉 Wire 自动把上面的实例，塞进 Handlers 结构体的对应字段里
	// "*" 表示自动匹配所有公有字段
	wire.Struct(new(Handlers), "*"),
)
