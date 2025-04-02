package fund

import (
	"log"
	"time"

	"github.com/andycai/unitool/models"
	"gorm.io/gorm"
)

// autoMigrate 自动迁移数据库
func autoMigrate() error {
	return app.DB.AutoMigrate(&Fund{}, &MarketIndex{})
}

// initData 初始化数据
func initData() error {
	if err := initMenus(); err != nil {
		return err
	}

	if err := initPermissions(); err != nil {
		return err
	}

	// 初始化市场指数
	indices := []MarketIndex{
		{Code: "000001", Name: "上证指数"},
		{Code: "399001", Name: "深证成指"},
		{Code: "399006", Name: "创业板指"},
		{Code: "000300", Name: "沪深300"},
	}

	for _, idx := range indices {
		var count int64
		app.DB.Model(&MarketIndex{}).Where("code = ?", idx.Code).Count(&count)
		if count == 0 {
			if err := app.DB.Create(&idx).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// initMenus 初始化菜单
func initMenus() error {
	return nil
	// 检查是否已初始化
	// if app.IsInitializedModule("fund:menu") {
	// 	log.Println("基金模块菜单已初始化，跳过")
	// 	return nil
	// }

	// // 开始事务
	// return app.DB.Transaction(func(tx *gorm.DB) error {
	// 	// 创建基金管理菜单
	// 	menu := models.Menu{
	// 		Name:      "基金管理",
	// 		Path:      "/admin/fund",
	// 		Icon:      "chart-line",
	// 		Sort:      50,
	// 		ParentID:  0,
	// 		CreatedAt: time.Now(),
	// 		UpdatedAt: time.Now(),
	// 	}

	// 	if err := tx.Create(&menu).Error; err != nil {
	// 		return err
	// 	}

	// 	// 标记模块已初始化
	// 	if err := tx.Create(&models.ModuleInit{
	// 		Module:      "fund:menu",
	// 		Initialized: 1,
	// 		CreatedAt:   time.Now(),
	// 		UpdatedAt:   time.Now(),
	// 	}).Error; err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })
}

// initPermissions 初始化权限
func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("fund:permission") {
		log.Println("基金模块权限已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建基金管理相关权限
		permissions := []models.Permission{
			{
				Name:        "基金管理查看",
				Code:        "fund:view",
				Description: "查看基金管理",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "基金管理配置",
				Code:        "fund:config",
				Description: "配置基金管理参数",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "基金数据同步",
				Code:        "fund:sync",
				Description: "同步基金数据",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "fund:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// SaveFund 保存基金信息
func SaveFund(fund *Fund) error {
	return app.DB.Save(fund).Error
}

// GetFund 获取基金信息
func GetFund(code string) (*Fund, error) {
	var fund Fund
	err := app.DB.Where("code = ?", code).First(&fund).Error
	if err != nil {
		return nil, err
	}
	return &fund, nil
}

// ListFunds 获取基金列表
func ListFunds(page, limit int) ([]Fund, int64, error) {
	var funds []Fund
	var total int64

	db := app.DB.Model(&Fund{})
	db.Count(&total)

	offset := (page - 1) * limit
	err := db.Order("updated_at desc").Offset(offset).Limit(limit).Find(&funds).Error
	if err != nil {
		return nil, 0, err
	}

	return funds, total, nil
}

// SaveMarketIndex 保存市场指数
func SaveMarketIndex(index *MarketIndex) error {
	return app.DB.Save(index).Error
}

// GetMarketIndices 获取所有市场指数
func GetMarketIndices() ([]MarketIndex, error) {
	var indices []MarketIndex
	err := app.DB.Find(&indices).Error
	if err != nil {
		return nil, err
	}
	return indices, nil
}

// DeleteAllFunds 删除所有基金数据
func DeleteAllFunds() error {
	return app.DB.Where("1 = 1").Delete(&Fund{}).Error
}

// BatchSaveFunds 批量保存基金信息
func BatchSaveFunds(funds []Fund) error {
	return app.DB.CreateInBatches(funds, 100).Error
}
