package handler

import (
	"strconv"

	"go-monolith-frame/internal/model/request"
	"go-monolith-frame/internal/service"
	"go-monolith-frame/pkg/errcode"
	"go-monolith-frame/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	// 1. 参数解析和验证
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}

	if id <= 0 {
		response.Error(c, errcode.ErrParams)
		return
	}

	// 2. 调用 service
	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 3. 返回统一格式响应
	response.Success(c, user)
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	// 1. 参数解析和验证
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 2. 调用 service
	if err := h.userService.Create(c.Request.Context(), &req); err != nil {
		response.Error(c, err)
		return
	}

	// 3. 返回统一格式响应
	response.Success(c, nil)
}
