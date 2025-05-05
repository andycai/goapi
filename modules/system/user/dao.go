package user

import (
	"log"
	"time"

	"github.com/andycai/goapi/enum"
	"github.com/andycai/goapi/models"
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
	// 检查是否已初始化
	if app.IsInitializedModule("user:menu") {
		log.Println("[用户模块]菜单数据已初始化，跳过")
		return nil
	}

	// 初始化数据
	now := time.Now()

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		parentMenus := []*models.Menu{
			{
				MenuID:     enum.MenuIdSystem,
				ParentID:   0,
				Name:       "系统管理",
				Path:       "/admin/system",
				Icon:       "system",
				Sort:       1,
				Permission: "",
				IsShow:     true,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				MenuID:     enum.MenuIdTools,
				ParentID:   0,
				Name:       "系统工具",
				Path:       "/admin/tools",
				Icon:       "tools",
				Sort:       2,
				Permission: "",
				IsShow:     true,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				MenuID:     enum.MenuIdGame,
				ParentID:   0,
				Name:       "游戏管理",
				Path:       "/admin/game",
				Icon:       "game",
				Sort:       3,
				Permission: "",
				IsShow:     true,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				MenuID:     enum.MenuIdWebApp,
				ParentID:   0,
				Name:       "Web应用",
				Path:       "/admin/webapp",
				Icon:       "webapp",
				Sort:       4,
				Permission: "",
				IsShow:     true,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}
		if err := tx.CreateInBatches(parentMenus, len(parentMenus)).Error; err != nil {
			return err
		}

		// 创建用户管理菜单
		userMenu := models.Menu{
			MenuID:     1001,
			ParentID:   enum.MenuIdSystem,
			Name:       "用户管理",
			Path:       "/admin/users",
			Icon:       "user",
			Sort:       1,
			Permission: "user:view",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if err := tx.Create(&userMenu).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "user:menu",
			Initialized: 1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func initUser() error {
	// 检查是否已初始化
	if app.IsInitializedModule("user:init") {
		log.Println("[用户模块]权限数据已初始化，跳过")
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
