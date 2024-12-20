package models

import (
	"time"
)

// ConfigProject 配置项目表
type ConfigProject struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`          // 项目名称
	Description string    `json:"description" gorm:"type:text"`           // 项目描述
	RootPath    string    `json:"root_path" gorm:"size:255;not null"`     // 配置根目录
	OutputPath  string    `json:"output_path" gorm:"size:255;not null"`   // 输出目录
	Status      string    `json:"status" gorm:"size:20;default:'active'"` // 状态：active, archived
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConfigTable 配置表定义
type ConfigTable struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ProjectID   uint      `json:"project_id" gorm:"index"`                // 所属项目ID
	Name        string    `json:"name" gorm:"size:100;not null"`          // 配置表名称
	Description string    `json:"description" gorm:"type:text"`           // 配置表描述
	FilePath    string    `json:"file_path" gorm:"size:255;not null"`     // 配置文件路径
	FileType    string    `json:"file_type" gorm:"size:20;not null"`      // 文件类型：excel, csv, json, xml, yaml
	SheetName   string    `json:"sheet_name" gorm:"size:100"`             // Excel工作表名称
	Validators  string    `json:"validators" gorm:"type:text"`            // 数据验证规则（JSON格式）
	Status      string    `json:"status" gorm:"size:20;default:'active'"` // 状态：active, archived
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConfigExport 配置导出记录
type ConfigExport struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProjectID uint      `json:"project_id" gorm:"index"`                 // 所属项目ID
	TableID   uint      `json:"table_id" gorm:"index"`                   // 配置表ID
	Format    string    `json:"format" gorm:"size:20;not null"`          // 导出格式：binary, json, bson, xml, lua, yaml
	Language  string    `json:"language" gorm:"size:20;not null"`        // 目标语言：cpp, java, go, cs, python, lua
	Status    string    `json:"status" gorm:"size:20;default:'pending'"` // 状态：pending, running, success, failed
	Output    string    `json:"output" gorm:"type:text"`                 // 导出输出
	Error     string    `json:"error" gorm:"type:text"`                  // 错误信息
	StartTime time.Time `json:"start_time"`                              // 开始时间
	EndTime   time.Time `json:"end_time"`                                // 结束时间
	Duration  int       `json:"duration"`                                // 执行时长(秒)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
