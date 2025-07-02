package middleware

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	needRefresh bool
	mu          sync.Mutex
)

// RefreshMiddleware 处理编辑后的刷新标记
func RefreshMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 检查当前请求是否设置了刷新标记
		if value := c.Locals("needRefresh"); value != nil {
			if needRefreshValue, ok := value.(bool); ok && needRefreshValue {
				mu.Lock()
				needRefresh = true
				mu.Unlock()
			}
		}
		return c.Next()
	}
}

// CheckRefreshMiddleware 检查是否需要刷新数据
func CheckRefreshMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		mu.Lock()
		if needRefresh {
			c.Set("X-Refresh-Required", "true")
			needRefresh = false
		}
		mu.Unlock()
		return c.Next()
	}
}
