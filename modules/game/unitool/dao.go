package unitool

import (
	"log"
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 数据访问层

func autoMigrate() error {
	return app.DB.AutoMigrate(
		&FindGuidLog{},
		&DuplicateGuid{},
	)
}

// 初始化数据
func initData() error {
	if err := initPermissions(); err != nil {
		return err
	}

	return nil
}

func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("unitool:permission") {
		log.Println("Unity工具模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建Unity工具相关权限
		permissions := []models.Permission{
			{
				Name:        "Unity工具查看",
				Code:        "unitool:view",
				Description: "查看Unity工具",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Unity工具使用",
				Code:        "unitool:use",
				Description: "使用Unity工具功能",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "unitool:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// SaveFindGuidLog 保存查找GUID的日志
func SaveFindGuidLog(findGuidLog *FindGuidLog) error {
	return app.DB.Create(findGuidLog).Error
}

// UpdateFindGuidLog 更新查找GUID的日志
func UpdateFindGuidLog(id uint, data map[string]interface{}) error {
	return app.DB.Model(&FindGuidLog{}).Where("id = ?", id).Updates(data).Error
}

// GetFindGuidLogByID 根据ID获取查找GUID的日志
func GetFindGuidLogByID(id uint) (*FindGuidLog, error) {
	var log FindGuidLog
	if err := app.DB.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// SaveDuplicateGuids 保存重复的GUID记录
func SaveDuplicateGuids(duplicates []DuplicateGuid) error {
	return app.DB.Create(&duplicates).Error
}

// GetDuplicateGuidsByLogID 根据日志ID获取重复的GUID记录
func GetDuplicateGuidsByLogID(logID uint) ([]DuplicateGuid, error) {
	var duplicates []DuplicateGuid
	if err := app.DB.Where("log_id = ?", logID).Find(&duplicates).Error; err != nil {
		return nil, err
	}
	return duplicates, nil
}

// GetFindGuidLogsFromDB 从数据库获取查找GUID的日志列表
func GetFindGuidLogsFromDB(page, limit int) ([]FindGuidLog, int64, error) {
	var logs []FindGuidLog
	var total int64

	// 获取总记录数
	if err := app.DB.Model(&FindGuidLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := app.DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func GetFindGuidLogById(id uint) (*FindGuidLog, error) {
	var findGuidLog FindGuidLog
	err := app.DB.First(&findGuidLog, id).Error
	return &findGuidLog, err
}
