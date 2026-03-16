package service

import (
	"context"

	"aigc/internal/model"
	"aigc/internal/repo"
	"github.com/google/wire"
)

// LogService 提供日志查询的业务逻辑
type LogService struct {
	repo repo.LogRepo
}

func NewLogService(r repo.LogRepo) *LogService {
	return &LogService{repo: r}
}

func (s *LogService) QueryLogs(ctx context.Context, limit int) ([]*model.Log, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.Query(ctx, limit)
}

var ProviderSet = wire.NewSet(NewLogService)
