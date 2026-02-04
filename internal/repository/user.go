package repository

import (
	"context"

	"go-monolith-frame/internal/errcode"
	"go-monolith-frame/internal/model/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, pageSize int, keyword string) ([]*entity.User, int64, error)
}

type userRepoImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建UserRepository实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepoImpl{
		db: db,
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

// Delete 删除用户（软删除）
func (r *userRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

// List 分页查询用户列表
func (r *userRepoImpl) List(ctx context.Context, page, pageSize int, keyword string) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.User{})

	// 关键词搜索
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ?", keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
