package models

import (
	"time"
)

// Parameter 参数配置
type Parameter struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Type       string    `json:"type" gorm:"size:50;not null;index"`
	Name       string    `json:"name" gorm:"size:100;not null"`
	Parameters string    `json:"parameters" gorm:"type:text"` // 存储JSON格式的参数数据
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy  uint      `json:"created_by"`
	UpdatedBy  uint      `json:"updated_by"`
}

// TableName 指定表名
func (Parameter) TableName() string {
	return "parameters"
}
