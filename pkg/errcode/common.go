package errcode

// 系统级错误 (SYS_ 前缀)
var (
	Success         = New("SUCCESS", "success")
	ErrServer       = New("SYS_SERVER_ERROR", "服务器内部错误")
	ErrParams       = New("SYS_INVALID_PARAMS", "参数错误")
	ErrNotFound     = New("SYS_NOT_FOUND", "资源不存在")
	ErrUnauthorized = New("SYS_UNAUTHORIZED", "未授权")
	ErrForbidden    = New("SYS_FORBIDDEN", "禁止访问")
	ErrTooManyReq   = New("SYS_TOO_MANY_REQUESTS", "请求过于频繁")
	ErrTimeout      = New("SYS_TIMEOUT", "请求超时")
)
