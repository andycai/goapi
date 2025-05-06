package dict

import (
	"errors"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 获取字典类型列表
func QueryTypes(page, limit int) ([]models.DictType, int64, error) {
	var dictTypes []models.DictType
	var total int64

	db := app.DB.Model(&models.DictType{})
	db.Count(&total)

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	if err := db.Order("id desc").Find(&dictTypes).Error; err != nil {
		return nil, 0, err
	}

	return dictTypes, total, nil
}

// 根据类型编码获取字典类型
func QueryTypeByType(typeCode string) (models.DictType, error) {
	var dictType models.DictType
	if err := app.DB.Where("type = ?", typeCode).First(&dictType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dictType, ErrDictTypeNotFound
		}
		return dictType, err
	}
	return dictType, nil
}
