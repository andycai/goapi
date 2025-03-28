package role

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
		&models.Role{},
		&models.RolePermission{},
	)
}

func initData() error {
	// 检查是否已初始化
	if app.IsInitializedModule("role") {
		log.Println("角色模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建角色相关权限
		permissions := []models.Permission{
			{
				Name:        "角色列表",
				Code:        "role:list",
				Description: "查看角色列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建角色",
				Code:        "role:create",
				Description: "创建新角色",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新角色",
				Code:        "role:update",
				Description: "更新角色信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除角色",
				Code:        "role:delete",
				Description: "删除角色",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "role",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
