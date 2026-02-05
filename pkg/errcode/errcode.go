package errcode

import "fmt"

// AppError 应用错误
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func New(code string, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func (e *AppError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("code: %s, msg: %s, details: %v", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("code: %s, msg: %s", e.Code, e.Message)
}

func (e *AppError) WithDetails(details any) *AppError {
	return &AppError{Code: e.Code, Message: e.Message, Details: details}
}

func (e *AppError) WithMsg(msg string) *AppError {
	return &AppError{Code: e.Code, Message: msg, Details: e.Details}
}

func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}
