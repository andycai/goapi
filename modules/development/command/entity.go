package command

import (
	"time"

	"gorm.io/gorm"
)

// CommandStatus represents the status of a command execution
type CommandStatus string

const (
	CommandStatusPending   CommandStatus = "pending"   // 待执行
	CommandStatusRunning   CommandStatus = "running"   // 执行中
	CommandStatusCompleted CommandStatus = "completed" // 执行完成
	CommandStatusFailed    CommandStatus = "failed"    // 执行失败
)

// Command represents a command script
type Command struct {
	gorm.Model
	Name        string   `gorm:"type:varchar(100);not null;comment:命令名称" json:"name"`
	Description string   `gorm:"type:varchar(500);comment:命令描述" json:"description"`
	Script      string   `gorm:"type:text;not null;comment:命令脚本" json:"script"`
	Category    string   `gorm:"type:varchar(50);comment:命令分类" json:"category"`
	Tags        []string `gorm:"type:json;comment:命令标签" json:"tags"`
	CreatedBy   uint     `gorm:"comment:创建者ID" json:"createdBy"`
	UpdatedBy   uint     `gorm:"comment:更新者ID" json:"updatedBy"`
}

// CommandExecution represents a command execution record
type CommandExecution struct {
	gorm.Model
	CommandID   uint          `gorm:"not null;comment:命令ID" json:"commandId"`
	Status      CommandStatus `gorm:"type:varchar(20);not null;default:'pending';comment:执行状态" json:"status"`
	Script      string        `gorm:"type:text;not null;comment:执行脚本" json:"script"`
	StartedAt   *time.Time    `gorm:"comment:开始时间" json:"startedAt"`
	CompletedAt *time.Time    `gorm:"comment:完成时间" json:"completedAt"`
	Output      string        `gorm:"type:text;comment:执行输出" json:"output"`
	Error       string        `gorm:"type:text;comment:错误信息" json:"error"`
	ExecutedBy  uint          `gorm:"comment:执行者ID" json:"executedBy"`
}

// TableName specifies the table name for Command
func (Command) TableName() string {
	return "commands"
}

// TableName specifies the table name for CommandExecution
func (CommandExecution) TableName() string {
	return "command_executions"
}
