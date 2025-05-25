package channel

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
		&models.Channel{},
		&models.PhysicalServer{},
		&models.ServerGroup{},
		&models.ServerGroupServer{},
		&models.Announcement{},
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
	if app.IsInitializedModule("channel:menu") {
		log.Println("[渠道模块]菜单数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建渠道菜单
		channelMenus := []*models.Menu{
			{
				MenuID:     3002,
				ParentID:   enum.MenuIdGame,
				Name:       "渠道管理",
				Path:       "/admin/channel",
				Icon:       "channel",
				Sort:       2,
				Permission: "channel:view",
				IsShow:     true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				MenuID:     3012,
				ParentID:   enum.MenuIdGame,
				Name:       "物理服务器",
				Path:       "/admin/physical_servers",
				Icon:       "physical_server",
				Sort:       12,
				Permission: "server:view",
				IsShow:     true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				MenuID:     3013,
				ParentID:   enum.MenuIdGame,
				Name:       "服务器分组",
				Path:       "/admin/server_groups",
				Icon:       "server_group",
				Sort:       13,
				Permission: "server:view",
				IsShow:     true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				MenuID:     3014,
				ParentID:   enum.MenuIdGame,
				Name:       "公告管理",
				Path:       "/admin/announcement",
				Icon:       "announcement",
				Sort:       14,
				Permission: "announcement:view",
				IsShow:     true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		if err := tx.CreateInBatches(channelMenus, len(channelMenus)).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "channel:menu",
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
	if app.IsInitializedModule("channel:permission") {
		log.Println("[渠道模块]权限数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建渠道相关权限
		permissions := []models.Permission{
			{
				Name:        "渠道查看",
				Code:        "channel:view",
				Description: "查看渠道列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "渠道管理",
				Code:        "channel:manage",
				Description: "管理渠道（创建、编辑等）",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "服务器查看",
				Code:        "server:view",
				Description: "查看服务器列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "服务器管理",
				Code:        "server:manage",
				Description: "管理服务器（创建、编辑等）",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "公告查看",
				Code:        "announcement:view",
				Description: "查看公告列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "公告管理",
				Code:        "announcement:manage",
				Description: "管理公告（创建、编辑等）",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "channel:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
