package models

import (
	"time"
)

// DictType 字典类型表
type DictType struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null;comment:字典名称"`
	Type      string    `json:"type" gorm:"size:100;not null;uniqueIndex;comment:字典类型"`
	Remark    string    `json:"remark" gorm:"type:text;comment:备注"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// DictData 字典数据表
type DictData struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	TypeID    int64     `json:"type_id" gorm:"not null;index;comment:字典类型ID"`
	Label     string    `json:"label" gorm:"size:100;not null;comment:字典标签"`
	Value     string    `json:"value" gorm:"size:100;not null;comment:字典值"`
	Sort      int       `json:"sort" gorm:"default:0;comment:排序"`
	Remark    string    `json:"remark" gorm:"type:text;comment:备注"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
