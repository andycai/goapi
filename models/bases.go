package models

import (
	"time"
)

// Entity 实体定义
type Entity struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`       // 实体名称
	TableName   string    `gorm:"size:100;not null" json:"table_name"` // 表名
	Description string    `gorm:"size:500" json:"description"`         // 描述
	CreatedAt   time.Time `json:"created_at"`                          // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                          // 更新时间
	CreatedBy   uint      `gorm:"not null" json:"created_by"`          // 创建者
	UpdatedBy   uint      `gorm:"not null" json:"updated_by"`          // 更新者
}

// Field 字段定义
type Field struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	EntityID    uint      `gorm:"not null" json:"entity_id"`     // 所属实体ID
	Name        string    `gorm:"size:100;not null" json:"name"` // 字段名称
	Type        string    `gorm:"size:50;not null" json:"type"`  // 字段类型
	Length      int       `json:"length"`                        // 字段长度
	IsNullable  bool      `gorm:"not null" json:"is_nullable"`   // 是否可为空
	IsUnique    bool      `gorm:"not null" json:"is_unique"`     // 是否唯一
	Default     string    `gorm:"size:500" json:"default"`       // 默认值
	Description string    `gorm:"size:500" json:"description"`   // 字段描述
	CreatedAt   time.Time `json:"created_at"`                    // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                    // 更新时间
	CreatedBy   uint      `gorm:"not null" json:"created_by"`    // 创建者
	UpdatedBy   uint      `gorm:"not null" json:"updated_by"`    // 更新者
}

// EntityData 实体数据
type EntityData struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	EntityID  uint      `gorm:"not null" json:"entity_id"`      // 所属实体ID
	Data      string    `gorm:"type:json;not null" json:"data"` // 数据
	CreatedAt time.Time `json:"created_at"`                     // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                     // 更新时间
	CreatedBy uint      `gorm:"not null" json:"created_by"`     // 创建者
	UpdatedBy uint      `gorm:"not null" json:"updated_by"`     // 更新者
}
