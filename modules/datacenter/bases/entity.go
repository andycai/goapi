package bases

import (
	"errors"
)

// 错误定义
var (
	ErrEntityNotFound = errors.New("实体不存在")
	ErrFieldNotFound  = errors.New("字段不存在")
	ErrDataNotFound   = errors.New("数据不存在")
)

// EntityResponse 实体响应
type EntityResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	TableName   string `json:"table_name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CreatedBy   uint   `json:"created_by"`
	UpdatedBy   uint   `json:"updated_by"`
}

// FieldResponse 字段响应
type FieldResponse struct {
	ID          uint   `json:"id"`
	EntityID    uint   `json:"entity_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Length      int    `json:"length"`
	IsNullable  bool   `json:"is_nullable"`
	IsUnique    bool   `json:"is_unique"`
	Default     string `json:"default"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CreatedBy   uint   `json:"created_by"`
	UpdatedBy   uint   `json:"updated_by"`
}

// EntityDataResponse 实体数据响应
type EntityDataResponse struct {
	ID        uint                   `json:"id"`
	EntityID  uint                   `json:"entity_id"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	CreatedBy uint                   `json:"created_by"`
	UpdatedBy uint                   `json:"updated_by"`
}

// EntityRequest 实体请求
type EntityRequest struct {
	Name        string `json:"name" binding:"required"`
	TableName   string `json:"table_name" binding:"required"`
	Description string `json:"description"`
}

// FieldRequest 字段请求
type FieldRequest struct {
	EntityID    uint   `json:"entity_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Length      int    `json:"length"`
	IsNullable  bool   `json:"is_nullable"`
	IsUnique    bool   `json:"is_unique"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

// EntityDataRequest 实体数据请求
type EntityDataRequest struct {
	EntityID uint                   `json:"entity_id" binding:"required"`
	Data     map[string]interface{} `json:"data" binding:"required"`
}

// FieldTypes 支持的字段类型
var FieldTypes = []string{
	"string",
	"text",
	"int",
	"float",
	"bool",
	"date",
	"datetime",
}
