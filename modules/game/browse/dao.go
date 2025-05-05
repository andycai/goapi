package browse

import (
	"log"
	"time"

	"github.com/andycai/goapi/enum"
	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 数据迁移
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
	// 检查是否已初始化
	if app.IsInitializedModule("browse:menu") {
		log.Println("文件浏览模块菜单已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建文件浏览菜单
		browseMenu := models.Menu{
			MenuID:     3001,
			ParentID:   enum.MenuIdGame,
			Name:       "文件浏览",
			Path:       "/admin/browse",
			Icon:       "browse",
			Sort:       1,
			Permission: "browse:view",
			IsShow:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := tx.Create(&browseMenu).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "browse:menu",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("browse:permission") {
		log.Println("文件浏览模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建文件浏览相关权限
		permissions := []models.Permission{
			{
				Name:        "文件浏览列表",
				Code:        "browse:view",
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
				Name:        "文件FTP上传",
				Code:        "browse:ftp",
				Description: "FTP上传文件",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "文件删除",
				Code:        "browse:delete",
				Description: "删除文件",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "browse:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
