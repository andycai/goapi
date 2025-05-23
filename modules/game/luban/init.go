package luban

import (
	"log"
	"time"

	"github.com/andycai/goapi/enum"
	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 数据迁移
func autoMigrate() error {
	return app.DB.AutoMigrate(&models.ConfigProject{}, &models.ConfigTable{}, &models.ConfigExport{})
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
	if app.IsInitializedModule("luban:menu") {
		log.Println("[Luban模块]菜单数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建Luban菜单
		lubanMenu := models.Menu{
			MenuID:     3011,
			ParentID:   enum.MenuIdGame,
			Name:       "Luban配置",
			Path:       "/admin/luban",
			Icon:       "luban",
			Sort:       11,
			Permission: "luban:view",
			IsShow:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := tx.Create(&lubanMenu).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "luban:menu",
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
	if app.IsInitializedModule("luban:permission") {
		log.Println("[Luban模块]权限数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建Luban相关权限
		permissions := []models.Permission{
			{
				Name:        "Luban配置列表",
				Code:        "luban:view",
				Description: "查看Luban配置列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Luban配置管理",
				Code:        "luban:manage",
				Description: "管理Luban配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Luban导出配置",
				Code:        "luban:export",
				Description: "导出Luban配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "luban:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// 获取项目列表
func getProjectList() ([]models.ConfigProject, error) {
	var projects []models.ConfigProject
	err := app.DB.Order("created_at desc").Find(&projects).Error
	return projects, err
}

// 获取项目详情
func getProjectByID(id uint) (*models.ConfigProject, error) {
	var project models.ConfigProject
	err := app.DB.First(&project, id).Error
	return &project, err
}

// 获取配置表列表
func getTableList(projectID uint) ([]models.ConfigTable, error) {
	var tables []models.ConfigTable
	query := app.DB.Order("created_at desc")
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	err := query.Find(&tables).Error
	return tables, err
}

// 获取配置表详情
func getTableByID(id uint) (*models.ConfigTable, error) {
	var table models.ConfigTable
	err := app.DB.First(&table, id).Error
	return &table, err
}

// 获取导出记录列表
func getExportList(projectID, tableID uint) ([]models.ConfigExport, error) {
	var exports []models.ConfigExport
	query := app.DB.Order("created_at desc")
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if tableID > 0 {
		query = query.Where("table_id = ?", tableID)
	}
	err := query.Find(&exports).Error
	return exports, err
}

// 获取导出记录详情
func getExportByID(id uint) (*models.ConfigExport, error) {
	var export models.ConfigExport
	err := app.DB.First(&export, id).Error
	return &export, err
}
