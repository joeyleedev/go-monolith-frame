package errcode

import "go-monolith-frame/pkg/errcode"

// 用户模块 2xxxx
var (
	ErrUserNotFound = errcode.New(20001, "用户不存在")
	ErrUserExisted  = errcode.New(20002, "用户已存在")
	ErrPassword     = errcode.New(20003, "密码错误")
	ErrUserDisabled = errcode.New(20004, "用户已禁用")
	ErrInvalidToken = errcode.New(20005, "无效的Token")
	ErrTokenExpired = errcode.New(20006, "Token已过期")
)
