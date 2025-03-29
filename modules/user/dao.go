package user

import (
	"log"
	"time"

	"github.com/andycai/unitool/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

// 数据迁移
func autoMigrate() error {
	return app.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.RolePermission{}, &models.ModuleInit{})
}

// 初始化数据
func initData() error {
	if err := initMenus(); err != nil {
		return err
	}

	if err := initUser(); err != nil {
		return err
	}

	return nil
}

func initMenus() error {
	return nil
}

func initUser() error {
	// 检查是否已初始化
	if app.IsInitializedModule("user:init") {
		log.Println("用户模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建基础权限
		permissions := []models.Permission{
			// 用户
			{
				Name:        "用户列表",
				Code:        "user:view",
				Description: "查看用户列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建用户",
				Code:        "user:create",
				Description: "创建新用户",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新用户",
				Code:        "user:update",
				Description: "更新用户信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除用户",
				Code:        "user:delete",
				Description: "删除用户",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 2. 创建管理员角色
		adminRole := models.Role{
			Name:        "超级管理员",
			Description: "系统超级管理员",
			Permissions: permissions, // 赋予所有权限
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := tx.Create(&adminRole).Error; err != nil {
			return err
		}

		// 3. 创建管理员用户
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		adminUser := models.User{
			Username:  "admin",
			Password:  string(hashedPassword),
			Nickname:  "系统管理员",
			RoleID:    adminRole.ID,
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(&adminUser).Error; err != nil {
			return err
		}

		// 4. 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "user:init",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
