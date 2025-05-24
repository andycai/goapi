package bases

import (
	"errors"
	"fmt"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 获取实体列表
func getEntities(limit, page int, search string) ([]models.Entity, int64, error) {
	var entities []models.Entity
	var total int64
	offset := (page - 1) * limit

	query := app.DB.Model(&models.Entity{})

	// 如果有搜索条件
	if search != "" {
		query = query.Where("name LIKE ? OR table_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Limit(limit).Offset(offset).Order("id DESC").Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// 获取单个实体
func getEntity(id uint) (models.Entity, error) {
	var entity models.Entity
	if err := app.DB.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, ErrEntityNotFound
		}
		return entity, err
	}
	return entity, nil
}

// 创建实体
func createEntity(entity *models.Entity) error {
	return app.DB.Create(entity).Error
}

// 更新实体
func updateEntity(entity *models.Entity) error {
	return app.DB.Save(entity).Error
}

// 删除实体
func deleteEntity(id uint) error {
	// 查找实体
	var entity models.Entity
	if err := app.DB.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEntityNotFound
		}
		return err
	}

	return app.DB.Delete(&entity).Error
}

// 获取字段列表
func getFields(entityID uint) ([]models.Field, error) {
	var fields []models.Field
	if err := app.DB.Where("entity_id = ?", entityID).Order("id").Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

// 获取单个字段
func getField(id uint) (models.Field, error) {
	var field models.Field
	if err := app.DB.First(&field, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return field, errors.New("字段不存在")
		}
		return field, err
	}
	return field, nil
}

// 创建字段
func createField(field *models.Field) error {
	return app.DB.Create(field).Error
}

// 更新字段
func updateField(field *models.Field) error {
	return app.DB.Save(field).Error
}

// 删除字段
func deleteField(id uint) error {
	// 查找字段
	var field models.Field
	if err := app.DB.First(&field, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("字段不存在")
		}
		return err
	}

	return app.DB.Delete(&field).Error
}

// 获取实体数据列表
func getEntityData(entityID uint, limit, page int) ([]models.EntityData, int64, error) {
	var dataList []models.EntityData
	var total int64
	offset := (page - 1) * limit

	query := app.DB.Model(&models.EntityData{}).Where("entity_id = ?", entityID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Limit(limit).Offset(offset).Order("id DESC").Find(&dataList).Error; err != nil {
		return nil, 0, err
	}

	return dataList, total, nil
}

// 获取单个实体数据
func getEntityDataByID(id uint) (models.EntityData, error) {
	var data models.EntityData
	if err := app.DB.First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("数据不存在")
		}
		return data, err
	}
	return data, nil
}

// 创建实体数据
func createEntityData(data *models.EntityData) error {
	return app.DB.Create(data).Error
}

// 更新实体数据
func updateEntityData(data *models.EntityData) error {
	return app.DB.Save(data).Error
}

// 删除实体数据
func deleteEntityData(id uint) error {
	// 查找数据
	var data models.EntityData
	if err := app.DB.First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("数据不存在")
		}
		return err
	}

	return app.DB.Delete(&data).Error
}

// 验证实体数据
func validateEntityData(entityID uint, data map[string]interface{}) error {
	fields, err := getFields(entityID)
	if err != nil {
		return err
	}

	for _, field := range fields {
		value, exists := data[field.Name]
		if !exists {
			if !field.IsNullable {
				return fmt.Errorf("字段 %s 不能为空", field.Name)
			}
			continue
		}

		// 根据字段类型验证数据
		switch field.Type {
		case "string", "text":
			if str, ok := value.(string); ok {
				if field.Length > 0 && len(str) > field.Length {
					return fmt.Errorf("字段 %s 长度不能超过 %d", field.Name, field.Length)
				}
			} else {
				return fmt.Errorf("字段 %s 必须是字符串类型", field.Name)
			}
		case "int":
			if _, ok := value.(float64); !ok {
				return fmt.Errorf("字段 %s 必须是整数类型", field.Name)
			}
		case "float":
			if _, ok := value.(float64); !ok {
				return fmt.Errorf("字段 %s 必须是浮点数类型", field.Name)
			}
		case "bool":
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("字段 %s 必须是布尔类型", field.Name)
			}
		case "datetime", "date", "time":
			if _, ok := value.(string); !ok {
				return fmt.Errorf("字段 %s 必须是日期时间类型", field.Name)
			}
		case "json":
			if _, ok := value.(map[string]interface{}); !ok {
				return fmt.Errorf("字段 %s 必须是JSON对象类型", field.Name)
			}
		}
	}

	return nil
}
