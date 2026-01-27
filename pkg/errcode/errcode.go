package errcode

import "fmt"

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	cause   error
}

func New(code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("code: %d, msg: %s, cause: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Message)
}

func (e *AppError) WithCause(err error) *AppError {
	return &AppError{Code: e.Code, Message: e.Message, cause: err}
}

func (e *AppError) WithMsg(msg string) *AppError {
	return &AppError{Code: e.Code, Message: msg, cause: e.cause}
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}
