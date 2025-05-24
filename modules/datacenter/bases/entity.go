package bases

import (
	"errors"
)

var (
	ErrEntityNotFound = errors.New("实体不存在")
	ErrFieldNotFound  = errors.New("字段不存在")
	ErrDataNotFound   = errors.New("数据不存在")
)

// Entity 实体定义
type Entity struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`        // 实体名称
	TableName   string `json:"table_name"`  // 表名
	Description string `json:"description"` // 描述
	CreatedAt   string `json:"created_at"`  // 创建时间
	UpdatedAt   string `json:"updated_at"`  // 更新时间
	CreatedBy   uint   `json:"created_by"`  // 创建者
	UpdatedBy   uint   `json:"updated_by"`  // 更新者
}

// Field 字段定义
type Field struct {
	ID          uint   `json:"id"`
	EntityID    uint   `json:"entity_id"`   // 所属实体ID
	Name        string `json:"name"`        // 字段名称
	Type        string `json:"type"`        // 字段类型
	Length      int    `json:"length"`      // 字段长度
	IsNullable  bool   `json:"is_nullable"` // 是否可为空
	IsUnique    bool   `json:"is_unique"`   // 是否唯一
	Default     string `json:"default"`     // 默认值
	Description string `json:"description"` // 字段描述
	CreatedAt   string `json:"created_at"`  // 创建时间
	UpdatedAt   string `json:"updated_at"`  // 更新时间
	CreatedBy   uint   `json:"created_by"`  // 创建者
	UpdatedBy   uint   `json:"updated_by"`  // 更新者
}

// EntityData 实体数据
type EntityData struct {
	ID        uint                   `json:"id"`
	EntityID  uint                   `json:"entity_id"`  // 所属实体ID
	Data      map[string]interface{} `json:"data"`       // 数据
	CreatedAt string                 `json:"created_at"` // 创建时间
	UpdatedAt string                 `json:"updated_at"` // 更新时间
	CreatedBy uint                   `json:"created_by"` // 创建者
	UpdatedBy uint                   `json:"updated_by"` // 更新者
}

// EntityRequest 创建/更新实体请求
type EntityRequest struct {
	Name        string `json:"name" binding:"required"`       // 实体名称
	TableName   string `json:"table_name" binding:"required"` // 表名
	Description string `json:"description"`                   // 描述
}

// FieldRequest 创建/更新字段请求
type FieldRequest struct {
	EntityID    uint   `json:"entity_id" binding:"required"` // 所属实体ID
	Name        string `json:"name" binding:"required"`      // 字段名称
	Type        string `json:"type" binding:"required"`      // 字段类型
	Length      int    `json:"length"`                       // 字段长度
	IsNullable  bool   `json:"is_nullable"`                  // 是否可为空
	IsUnique    bool   `json:"is_unique"`                    // 是否唯一
	Default     string `json:"default"`                      // 默认值
	Description string `json:"description"`                  // 字段描述
}

// EntityDataRequest 创建/更新实体数据请求
type EntityDataRequest struct {
	EntityID uint                   `json:"entity_id" binding:"required"` // 所属实体ID
	Data     map[string]interface{} `json:"data" binding:"required"`      // 数据
}

// EntityResponse 实体响应
type EntityResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`        // 实体名称
	TableName   string `json:"table_name"`  // 表名
	Description string `json:"description"` // 描述
	CreatedAt   string `json:"created_at"`  // 创建时间
	UpdatedAt   string `json:"updated_at"`  // 更新时间
	CreatedBy   uint   `json:"created_by"`  // 创建者
	UpdatedBy   uint   `json:"updated_by"`  // 更新者
}

// FieldResponse 字段响应
type FieldResponse struct {
	ID          uint   `json:"id"`
	EntityID    uint   `json:"entity_id"`   // 所属实体ID
	Name        string `json:"name"`        // 字段名称
	Type        string `json:"type"`        // 字段类型
	Length      int    `json:"length"`      // 字段长度
	IsNullable  bool   `json:"is_nullable"` // 是否可为空
	IsUnique    bool   `json:"is_unique"`   // 是否唯一
	Default     string `json:"default"`     // 默认值
	Description string `json:"description"` // 字段描述
	CreatedAt   string `json:"created_at"`  // 创建时间
	UpdatedAt   string `json:"updated_at"`  // 更新时间
	CreatedBy   uint   `json:"created_by"`  // 创建者
	UpdatedBy   uint   `json:"updated_by"`  // 更新者
}

// EntityDataResponse 实体数据响应
type EntityDataResponse struct {
	ID        uint                   `json:"id"`
	EntityID  uint                   `json:"entity_id"`  // 所属实体ID
	Data      map[string]interface{} `json:"data"`       // 数据
	CreatedAt string                 `json:"created_at"` // 创建时间
	UpdatedAt string                 `json:"updated_at"` // 更新时间
	CreatedBy uint                   `json:"created_by"` // 创建者
	UpdatedBy uint                   `json:"updated_by"` // 更新者
}

// EntityListResponse 实体列表响应
type EntityListResponse struct {
	Entities []EntityResponse `json:"entities"`
	Total    int64            `json:"total"`
}

// FieldListResponse 字段列表响应
type FieldListResponse struct {
	Fields []FieldResponse `json:"fields"`
	Total  int64           `json:"total"`
}

// EntityDataListResponse 实体数据列表响应
type EntityDataListResponse struct {
	Data  []EntityDataResponse `json:"data"`
	Total int64                `json:"total"`
}

// FieldTypes 支持的字段类型
var FieldTypes = []string{
	"string",
	"text",
	"int",
	"float",
	"bool",
	"datetime",
	"date",
	"time",
	"json",
}
