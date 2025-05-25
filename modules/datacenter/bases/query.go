package bases

import (
	"encoding/json"
	"fmt"

	"github.com/andycai/goapi/models"
)

// QueryEntities 查询实体列表
func QueryEntities(limit, page int, search string) ([]*EntityResponse, int64, error) {
	entities, total, err := getEntities(limit, page, search)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*EntityResponse, len(entities))
	for i, entity := range entities {
		responses[i] = &EntityResponse{
			ID:          entity.ID,
			Name:        entity.Name,
			TableName:   entity.TableName,
			Description: entity.Description,
			CreatedAt:   entity.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   entity.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:   entity.CreatedBy,
			UpdatedBy:   entity.UpdatedBy,
		}
	}

	return responses, total, nil
}

// QueryEntity 查询单个实体
func QueryEntity(id uint) (*EntityResponse, error) {
	entity, err := getEntity(id)
	if err != nil {
		return nil, err
	}

	return &EntityResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		TableName:   entity.TableName,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
	}, nil
}

// QueryFields 查询字段列表
func QueryFields(entityID uint, limit, page int, search string) ([]*FieldResponse, int64, error) {
	fields, total, err := getFields(entityID, limit, page, search)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*FieldResponse, len(fields))
	for i, field := range fields {
		responses[i] = &FieldResponse{
			ID:          field.ID,
			EntityID:    field.EntityID,
			Name:        field.Name,
			Type:        field.Type,
			Length:      field.Length,
			IsNullable:  field.IsNullable,
			IsUnique:    field.IsUnique,
			Default:     field.Default,
			Description: field.Description,
			CreatedAt:   field.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   field.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:   field.CreatedBy,
			UpdatedBy:   field.UpdatedBy,
		}
	}

	return responses, total, nil
}

// QueryField 查询单个字段
func QueryField(id uint) (*FieldResponse, error) {
	field, err := getField(id)
	if err != nil {
		return nil, err
	}

	return &FieldResponse{
		ID:          field.ID,
		EntityID:    field.EntityID,
		Name:        field.Name,
		Type:        field.Type,
		Length:      field.Length,
		IsNullable:  field.IsNullable,
		IsUnique:    field.IsUnique,
		Default:     field.Default,
		Description: field.Description,
		CreatedAt:   field.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   field.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:   field.CreatedBy,
		UpdatedBy:   field.UpdatedBy,
	}, nil
}

// QueryEntityData 查询实体数据列表
func QueryEntityData(entityID uint, limit, page int) ([]*EntityDataResponse, int64, error) {
	data, total, err := getEntityData(entityID, limit, page)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*EntityDataResponse, len(data))
	for i, d := range data {
		// 解析JSON数据
		var fields map[string]interface{}
		if d.Data != "" {
			if err := json.Unmarshal([]byte(d.Data), &fields); err != nil {
				// 解析失败时使用空对象
				fields = make(map[string]interface{})
			}
		} else {
			fields = make(map[string]interface{})
		}

		responses[i] = &EntityDataResponse{
			ID:        d.ID,
			EntityID:  d.EntityID,
			Data:      fields,
			CreatedAt: d.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: d.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: d.CreatedBy,
			UpdatedBy: d.UpdatedBy,
		}
	}

	return responses, total, nil
}

// QueryEntityDataByID 查询单个实体数据
func QueryEntityDataByID(id uint) (*EntityDataResponse, error) {
	data, err := getEntityDataByID(id)
	if err != nil {
		return nil, err
	}

	// 解析JSON数据
	var fields map[string]interface{}
	if data.Data != "" {
		if err := json.Unmarshal([]byte(data.Data), &fields); err != nil {
			// 解析失败时使用空对象
			fields = make(map[string]interface{})
		}
	} else {
		fields = make(map[string]interface{})
	}

	return &EntityDataResponse{
		ID:        data.ID,
		EntityID:  data.EntityID,
		Data:      fields,
		CreatedAt: data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy: data.CreatedBy,
		UpdatedBy: data.UpdatedBy,
	}, nil
}

// ValidateEntityData 验证实体数据
func ValidateEntityData(entityID uint, data map[string]interface{}) error {
	// 查询实体的所有字段
	fields, err := getEntityFields(entityID)
	if err != nil {
		return err
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

		// 验证唯一性
		if field.IsUnique {
			var count int64
			if err := app.DB.Model(&models.EntityData{}).
				Where("entity_id = ? AND data->>'$.%s' = ?", entityID, field.Name, value).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return fmt.Errorf("字段 %s 的值必须唯一", field.Name)
			}
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
	case "date", "datetime":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("期望日期时间类型")
		}
	default:
		return fmt.Errorf("不支持的字段类型: %s", fieldType)
	}

	return nil
}
