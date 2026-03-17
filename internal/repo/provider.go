package repo

import (
	"github.com/google/wire"
)

type Repos struct {
	LogRepo LogRepo
}

var ProviderSet = wire.NewSet(
	NewSqliteLogRepo,
	wire.Struct(new(Repos), "*"),
)
