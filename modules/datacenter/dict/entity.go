package dict

import (
	"errors"
	"time"
)

var (
	ErrDictTypeNotFound      = errors.New("字典类型不存在")
	ErrDictTypeAlreadyExists = errors.New("字典类型已存在")
	ErrDictDataNotFound      = errors.New("字典数据不存在")
)

// DictType 字典类型
type DictType struct {
	ID        int64     `json:"id"`         // ID
	Name      string    `json:"name"`       // 名称
	Type      string    `json:"type"`       // 类型编码
	Remark    string    `json:"remark"`     // 备注
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// DictData 字典数据
type DictData struct {
	ID        int64     `json:"id"`         // ID
	Type      string    `json:"type"`       // 字典类型
	Label     string    `json:"label"`      // 标签
	Value     string    `json:"value"`      // 值
	Sort      int       `json:"sort"`       // 排序
	Remark    string    `json:"remark"`     // 备注
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}
