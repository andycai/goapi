package dict

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
		&models.DictType{},
		&models.DictData{},
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
	if app.IsInitializedModule("dict:menu") {
		log.Println("[字典模块]菜单数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建字典管理菜单
		dictMenu := models.Menu{
			MenuID:     1006,
			ParentID:   enum.MenuIdSystem,
			Name:       "字典管理",
			Path:       "/admin/dict",
			Icon:       "dict",
			Sort:       6,
			Permission: "dict:view",
			IsShow:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := tx.Create(&dictMenu).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "dict:menu",
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
	if app.IsInitializedModule("dict:permission") {
		log.Println("[字典模块]权限数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建字典管理相关权限
		permissions := []models.Permission{
			{
				Name:        "字典查看",
				Code:        "dict:view",
				Description: "查看字典类型和字典数据",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "字典添加",
				Code:        "dict:add",
				Description: "添加字典类型和字典数据",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "字典编辑",
				Code:        "dict:edit",
				Description: "编辑字典类型和字典数据",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "字典删除",
				Code:        "dict:delete",
				Description: "删除字典类型和字典数据",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记权限已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "dict:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
