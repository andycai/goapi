package bases

import (
	"github.com/andycai/goapi/models"
)

// autoMigrate 自动迁移数据库表
func autoMigrate() error {
	return app.DB.AutoMigrate(
		&models.Entity{},
		&models.Field{},
		&models.EntityData{},
	)
}

// initData 初始化数据
func initData() error {
	// 检查是否已存在基础实体
	var count int64
	if err := app.DB.Model(&models.Entity{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	// 创建基础实体
	entity := &models.Entity{
		Name:        "用户",
		TableName:   "users",
		Description: "系统用户",
	}

	if err := app.DB.Create(entity).Error; err != nil {
		return err
	}

	// 创建基础字段
	fields := []models.Field{
		{
			EntityID:    entity.ID,
			Name:        "username",
			Type:        "string",
			Length:      50,
			IsNullable:  false,
			IsUnique:    true,
			Description: "用户名",
		},
		{
			EntityID:    entity.ID,
			Name:        "password",
			Type:        "string",
			Length:      100,
			IsNullable:  false,
			Description: "密码",
		},
		{
			EntityID:    entity.ID,
			Name:        "email",
			Type:        "string",
			Length:      100,
			IsNullable:  false,
			IsUnique:    true,
			Description: "邮箱",
		},
		{
			EntityID:    entity.ID,
			Name:        "status",
			Type:        "int",
			IsNullable:  false,
			Default:     "1",
			Description: "状态：1-正常，0-禁用",
		},
	}

	for _, field := range fields {
		if err := app.DB.Create(&field).Error; err != nil {
			return err
		}
	}

	return nil
}
