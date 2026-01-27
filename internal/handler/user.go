package handler

import (
	"strconv"

	"go-monolith-frame/internal/middleware"
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

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	// 1. 参数解析和验证
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 2. 调用 service
	loginResp, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 3. 返回统一格式响应
	response.Success(c, loginResp)
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
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

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 2. 调用 service
	if err := h.userService.Update(c.Request.Context(), id, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 3. 返回统一格式响应
	response.Success(c, nil)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
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
	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	// 3. 返回统一格式响应
	response.Success(c, nil)
}

// ListUsers 查询用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 1. 参数解析和验证
	var req request.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 设置默认值
	req.SetDefaults()

	// 2. 调用 service
	listResp, err := h.userService.List(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 3. 返回统一格式响应
	response.Success(c, listResp)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
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

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 2. 权限验证：只能修改自己的密码
	currentUserID := middleware.GetUserID(c)
	if currentUserID != id {
		response.Error(c, errcode.ErrUnauthorized.WithMsg("只能修改自己的密码"))
		return
	}

	// 3. 调用 service
	if err := h.userService.ChangePassword(c.Request.Context(), id, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 4. 返回统一格式响应
	response.Success(c, nil)
}
