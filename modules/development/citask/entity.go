package citask

import "time"

// TaskProgress 任务进度
type TaskProgress struct {
	ID        uint      `json:"id"`
	TaskID    uint      `json:"task_id"`
	TaskName  string    `json:"task_name"`
	Status    string    `json:"status"`
	Output    string    `json:"output"`
	Error     string    `json:"error"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"`
	Progress  int       `json:"progress"` // 0-100
}
