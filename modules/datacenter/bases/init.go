package bases

import (
	"log"
	"time"

	"github.com/andycai/goapi/enum"
	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// autoMigrate 自动迁移数据库表
func autoMigrate() error {
	return app.DB.AutoMigrate(
		&models.Entity{},
		&models.Field{},
		&models.EntityData{},
	)
}

// initData 初始化数据
func initData() error {
	if err := initMenus(); err != nil {
		return err
	}
	return initPermissions()
}

// initMenus 初始化菜单
func initMenus() error {
	// 检查是否已初始化
	if app.IsInitializedModule("bases:menu") {
		log.Println("[数据模块]菜单数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建基础数据管理菜单
		basesMenu := models.Menu{
			MenuID:     1010,
			ParentID:   enum.MenuIdSystem,
			Name:       "数据管理",
			Path:       "/admin/bases",
			Icon:       "bases",
			Sort:       10,
			Permission: "bases:view",
			IsShow:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := tx.Create(&basesMenu).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "bases:menu",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// initPermissions 初始化权限
func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("bases:permission") {
		log.Println("[基础数据模块]权限数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建基础数据管理相关权限
		permissions := []models.Permission{
			{
				Name:        "基础数据查看",
				Code:        "bases:view",
				Description: "查看基础数据列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "实体管理查看",
				Code:        "bases:entity:list",
				Description: "查看实体列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "实体管理创建",
				Code:        "bases:entity:create",
				Description: "创建新实体",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "实体管理编辑",
				Code:        "bases:entity:edit",
				Description: "编辑实体",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "实体管理删除",
				Code:        "bases:entity:delete",
				Description: "删除实体",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "字段管理查看",
				Code:        "bases:field:list",
				Description: "查看字段列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "字段管理创建",
				Code:        "bases:field:create",
				Description: "创建新字段",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "字段管理编辑",
				Code:        "bases:field:edit",
				Description: "编辑字段",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "字段管理删除",
				Code:        "bases:field:delete",
				Description: "删除字段",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "bases:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
