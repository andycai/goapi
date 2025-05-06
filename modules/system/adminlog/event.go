package adminlog

// 操作日志事件
type AddOperationLogEvent struct {
	Action     string
	Resource   string
	ResourceID uint
	Details    string
}
