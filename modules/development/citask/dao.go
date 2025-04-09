package citask

import (
	"log"
	"time"

	"github.com/andycai/unitool/models"
	"gorm.io/gorm"
)

// 数据迁移
func autoMigrate() error {
	return app.DB.AutoMigrate(&models.Task{}, &models.TaskLog{})
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
	if app.IsInitializedModule("citask:permission") {
		log.Println("任务管理模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建任务管理相关权限
		permissions := []models.Permission{
			{
				Name:        "构建任务列表",
				Code:        "citask:view",
				Description: "查看构建任务列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建构建任务",
				Code:        "citask:create",
				Description: "创建新构建任务",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新构建任务",
				Code:        "citask:update",
				Description: "更新构建任务信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除构建任务",
				Code:        "citask:delete",
				Description: "删除构建任务",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "执行构建任务",
				Code:        "citask:run",
				Description: "执行构建任务",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "citask:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
