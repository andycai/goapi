package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"uniqueIndex;size:50" json:"username"`
	Password      string         `gorm:"size:100" json:"-"` // 密码不返回给前端
	Nickname      string         `gorm:"size:50" json:"nickname"`
	RoleID        uint           `json:"role_id"`
	Role          Role           `gorm:"foreignKey:RoleID" json:"role"`
	Status        int            `gorm:"default:1" json:"status"` // 1:启用 0:禁用
	LastLogin     time.Time      `json:"last_login"`
	HasChangedPwd bool           `gorm:"default:false" json:"has_changed_pwd"` // 是否已修改初始密码
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// Role 角色表
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Permission 权限表
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:50" json:"code"` // 权限编码
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RolePermission 角色-权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

// ModuleInit 模块是否初始化的数据库表
type ModuleInit struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Module      string         `gorm:"size:50;index" json:"module"`
	Initialized uint8          `gorm:"type:tinyint(1);default:0" json:"initialized" `
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
