package adminlog

// 操作日志事件
type AddOperationLogEvent struct {
	UserID     uint
	Username   string
	IP         string
	UserAgent  string
	Action     string
	Resource   string
	ResourceID uint
	Details    string
}
