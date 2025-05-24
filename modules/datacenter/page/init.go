package page

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
		&models.Page{},
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
	if app.IsInitializedModule("page:menu") {
		log.Println("[页面模块]菜单数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建页面管理菜单
		pageMenu := models.Menu{
			MenuID:     1008,
			ParentID:   enum.MenuIdSystem,
			Name:       "页面管理",
			Path:       "/admin/page",
			Icon:       "page",
			Sort:       8,
			Permission: "page:view",
			IsShow:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := tx.Create(&pageMenu).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "page:menu",
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
	if app.IsInitializedModule("page:permission") {
		log.Println("[页面模块]权限数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建页面管理相关权限
		permissions := []models.Permission{
			{
				Name:        "页面查看",
				Code:        "page:view",
				Description: "查看页面列表和详情",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "页面添加",
				Code:        "page:add",
				Description: "添加新页面",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "页面编辑",
				Code:        "page:edit",
				Description: "编辑页面",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "页面删除",
				Code:        "page:delete",
				Description: "删除页面",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记权限已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "page:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
