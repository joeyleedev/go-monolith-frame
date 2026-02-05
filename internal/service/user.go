package service

import (
	"context"
	"fmt"
	"time"

	"go-monolith-frame/internal/config"
	bizerr "go-monolith-frame/internal/errcode"
	"go-monolith-frame/internal/model/entity"
	"go-monolith-frame/internal/model/request"
	"go-monolith-frame/internal/model/response"
	"go-monolith-frame/internal/repository"
	"go-monolith-frame/pkg/cache"
	"go-monolith-frame/pkg/errcode"
	"go-monolith-frame/pkg/utils"
)

type UserService struct {
	userRepo repository.UserRepository
	cache    cache.Cache
}

func NewUserService(repo repository.UserRepository, cache cache.Cache) *UserService {
	return &UserService{
		userRepo: repo,
		cache:    cache,
	}
}

// GetByID 获取用户信息
func (s *UserService) GetByID(ctx context.Context, id int64) (*response.UserResponse, error) {
	// 先查缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	var cachedUser entity.User
	err := s.cache.GetObject(ctx, cacheKey, &cachedUser)

	if err == nil {
		// 缓存命中
		return &response.UserResponse{
			ID:       cachedUser.ID,
			Username: cachedUser.Username,
			Email:    cachedUser.Email,
			Status:   cachedUser.Status,
		}, nil
	}

	// 缓存未命中，查数据库
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.cache.SetObject(ctx, cacheKey, user, 10*time.Minute); err != nil {
		// 缓存写入失败不影响主流程，可以记录日志
	}

	return &response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}, nil
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *request.CreateUserRequest) error {
	// 检查邮箱是否已存在
	existUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existUser != nil {
		return bizerr.ErrUserExisted
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errcode.ErrServer.WithMsg("密码加密失败").WithDetails(err)
	}

	// 创建用户实体
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   entity.UserStatusEnabled, // 默认启用
	}

	// 保存到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		return errcode.ErrServer.WithMsg("创建用户失败").WithDetails(err)
	}

	// Cache Aside 模式：创建后删除相关缓存（如果存在）
	// 这里可以删除邮箱相关的缓存，确保数据一致性
	emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
	_ = s.cache.Delete(ctx, emailCacheKey)

	return nil
}

// Update 更新用户信息
func (s *UserService) Update(ctx context.Context, id int64, req *request.UpdateUserRequest) error {
	// 先更新数据库
	user := &entity.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
	}

	// 如果提供了status字段，则更新status
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errcode.ErrServer.WithMsg("更新用户失败").WithDetails(err)
	}

	// Cache Aside 模式：更新后删除缓存，下次读取时会重新加载最新数据
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, id int64) error {
	// 先删除数据库记录
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return errcode.ErrServer.WithMsg("删除用户失败").WithDetails(err)
	}

	// Cache Aside 模式：删除后也删除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	// 1. 根据邮箱查找用户
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errcode.ErrUnauthorized.WithMsg("用户名或密码错误")
	}

	// 2. 验证密码
	if !utils.VerifyPassword(user.Password, req.Password) {
		return nil, bizerr.ErrPassword
	}

	// 3. 检查用户状态
	if user.Status != entity.UserStatusEnabled {
		return nil, bizerr.ErrUserDisabled
	}

	// 4. 生成JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, errcode.ErrServer.WithMsg("生成Token失败").WithDetails(err)
	}

	// 5. 返回登录响应
	cfg := config.Get()
	expiresAt := time.Now().Add(cfg.JWT.ExpireTime).Unix()

	return &response.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: &response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Status:   user.Status,
		},
	}, nil
}

// List 分页查询用户列表
func (s *UserService) List(ctx context.Context, req *request.ListUsersRequest) (*response.ListUsersResponse, error) {
	// 调用repository查询
	users, total, err := s.userRepo.List(ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, errcode.ErrServer.WithMsg("查询用户列表失败").WithDetails(err)
	}

	// 转换为响应格式
	items := make([]response.UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Status:   user.Status,
		})
	}

	return &response.ListUsersResponse{
		Total: total,
		Items: items,
	}, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, id int64, req *request.ChangePasswordRequest) error {
	// 1. 查找用户
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 验证旧密码
	if !utils.VerifyPassword(user.Password, req.OldPassword) {
		return bizerr.ErrPassword
	}

	// 3. 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errcode.ErrServer.WithMsg("密码加密失败").WithDetails(err)
	}

	// 4. 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errcode.ErrServer.WithMsg("修改密码失败").WithDetails(err)
	}

	// 5. 清除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}

// UpdateStatus 更新用户状态
func (s *UserService) UpdateStatus(ctx context.Context, id int64, status int8) error {
	// 检查状态值是否有效
	if status != entity.UserStatusEnabled && status != entity.UserStatusDisabled {
		return errcode.ErrParams.WithMsg("无效的用户状态")
	}

	// 查找用户
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新状态
	user.Status = status
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errcode.ErrServer.WithMsg("更新用户状态失败").WithDetails(err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}
