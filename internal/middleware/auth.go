package middleware

import (
	"strings"

	bizerr "go-monolith-frame/internal/errcode"
	"go-monolith-frame/pkg/errcode"
	"go-monolith-frame/pkg/response"
	"go-monolith-frame/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	userIDKey   = "user_id"
	usernameKey = "username"
	emailKey    = "email"
)

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errcode.ErrUnauthorized.WithMsg("未提供认证Token"))
			c.Abort()
			return
		}

		// 2. 解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, errcode.ErrUnauthorized.WithMsg("Token格式错误"))
			c.Abort()
			return
		}

		token := parts[1]

		// 3. 验证Token
		claims, err := utils.ParseToken(token)
		if err != nil {
			if strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "token is expired") {
				response.Error(c, bizerr.ErrTokenExpired)
			} else {
				response.Error(c, bizerr.ErrInvalidToken.WithMsg(err.Error()))
			}
			c.Abort()
			return
		}

		// 4. 将用户信息注入到Context
		c.Set(userIDKey, claims.UserID)
		c.Set(usernameKey, claims.Username)
		c.Set(emailKey, claims.Email)

		c.Next()
	}
}

// GetUserID 从Context获取用户ID
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(userIDKey)
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUsername 从Context获取用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get(usernameKey)
	if !exists {
		return ""
	}
	return username.(string)
}

// GetEmail 从Context获取邮箱
func GetEmail(c *gin.Context) string {
	email, exists := c.Get(emailKey)
	if !exists {
		return ""
	}
	return email.(string)
}
