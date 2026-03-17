package data

import (
	"github.com/google/wire"
)

// ProviderSet 暴露给外层的 Wire 统筹组装
var ProviderSet = wire.NewSet(InitDB)
