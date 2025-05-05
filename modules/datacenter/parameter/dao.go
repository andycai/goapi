package parameter

import (
	"log"
	"time"

	"github.com/andycai/goapi/enum"
	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 数据访问层

func autoMigrate() error {
	return app.DB.AutoMigrate(
		&models.Parameter{},
	)
}

// 初始化数据
func initData() error {
	if err := initMenus(); err != nil {
		return err
	}

	if err := initPermissions(); err != nil {
		return err
	}

	return nil
}

func initMenus() error {
	// 检查是否已初始化
	if app.IsInitializedModule("parameter:menu") {
		log.Println("参数配置模块菜单已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建字典管理菜单
		parameterMenu := models.Menu{
			MenuID:     1007,
			ParentID:   enum.MenuIdSystem,
			Name:       "参数配置",
			Path:       "/admin/parameter",
			Icon:       "parameter",
			Sort:       7,
			Permission: "parameter:view",
			IsShow:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := tx.Create(&parameterMenu).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "parameter:menu",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("parameter:permission") {
		log.Println("参数管理模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建参数管理相关权限
		permissions := []models.Permission{
			{
				Name:        "参数管理查看",
				Code:        "parameter:view",
				Description: "查看参数列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "参数管理创建",
				Code:        "parameter:create",
				Description: "创建新参数",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "参数管理编辑",
				Code:        "parameter:edit",
				Description: "编辑参数",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "参数管理删除",
				Code:        "parameter:delete",
				Description: "删除参数",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "parameter:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
