package bases

import (
	"errors"
	"fmt"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// QueryEntities 获取实体列表
func QueryEntities(limit, page int, search string) ([]models.Entity, int64, error) {
	var entities []models.Entity
	var total int64

	query := app.DB.Model(&models.Entity{})

	if search != "" {
		query = query.Where("name LIKE ? OR table_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// QueryEntity 获取单个实体
func QueryEntity(id uint) (*models.Entity, error) {
	var entity models.Entity
	if err := app.DB.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, err
	}
	return &entity, nil
}

// QueryFields 获取字段列表
func QueryFields(entityID uint, limit, page int, search string) ([]models.Field, int64, error) {
	var fields []models.Field
	var total int64

	query := app.DB.Model(&models.Field{}).Where("entity_id = ?", entityID)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&fields).Error; err != nil {
		return nil, 0, err
	}

	return fields, total, nil
}

// QueryField 获取单个字段
func QueryField(id uint) (*models.Field, error) {
	var field models.Field
	if err := app.DB.First(&field, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFieldNotFound
		}
		return nil, err
	}
	return &field, nil
}

// QueryEntityData 获取实体数据列表
func QueryEntityData(entityID uint, limit, page int) ([]models.EntityData, int64, error) {
	var data []models.EntityData
	var total int64

	query := app.DB.Model(&models.EntityData{}).Where("entity_id = ?", entityID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// QueryEntityDataByID 获取单个实体数据
func QueryEntityDataByID(id uint) (*models.EntityData, error) {
	var data models.EntityData
	if err := app.DB.First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDataNotFound
		}
		return nil, err
	}
	return &data, nil
}

// ValidateEntityData 验证实体数据
func ValidateEntityData(entityID uint, data map[string]interface{}) error {
	// 获取实体的所有字段
	var fields []models.Field
	if err := app.DB.Where("entity_id = ?", entityID).Find(&fields).Error; err != nil {
		return fmt.Errorf("获取字段失败: %v", err)
	}

	// 验证每个字段
	for _, field := range fields {
		value, exists := data[field.Name]
		if !exists {
			if !field.IsNullable {
				return fmt.Errorf("字段 %s 不能为空", field.Name)
			}
			continue
		}

		// 验证字段类型
		if err := validateFieldValue(field.Type, value); err != nil {
			return fmt.Errorf("字段 %s 的值类型错误: %v", field.Name, err)
		}
	}

	return nil
}

// validateFieldValue 验证字段值类型
func validateFieldValue(fieldType string, value interface{}) error {
	switch fieldType {
	case "string", "text":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("期望字符串类型")
		}
	case "int":
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("期望整数类型")
		}
	case "float":
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("期望浮点数类型")
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("期望布尔类型")
		}
	case "date":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("期望日期类型")
		}
	case "datetime":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("期望日期时间类型")
		}
	default:
		return fmt.Errorf("不支持的字段类型: %s", fieldType)
	}

	return nil
}
