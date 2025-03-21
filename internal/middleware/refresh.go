package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	needRefresh bool
	mu          sync.Mutex
)

// RefreshMiddleware 处理编辑后的刷新标记
func RefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查当前请求是否设置了刷新标记
		if value, exists := c.Get("needRefresh"); exists && value.(bool) {
			mu.Lock()
			needRefresh = true
			mu.Unlock()
		}
		c.Next()
	}
}

// CheckRefreshMiddleware 检查是否需要刷新数据
func CheckRefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mu.Lock()
		if needRefresh {
			c.Header("X-Refresh-Required", "true")
			needRefresh = false
		}
		mu.Unlock()
		c.Next()
	}
}
