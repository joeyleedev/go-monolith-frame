package middleware

import (
	"fmt"

	"go-monolith-frame/pkg/logger"
	"go-monolith-frame/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery Panic恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic信息
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// 返回错误响应
				response.Error(c, fmt.Errorf("internal server error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
