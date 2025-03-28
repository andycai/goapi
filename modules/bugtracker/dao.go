package bugtracker

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
		&models.Project{},
		&models.Iteration{},
		&models.Issue{},
		&models.Comment{},
	)
}

func initData() error {
	// 检查是否已初始化
	if app.IsInitializedModule("bugtracker") {
		log.Println("Bug跟踪模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建Bug跟踪相关权限
		permissions := []models.Permission{
			{
				Name:        "Bug跟踪列表",
				Code:        "bugtracker:list",
				Description: "查看Bug跟踪列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建Bug",
				Code:        "bugtracker:create",
				Description: "创建新Bug",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新Bug",
				Code:        "bugtracker:update",
				Description: "更新Bug信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除Bug",
				Code:        "bugtracker:delete",
				Description: "删除Bug",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "bugtracker",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
