package errcode

// 系统级错误 1xxxx
var (
	Success         = New(0, "success")
	ErrServer       = New(10001, "服务器内部错误")
	ErrParams       = New(10002, "参数错误")
	ErrNotFound     = New(10003, "资源不存在")
	ErrUnauthorized = New(10004, "未授权")
	ErrForbidden    = New(10005, "禁止访问")
	ErrTooManyReq   = New(10006, "请求过于频繁")
	ErrTimeout      = New(10007, "请求超时")
)
