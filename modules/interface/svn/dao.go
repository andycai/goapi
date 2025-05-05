package svn

import (
	"log"
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
	if app.IsInitializedModule("svn:permission") {
		log.Println("[svn模块]权限数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建svn相关权限
		permissions := []models.Permission{
			{
				Name:        "SVN列表",
				Code:        "svn:view",
				Description: "查看SVN列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN检出",
				Code:        "svn:checkout",
				Description: "检出SVN仓库",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN更新",
				Code:        "svn:update",
				Description: "更新SVN仓库",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN提交",
				Code:        "svn:commit",
				Description: "提交SVN更改",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN状态",
				Code:        "svn:status",
				Description: "查看SVN仓库状态",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN信息",
				Code:        "svn:info",
				Description: "查看SVN仓库信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN日志",
				Code:        "svn:log",
				Description: "查看SVN提交日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN还原",
				Code:        "svn:revert",
				Description: "还原SVN更改",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN添加",
				Code:        "svn:add",
				Description: "添加文件到SVN",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "SVN删除",
				Code:        "svn:delete",
				Description: "从SVN删除文件",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "svn:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
