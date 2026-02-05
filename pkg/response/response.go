package response

import (
	"errors"
	"net/http"
	"time"

	"go-monolith-frame/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success    bool            `json:"success"`
	Data       interface{}     `json:"data,omitempty"`
	Error      *ErrorInfo      `json:"error,omitempty"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
	Meta       MetaInfo        `json:"meta"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type PaginationInfo struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

type MetaInfo struct {
	TimeStamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    buildMeta(c),
	})
}

func SuccessWithPagination(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Pagination: &PaginationInfo{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
		Meta: buildMeta(c),
	})
}

func Error(c *gin.Context, err error) {
	var bizErr *errcode.AppError
	if errors.As(err, &bizErr) {
		c.JSON(http.StatusOK, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    bizErr.Code,
				Message: bizErr.Message,
				Details: bizErr.Details,
			},
			Meta: buildMeta(c),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errcode.ErrServer.Code,
			Message: errcode.ErrServer.Message,
		},
		Meta: buildMeta(c),
	})
}

func buildMeta(c *gin.Context) MetaInfo {
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}
	return MetaInfo{
		TimeStamp: time.Now().Format(time.RFC3339),
		RequestID: requestID,
	}
}
