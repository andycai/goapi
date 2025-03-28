package browse

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
	if app.IsInitializedModule("browse") {
		log.Println("文件浏览模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建文件浏览相关权限
		permissions := []models.Permission{
			{
				Name:        "文件浏览列表",
				Code:        "browse:list",
				Description: "查看文件列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "文件下载",
				Code:        "browse:download",
				Description: "下载文件",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "文件上传",
				Code:        "browse:upload",
				Description: "上传文件",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "browse",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
