package patch

import "time"

// PatchConfig 补丁配置
type PatchConfig struct {
	PatchPath          string `json:"patch_path"`          // 补丁路径
	ZipPath            string `json:"zip_path"`            // ZIP存放路径
	Branch             string `json:"branch"`              // 分支
	Platform           string `json:"platform"`            // 平台
	DefaultOldVersion  string `json:"default_old_version"` // 默认旧版本号
	DefaultNewVersion  string `json:"default_new_version"` // 默认新版本号
	DefaultDescription string `json:"default_description"` // 默认补丁描述
	ConfigPath         string `json:"config_path"`         // 配置文件路径
}

// PatchRecordResp 补丁记录
type PatchRecordResp struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OldVersion  string    `json:"old_version"` // 旧版本号
	NewVersion  string    `json:"new_version"` // 新版本号
	Version     string    `json:"version"`     // 版本组合
	Branch      string    `json:"branch"`      // 分支
	Platform    string    `json:"platform"`    // 平台
	PatchFile   string    `json:"patch_file"`  // 补丁文件路径
	Size        int64     `json:"size"`        // 文件大小
	FileCount   int       `json:"file_count"`  // 变更文件数量
	Status      int       `json:"status"`      // 状态: 0-待应用, 1-已应用, 2-应用失败
	Description string    `json:"description"` // 补丁描述
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FileChange 文件变更
type FileChange struct {
	Path       string `json:"path"`        // 文件路径
	ChangeType string `json:"change_type"` // 变更类型 (A:新增, M:修改, D:删除)
	Checksum   string `json:"checksum"`    // 文件校验和
	Size       int64  `json:"size"`        // 文件大小
}
