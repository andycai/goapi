package parameter

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/andycai/goapi/models"
)

// isValidFieldType 验证字段类型是否有效
func isValidFieldType(fieldType string) bool {
	switch fieldType {
	case "string", "number", "boolean":
		return true
	default:
		return false
	}
}

// isValidFieldValue 验证字段值类型是否与声明的类型匹配
func isValidFieldValue(fieldType string, value interface{}) bool {
	if value == nil {
		return false
	}

	switch fieldType {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		switch v := value.(type) {
		case float64, float32, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
			return true
		case json.Number:
			_, err := v.Float64()
			return err == nil
		default:
			return false
		}
	case "boolean":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}

// CommandCreateParameter creates a new parameter
func CommandCreateParameter(req *ParameterRequest, userID uint) (*ParameterResponse, error) {
	// 验证参数
	if req.Name == "" {
		return nil, errors.New("参数名称不能为空")
	}

	if req.Type == "" {
		return nil, errors.New("参数类型不能为空")
	}

	// 验证字段名称唯一性
	fieldNames := make(map[string]bool)
	for _, field := range req.Parameters {
		if field.Name == "" {
			return nil, errors.New("字段名称不能为空")
		}

		// 检查字段名称是否重复
		if _, exists := fieldNames[field.Name]; exists {
			return nil, fmt.Errorf("字段名称 '%s' 重复", field.Name)
		}
		fieldNames[field.Name] = true

		// 验证字段类型
		if !isValidFieldType(field.Type) {
			return nil, fmt.Errorf("字段 '%s' 的类型无效，必须是 string, number 或 boolean", field.Name)
		}

		// 验证字段值类型
		if !isValidFieldValue(field.Type, field.Value) {
			return nil, fmt.Errorf("字段 '%s' 的值类型与声明的类型不匹配", field.Name)
		}
	}

	// 序列化参数字段
	parametersJSON, err := json.Marshal(req.Parameters)
	if err != nil {
		return nil, err
	}

	// 创建参数记录
	param := models.Parameter{
		Type:       req.Type,
		Name:       req.Name,
		Parameters: string(parametersJSON),
		CreatedBy:  userID,
		UpdatedBy:  userID,
	}

	if err := createParameter(&param); err != nil {
		return nil, err
	}

	// 返回创建的参数
	return &ParameterResponse{
		ID:         param.ID,
		Type:       param.Type,
		Name:       param.Name,
		Parameters: req.Parameters,
		CreatedAt:  param.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  param.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:  param.CreatedBy,
		UpdatedBy:  param.UpdatedBy,
	}, nil
}

// CommandUpdateParameter updates an existing parameter
func CommandUpdateParameter(id uint, req *ParameterRequest, userID uint) (*ParameterResponse, error) {
	// 查找参数
	param, err := getParameter(id)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if req.Name == "" {
		return nil, errors.New("参数名称不能为空")
	}

	if req.Type == "" {
		return nil, errors.New("参数类型不能为空")
	}

	// 验证字段名称唯一性
	fieldNames := make(map[string]bool)
	for _, field := range req.Parameters {
		if field.Name == "" {
			return nil, errors.New("字段名称不能为空")
		}

		// 检查字段名称是否重复
		if _, exists := fieldNames[field.Name]; exists {
			return nil, fmt.Errorf("字段名称 '%s' 重复", field.Name)
		}
		fieldNames[field.Name] = true

		// 验证字段类型
		if !isValidFieldType(field.Type) {
			return nil, fmt.Errorf("字段 '%s' 的类型无效，必须是 string, number 或 boolean", field.Name)
		}

		// 验证字段值类型
		if !isValidFieldValue(field.Type, field.Value) {
			return nil, fmt.Errorf("字段 '%s' 的值类型与声明的类型不匹配", field.Name)
		}
	}

	// 序列化参数字段
	parametersJSON, err := json.Marshal(req.Parameters)
	if err != nil {
		return nil, err
	}

	// 更新参数
	param.Type = req.Type
	param.Name = req.Name
	param.Parameters = string(parametersJSON)
	param.UpdatedBy = userID
	param.UpdatedAt = time.Now()

	if err := updateParameter(&param); err != nil {
		return nil, err
	}

	// 返回更新后的参数
	return &ParameterResponse{
		ID:         param.ID,
		Type:       param.Type,
		Name:       param.Name,
		Parameters: req.Parameters,
		CreatedAt:  param.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  param.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:  param.CreatedBy,
		UpdatedBy:  param.UpdatedBy,
	}, nil
}

// CommandDeleteParameter deletes a parameter
func CommandDeleteParameter(id uint) error {
	return deleteParameter(id)
}
