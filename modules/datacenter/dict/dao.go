package dict

import (
	"errors"
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 字典类型数据库操作

// 获取字典类型列表
func getDictTypeList(page, limit int) ([]models.DictType, int64, error) {
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
func getDictTypeByType(typeCode string) (models.DictType, error) {
	var dictType models.DictType
	if err := app.DB.Where("type = ?", typeCode).First(&dictType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dictType, ErrDictTypeNotFound
		}
		return dictType, err
	}
	return dictType, nil
}

// 根据ID获取字典类型
func getDictTypeByID(id int64) (models.DictType, error) {
	var dictType models.DictType
	if err := app.DB.First(&dictType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dictType, ErrDictTypeNotFound
		}
		return dictType, err
	}
	return dictType, nil
}

// 添加字典类型
func addDictType(dictType *models.DictType) error {
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

// 更新字典类型
func updateDictType(dictType *models.DictType) error {
	// 检查字典类型是否存在
	var existingType models.DictType
	if err := app.DB.First(&existingType, dictType.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDictTypeNotFound
		}
		return err
	}

	// 检查是否有其他记录使用相同的类型编码
	var count int64
	if err := app.DB.Model(&models.DictType{}).Where("type = ? AND id != ?", dictType.Type, dictType.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrDictTypeAlreadyExists
	}

	// 更新字段
	updates := map[string]interface{}{
		"name":       dictType.Name,
		"type":       dictType.Type,
		"remark":     dictType.Remark,
		"updated_at": time.Now(),
	}

	return app.DB.Model(dictType).Updates(updates).Error
}

// 删除字典类型
func deleteDictType(id int64) error {
	// 检查字典类型是否存在
	var dictType models.DictType
	if err := app.DB.First(&dictType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDictTypeNotFound
		}
		return err
	}

	// 删除相关的字典数据
	if err := app.DB.Where("type = ?", dictType.Type).Delete(&models.DictData{}).Error; err != nil {
		return err
	}

	// 删除字典类型
	return app.DB.Delete(&dictType).Error
}

// 字典数据数据库操作

// 根据字典类型获取字典数据列表
func getDictDataList(typeCode string, page, limit int) ([]models.DictData, int64, error) {
	var dictData []models.DictData
	var total int64

	db := app.DB.Model(&models.DictData{}).Where("type = ?", typeCode)
	db.Count(&total)

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	if err := db.Order("sort").Find(&dictData).Error; err != nil {
		return nil, 0, err
	}

	return dictData, total, nil
}

// 根据字典类型ID获取字典数据列表
func getDictDataListByTypeID(typeID int64, page, limit int) ([]models.DictData, int64, error) {
	var dictData []models.DictData
	var total int64

	db := app.DB.Model(&models.DictData{}).Where("type_id = ?", typeID)
	db.Count(&total)

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	if err := db.Order("sort").Find(&dictData).Error; err != nil {
		return nil, 0, err
	}

	return dictData, total, nil
}

// 根据ID获取字典数据
func getDictDataByID(id int64) (models.DictData, error) {
	var dictData models.DictData
	if err := app.DB.First(&dictData, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dictData, ErrDictDataNotFound
		}
		return dictData, err
	}
	return dictData, nil
}

// 获取所有字典数据（不分页）
func getAllDictData(typeCode string) ([]models.DictData, error) {
	var dictData []models.DictData
	if err := app.DB.Where("type = ?", typeCode).Order("sort").Find(&dictData).Error; err != nil {
		return nil, err
	}
	return dictData, nil
}

// 添加字典数据
func addDictData(dictData *models.DictData) error {
	// 检查字典类型是否存在
	var count int64
	if err := app.DB.Model(&models.DictType{}).Where("id = ?", dictData.TypeID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrDictTypeNotFound
	}

	dictData.CreatedAt = time.Now()
	dictData.UpdatedAt = time.Now()

	return app.DB.Create(dictData).Error
}

// 更新字典数据
func updateDictData(dictData *models.DictData) error {
	// 检查字典数据是否存在
	var existingData models.DictData
	if err := app.DB.First(&existingData, dictData.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDictDataNotFound
		}
		return err
	}

	// 更新字段
	updates := map[string]interface{}{
		"type_id":    dictData.TypeID,
		"label":      dictData.Label,
		"value":      dictData.Value,
		"sort":       dictData.Sort,
		"remark":     dictData.Remark,
		"updated_at": time.Now(),
	}

	return app.DB.Model(dictData).Updates(updates).Error
}

// 删除字典数据
func deleteDictData(id int64) error {
	// 检查字典数据是否存在
	var dictData models.DictData
	if err := app.DB.First(&dictData, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDictDataNotFound
		}
		return err
	}

	// 删除字典数据
	return app.DB.Delete(&dictData).Error
}
