package service

import (
	"context"
	"fmt"
	"time"

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
		return errcode.ErrServer.WithMsg("密码加密失败").WithCause(err)
	}

	// 创建用户实体
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// 保存到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		return errcode.ErrServer.WithMsg("创建用户失败").WithCause(err)
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

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errcode.ErrServer.WithMsg("更新用户失败").WithCause(err)
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
		return errcode.ErrServer.WithMsg("删除用户失败").WithCause(err)
	}

	// Cache Aside 模式：删除后也删除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}
