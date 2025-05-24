package parameter

import "errors"

var (
	ErrParameterNotFound = errors.New("参数不存在")
)

// ParameterField 参数字段
type ParameterField struct {
	Name  string      `json:"name"`  // 字段名称
	Type  string      `json:"type"`  // 字段类型：string, number, boolean
	Value interface{} `json:"value"` // 字段值
}

// ParameterRequest 创建/更新参数请求
type ParameterRequest struct {
	Type       string           `json:"type"`       // 参数类型
	Name       string           `json:"name"`       // 参数名称
	Parameters []ParameterField `json:"parameters"` // 参数字段列表
}

// ParameterResponse 参数响应
type ParameterResponse struct {
	ID         uint             `json:"id"`
	Type       string           `json:"type"`
	Name       string           `json:"name"`
	Parameters []ParameterField `json:"parameters"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	CreatedBy  uint             `json:"created_by"`
	UpdatedBy  uint             `json:"updated_by"`
}

// ParameterListResponse 参数列表响应
type ParameterListResponse struct {
	Parameters []ParameterResponse `json:"parameters"`
	Total      int64               `json:"total"`
}
