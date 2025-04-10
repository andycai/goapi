package unitool

import (
	"time"
)

// FindGuidLog 查找GUID的操作日志
type FindGuidLog struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	TargetPath      string    `gorm:"type:varchar(500);not null" json:"target_path"`      // 目标目录路径
	NotificationURL string    `gorm:"type:varchar(500);not null" json:"notification_url"` // 通知URL
	DuplicateCount  int       `gorm:"not null" json:"duplicate_count"`                    // 重复GUID数量
	Status          int       `gorm:"not null;default:0" json:"status"`                   // 状态：0-进行中，1-完成，2-失败
	Message         string    `gorm:"type:text" json:"message"`                           // 处理结果消息
	CreatedAt       time.Time `gorm:"not null" json:"created_at"`                         // 创建时间
	UpdatedAt       time.Time `gorm:"not null" json:"updated_at"`                         // 更新时间
}

// DuplicateGuid 重复的GUID记录
type DuplicateGuid struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	LogID     uint      `gorm:"not null" json:"log_id"`                      // 关联的日志ID
	GUID      string    `gorm:"type:varchar(100);not null" json:"guid"`      // GUID值
	FilePath  string    `gorm:"type:varchar(500);not null" json:"file_path"` // 文件路径
	CreatedAt time.Time `gorm:"not null" json:"created_at"`                  // 创建时间
}

// FindGuidRequest 查找GUID的请求参数
type FindGuidRequest struct {
	TargetPath      string `json:"targetPath" binding:"required"`      // 目标目录路径
	NotificationURL string `json:"notificationUrl" binding:"required"` // 通知URL
}
