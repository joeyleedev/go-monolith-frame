package errcode

import "go-monolith-frame/pkg/errcode"

// 用户模块错误 (USER_ 前缀)
var (
	ErrUserNotFound = errcode.New("USER_NOT_FOUND", "用户不存在")
	ErrUserExisted  = errcode.New("USER_EXISTS", "用户已存在")
	ErrPassword     = errcode.New("USER_INVALID_PASSWORD", "密码错误")
	ErrUserDisabled = errcode.New("USER_DISABLED", "用户已禁用")
	ErrInvalidToken = errcode.New("USER_INVALID_TOKEN", "无效的Token")
	ErrTokenExpired = errcode.New("USER_TOKEN_EXPIRED", "Token已过期")
)
