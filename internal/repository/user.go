package repository

import (
	"context"

	"go-monolith-frame/internal/errcode"
	"go-monolith-frame/internal/model/entity"
	"go-monolith-frame/pkg/mysql"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepoImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建UserRepository实例
func NewUserRepository() UserRepository {
	return &userRepoImpl{
		db: mysql.DB,
	}
}

// FindByID 根据ID查找用户
func (r *userRepoImpl) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepoImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Save 保存用户
func (r *userRepoImpl) Save(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Create 创建用户
func (r *userRepoImpl) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *userRepoImpl) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Updates(user).Error
}

// Delete 删除用户
func (r *userRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}
