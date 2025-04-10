package git

import (
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

func autoMigrate() error {
	return nil
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
	if app.IsInitializedModule("git:permission") {
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建git相关权限
		permissions := []models.Permission{
			{
				Name:        "Git列表",
				Code:        "git:view",
				Description: "查看Git列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git克隆",
				Code:        "git:clone",
				Description: "克隆Git仓库",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git拉取",
				Code:        "git:pull",
				Description: "拉取Git仓库更新",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git推送",
				Code:        "git:push",
				Description: "推送Git仓库更新",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git状态",
				Code:        "git:status",
				Description: "查看Git仓库状态",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git日志",
				Code:        "git:log",
				Description: "查看Git提交日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git提交",
				Code:        "git:commit",
				Description: "提交Git更改",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git检出",
				Code:        "git:checkout",
				Description: "检出Git分支或标签",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git分支",
				Code:        "git:branch",
				Description: "管理Git分支",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git合并",
				Code:        "git:merge",
				Description: "合并Git分支",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Git重置",
				Code:        "git:reset",
				Description: "重置Git仓库状态",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "git:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
