package permission

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
		&models.Permission{},
	)
}

func initData() error {
	// 检查是否已初始化
	if app.IsInitializedModule("permission") {
		log.Println("权限模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建权限管理相关权限
		permissions := []models.Permission{
			{
				Name:        "权限列表",
				Code:        "permission:view",
				Description: "查看权限列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建权限",
				Code:        "permission:create",
				Description: "创建新权限",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新权限",
				Code:        "permission:update",
				Description: "更新权限信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除权限",
				Code:        "permission:delete",
				Description: "删除权限",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
