package patch

import (
	"log"
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 数据访问层，目前暂时没有特殊的数据访问逻辑
// 所有数据库操作都在 service 层完成

func autoMigrate() error {
	return app.DB.AutoMigrate(
		&models.PatchRecord{},
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
	return nil
}

func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("patch:permission") {
		log.Println("补丁模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建补丁相关权限
		permissions := []models.Permission{
			{
				Name:        "补丁管理查看",
				Code:        "patch:view",
				Description: "查看补丁管理",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "补丁配置",
				Code:        "patch:config",
				Description: "配置补丁参数",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "补丁生成",
				Code:        "patch:generate",
				Description: "生成补丁包",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "补丁应用",
				Code:        "patch:apply",
				Description: "应用补丁包",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "patch:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
