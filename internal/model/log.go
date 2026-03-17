package model

import "time"

// Log 简单的日志模型，用于示例
type Log struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"created_at"`
}

// 👇 加上这个终极杀器，强制指定表名！
func (Log) TableName() string {
	return "log"
}
