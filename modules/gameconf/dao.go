package gameconf

import (
	"log"
	"time"

	"github.com/andycai/unitool/models"
	"gorm.io/gorm"
)

// 数据访问层，目前暂时没有特殊的数据访问逻辑
// 所有数据库操作都在 service 层完成

func autoMigrate() error {
	return app.DB.AutoMigrate(
		&models.GameConfProject{},
		&models.GameConfTable{},
		&models.GameConfExport{},
	)
}

func initData() error {
	// 检查是否已初始化
	if app.IsInitializedModule("gameconf") {
		log.Println("游戏配置模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建游戏配置相关权限
		permissions := []models.Permission{
			{
				Name:        "游戏配置列表",
				Code:        "gameconf:list",
				Description: "查看游戏配置列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建游戏配置",
				Code:        "gameconf:create",
				Description: "创建新游戏配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新游戏配置",
				Code:        "gameconf:update",
				Description: "更新游戏配置信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除游戏配置",
				Code:        "gameconf:delete",
				Description: "删除游戏配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "gameconf",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
