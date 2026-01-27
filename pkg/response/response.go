package response

import (
	"errors"
	"net/http"

	"go-monolith-frame/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func Error(c *gin.Context, err error) {
	var bizErr *errcode.AppError
	if errors.As(err, &bizErr) {
		c.JSON(http.StatusOK, Response{
			Code:    bizErr.Code,
			Message: bizErr.Message,
		})
		return
	}
	// 未知错误，返回内部服务器错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:    errcode.ErrServer.Code,
		Message: "内部服务器错误",
	})
}
