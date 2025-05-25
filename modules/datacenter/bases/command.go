package bases

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/andycai/goapi/models"
)

// isValidFieldType 验证字段类型是否有效
func isValidFieldType(fieldType string) bool {
	for _, t := range FieldTypes {
		if t == fieldType {
			return true
		}
	}
	return false
}

// isValidFieldValue 验证字段值类型是否与声明的类型匹配
func isValidFieldValue(fieldType string, value interface{}) bool {
	if value == nil {
		return false
	}

	switch fieldType {
	case "string", "text":
		_, ok := value.(string)
		return ok
	case "int":
		switch v := value.(type) {
		case float64, float32, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
			return true
		case json.Number:
			_, err := v.Float64()
			return err == nil
		default:
			return false
		}
	case "float":
		switch v := value.(type) {
		case float64, float32:
			return true
		case json.Number:
			_, err := v.Float64()
			return err == nil
		default:
			return false
		}
	case "bool":
		_, ok := value.(bool)
		return ok
	case "datetime", "date", "time":
		_, ok := value.(string)
		return ok
	case "json":
		_, ok := value.(map[string]interface{})
		return ok
	default:
		return false
	}
}

// CommandCreateEntity 创建实体
func CommandCreateEntity(req *EntityRequest, userID uint) (*EntityResponse, error) {
	// 验证参数
	if req.Name == "" {
		return nil, errors.New("实体名称不能为空")
	}

	if req.BasesName == "" {
		return nil, errors.New("表名不能为空")
	}

	// 创建实体
	entity := &models.Entity{
		Name:        req.Name,
		BasesName:   req.BasesName,
		Description: req.Description,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	if err := createEntity(entity); err != nil {
		return nil, err
	}

	// 返回创建的实体
	return &EntityResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		BasesName:   entity.BasesName,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
	}, nil
}

// CommandUpdateEntity 更新实体
func CommandUpdateEntity(id uint, req *EntityRequest, userID uint) (*EntityResponse, error) {
	// 查找实体
	entity, err := getEntity(id)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if req.Name == "" {
		return nil, errors.New("实体名称不能为空")
	}

	if req.BasesName == "" {
		return nil, errors.New("表名不能为空")
	}

	// 更新实体
	entity.Name = req.Name
	entity.BasesName = req.BasesName
	entity.Description = req.Description
	entity.UpdatedBy = userID
	entity.UpdatedAt = time.Now()

	if err := updateEntity(&entity); err != nil {
		return nil, err
	}

	// 返回更新后的实体
	return &EntityResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		BasesName:   entity.BasesName,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
	}, nil
}

// CommandDeleteEntity 删除实体
func CommandDeleteEntity(id uint) error {
	return deleteEntity(id)
}

// CommandCreateField 创建字段
func CommandCreateField(req *FieldRequest, userID uint) (*FieldResponse, error) {
	// 验证参数
	if req.Name == "" {
		return nil, errors.New("字段名称不能为空")
	}

	if req.Type == "" {
		return nil, errors.New("字段类型不能为空")
	}

	// 验证字段类型是否支持
	if !isValidFieldType(req.Type) {
		return nil, fmt.Errorf("不支持的字段类型: %s", req.Type)
	}

	// 创建字段
	field := &models.Field{
		EntityID:    req.EntityID,
		Name:        req.Name,
		Type:        req.Type,
		Length:      req.Length,
		IsNullable:  req.IsNullable,
		IsUnique:    req.IsUnique,
		Default:     req.Default,
		Description: req.Description,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	if err := createField(field); err != nil {
		return nil, err
	}

	// 返回创建的字段
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

// CommandUpdateField 更新字段
func CommandUpdateField(id uint, req *FieldRequest, userID uint) (*FieldResponse, error) {
	// 查找字段
	field, err := getField(id)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if req.Name == "" {
		return nil, errors.New("字段名称不能为空")
	}

	if req.Type == "" {
		return nil, errors.New("字段类型不能为空")
	}

	// 验证字段类型是否支持
	if !isValidFieldType(req.Type) {
		return nil, fmt.Errorf("不支持的字段类型: %s", req.Type)
	}

	// 更新字段
	field.Name = req.Name
	field.Type = req.Type
	field.Length = req.Length
	field.IsNullable = req.IsNullable
	field.IsUnique = req.IsUnique
	field.Default = req.Default
	field.Description = req.Description
	field.UpdatedBy = userID
	field.UpdatedAt = time.Now()

	if err := updateField(&field); err != nil {
		return nil, err
	}

	// 返回更新后的字段
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

// CommandDeleteField 删除字段
func CommandDeleteField(id uint) error {
	return deleteField(id)
}

// CommandCreateEntityData 创建实体数据
func CommandCreateEntityData(req *EntityDataRequest, userID uint) (*EntityDataResponse, error) {
	// 验证数据
	if err := validateEntityData(req.EntityID, req.Data); err != nil {
		return nil, err
	}

	// 序列化数据
	dataJSON, err := json.Marshal(req.Data)
	if err != nil {
		return nil, err
	}

	// 创建数据
	data := &models.EntityData{
		EntityID:  req.EntityID,
		Data:      string(dataJSON),
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if err := createEntityData(data); err != nil {
		return nil, err
	}

	// 返回创建的数据
	return &EntityDataResponse{
		ID:        data.ID,
		EntityID:  data.EntityID,
		Data:      req.Data,
		CreatedAt: data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy: data.CreatedBy,
		UpdatedBy: data.UpdatedBy,
	}, nil
}

// CommandUpdateEntityData 更新实体数据
func CommandUpdateEntityData(id uint, req *EntityDataRequest, userID uint) (*EntityDataResponse, error) {
	// 查找数据
	data, err := getEntityDataByID(id)
	if err != nil {
		return nil, err
	}

	// 验证数据
	if err := validateEntityData(req.EntityID, req.Data); err != nil {
		return nil, err
	}

	// 序列化数据
	dataJSON, err := json.Marshal(req.Data)
	if err != nil {
		return nil, err
	}

	// 更新数据
	data.Data = string(dataJSON)
	data.UpdatedBy = userID
	data.UpdatedAt = time.Now()

	if err := updateEntityData(&data); err != nil {
		return nil, err
	}

	// 返回更新后的数据
	return &EntityDataResponse{
		ID:        data.ID,
		EntityID:  data.EntityID,
		Data:      req.Data,
		CreatedAt: data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy: data.CreatedBy,
		UpdatedBy: data.UpdatedBy,
	}, nil
}

// CommandDeleteEntityData 删除实体数据
func CommandDeleteEntityData(id uint) error {
	return deleteEntityData(id)
}
