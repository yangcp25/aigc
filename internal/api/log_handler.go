package api

import (
	"net/http"
	"strconv"

	"aigc/internal/service"
	"github.com/gin-gonic/gin"
)

// LogHandler 返回日志相关的 HTTP handler
type LogHandler struct {
	service *service.LogService
}

func NewLogHandler(s *service.LogService) *LogHandler {
	return &LogHandler{service: s}
}

// GET /api/v1/logs?limit=20
func (h *LogHandler) GetLogs(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 20
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}
	logs, err := h.service.QueryLogs(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// List 兼容新的路由命名（handlers.LogHandler.List）
func (h *LogHandler) List(c *gin.Context) {
	// 复用 GetLogs 的实现
	h.GetLogs(c)
}
