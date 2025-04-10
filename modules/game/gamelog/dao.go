package gamelog

import (
	"log"
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 数据迁移
func autoMigrate() error {
	return app.DB.AutoMigrate(&models.GameLog{})
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
	return nil
}

func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("gamelog:permission") {
		log.Println("游戏日志模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建游戏日志相关权限
		permissions := []models.Permission{
			{
				Name:        "游戏日志列表",
				Code:        "gamelog:view",
				Description: "查看游戏日志列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "游戏日志收集",
				Code:        "gamelog:create",
				Description: "收集游戏日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "游戏日志删除",
				Code:        "gamelog:delete",
				Description: "删除游戏日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "gamelog:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
