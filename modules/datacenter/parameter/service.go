package parameter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 获取参数列表
func getParameters(limit, page int, search string) ([]ParameterResponse, int64, error) {
	var parameters []models.Parameter
	var total int64
	offset := (page - 1) * limit

	query := app.DB.Model(&models.Parameter{})

	// 如果有搜索条件
	if search != "" {
		query = query.Where("name LIKE ? OR type LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Limit(limit).Offset(offset).Order("id DESC").Find(&parameters).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	response := make([]ParameterResponse, 0, len(parameters))
	for _, param := range parameters {
		paramResp := ParameterResponse{
			ID:        param.ID,
			Type:      param.Type,
			Name:      param.Name,
			CreatedAt: param.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: param.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: param.CreatedBy,
			UpdatedBy: param.UpdatedBy,
		}

		// 解析JSON参数
		var fields []ParameterField
		if param.Parameters != "" {
			if err := json.Unmarshal([]byte(param.Parameters), &fields); err != nil {
				// 解析失败时使用空数组
				fields = []ParameterField{}
			}
		} else {
			fields = []ParameterField{}
		}

		paramResp.Parameters = fields
		response = append(response, paramResp)
	}

	return response, total, nil
}

// 获取单个参数
func getParameter(id uint) (*ParameterResponse, error) {
	var param models.Parameter
	if err := app.DB.First(&param, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("参数不存在")
		}
		return nil, err
	}

	response := &ParameterResponse{
		ID:        param.ID,
		Type:      param.Type,
		Name:      param.Name,
		CreatedAt: param.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: param.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy: param.CreatedBy,
		UpdatedBy: param.UpdatedBy,
	}

	// 解析JSON参数
	var fields []ParameterField
	if param.Parameters != "" {
		if err := json.Unmarshal([]byte(param.Parameters), &fields); err != nil {
			// 解析失败时使用空数组
			fields = []ParameterField{}
		}
	} else {
		fields = []ParameterField{}
	}

	response.Parameters = fields
	return response, nil
}

// 创建参数
func createParameter(req *ParameterRequest, userID uint) (*ParameterResponse, error) {
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

	if err := app.DB.Create(&param).Error; err != nil {
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

// 更新参数
func updateParameter(id uint, req *ParameterRequest, userID uint) (*ParameterResponse, error) {
	// 查找参数
	var param models.Parameter
	if err := app.DB.First(&param, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("参数不存在")
		}
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

	if err := app.DB.Save(&param).Error; err != nil {
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

// 删除参数
func deleteParameter(id uint) error {
	// 查找参数
	var param models.Parameter
	if err := app.DB.First(&param, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("参数不存在")
		}
		return err
	}

	// 删除参数
	return app.DB.Delete(&param).Error
}

// 验证字段类型是否有效
func isValidFieldType(fieldType string) bool {
	validTypes := map[string]bool{
		"string":  true,
		"number":  true,
		"boolean": true,
	}

	return validTypes[strings.ToLower(fieldType)]
}

// 验证字段值是否与类型匹配
func isValidFieldValue(fieldType string, value interface{}) bool {
	// 如果值为nil，认为是有效的
	if value == nil {
		return true
	}

	switch strings.ToLower(fieldType) {
	case "string":
		// 对于字符串类型，接受字符串值
		_, ok := value.(string)
		return ok
	case "number":
		// 对于数字类型，接受数字值（浮点数或整数）
		switch value.(type) {
		case float64, float32, int, int32, int64, uint, uint32, uint64:
			return true
		default:
			// 尝试将JSON中的数字（通常是float64）转换为其他类型
			_, ok := value.(float64)
			return ok
		}
	case "boolean":
		// 对于布尔类型，接受布尔值
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}
