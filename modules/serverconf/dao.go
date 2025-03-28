package serverconf

import (
	"log"
	"time"

	"github.com/andycai/unitool/models"
	"gorm.io/gorm"
)

// 数据迁移
func autoMigrate() error {
	return nil
}

// 初始化数据
func initData() error {
	// 检查是否已初始化
	if app.IsInitializedModule("serverconf") {
		log.Println("服务器配置模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建服务器配置相关权限
		permissions := []models.Permission{
			{
				Name:        "服务器配置列表",
				Code:        "serverconf:list",
				Description: "查看服务器配置列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "服务器配置管理",
				Code:        "serverconf:manage",
				Description: "管理服务器配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "服务器配置同步",
				Code:        "serverconf:sync",
				Description: "同步服务器配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "serverconf",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
