package dict

import (
	"time"

	"github.com/andycai/goapi/models"
)

// 添加字典类型
func CommandAddType(dictType *models.DictType) error {
	// 检查类型是否已存在
	var count int64
	if err := app.DB.Model(&models.DictType{}).Where("type = ?", dictType.Type).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrDictTypeAlreadyExists
	}

	dictType.CreatedAt = time.Now()
	dictType.UpdatedAt = time.Now()

	return app.DB.Create(dictType).Error
}
