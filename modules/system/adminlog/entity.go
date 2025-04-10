package adminlog

import "time"

// AdminLogEntry 存储操作日志信息的结构体
type AdminLogEntry struct {
	UserID     uint      // 用户ID
	Username   string    // 用户名
	Action     string    // 操作类型
	Resource   string    // 资源类型
	ResourceID uint      // 资源ID
	Details    string    // 操作详情
	IP         string    // IP地址
	UserAgent  string    // 用户代理
	CreatedAt  time.Time // 创建时间
}
