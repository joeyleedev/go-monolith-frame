package entity

import (
	"time"

	"gorm.io/gorm"
)

// 用户状态常量
const (
	UserStatusDisabled = 0 // 禁用
	UserStatusEnabled  = 1 // 启用
)

type User struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"type:varchar(50);not null" json:"username"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // 密码不返回给前端
	Status    int8           `gorm:"type:tinyint;not null;default:1;comment:用户状态:0禁用,1启用" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除支持
}
