package repo

import (
	"context"
	"database/sql"
	"time"

	"aigc/internal/model"
	"github.com/google/wire"
)

// Repos 聚合了所有的业务 Repo
type Repos struct {
	LogRepo LogRepo
}

// LogRepo 日志仓库接口
type LogRepo interface {
	Query(ctx context.Context, limit int) ([]*model.Log, error)
	Seed(ctx context.Context) error
}

// sqliteLogRepo 基于 *sql.DB 的实现
type sqliteLogRepo struct {
	db *sql.DB
}

func NewSqliteLogRepo(db *sql.DB) LogRepo {
	return &sqliteLogRepo{db: db}
}

// ensure table exists
func (r *sqliteLogRepo) ensureTable(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT,
		level TEXT,
		created_at TEXT
	)
	`)
	return err
}

func (r *sqliteLogRepo) Seed(ctx context.Context) error {
	if err := r.ensureTable(ctx); err != nil {
		return err
	}
	// insert some sample logs
	_, err := r.db.ExecContext(ctx, `INSERT INTO logs (message, level, created_at) VALUES
	(?, ?, ?),
	(?, ?, ?)
	`,
		"system started", "info", time.Now().Format(time.RFC3339Nano),
		"failed to connect", "error", time.Now().Format(time.RFC3339Nano),
	)
	return err
}

func (r *sqliteLogRepo) Query(ctx context.Context, limit int) ([]*model.Log, error) {
	if err := r.ensureTable(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id, message, level, created_at FROM logs ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*model.Log, 0)
	for rows.Next() {
		var l model.Log
		var t string
		if err := rows.Scan(&l.ID, &l.Message, &l.Level, &t); err != nil {
			return nil, err
		}
		// parse time if possible
		if parsed, err := time.Parse(time.RFC3339Nano, t); err == nil {
			l.CreatedAt = parsed
		} else {
			l.CreatedAt = time.Now()
		}
		out = append(out, &l)
	}
	return out, rows.Err()
}

var ProviderSet = wire.NewSet(
	NewSqliteLogRepo,
	wire.Struct(new(Repos), "*"),
)
