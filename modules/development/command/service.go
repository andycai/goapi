package command

import (
	"errors"
	"time"
)

// initService 初始化服务
func initService() {
	// 初始化服务逻辑
}

// ExecuteCommand 执行命令
func ExecuteCommand(id uint) (*CommandExecution, error) {
	cmd, err := GetCommandByID(id)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, errors.New("命令不存在")
	}

	now := time.Now()
	// 创建执行记录
	execution := &CommandExecution{
		CommandID:  id,
		Status:     CommandStatusPending,
		Script:     cmd.Script,
		StartedAt:  &now,
		ExecutedBy: 0, // TODO: 添加执行者ID
	}

	if err := SaveCommandExecution(execution); err != nil {
		return nil, err
	}

	// 更新状态为运行中
	execution.Status = CommandStatusRunning
	if err := UpdateCommandExecution(execution.ID, map[string]interface{}{
		"status": execution.Status,
	}); err != nil {
		return nil, err
	}

	// TODO: 实际执行命令的逻辑
	execution.Status = CommandStatusCompleted
	execution.CompletedAt = &now
	execution.Output = "命令执行成功"

	// 更新执行结果
	if err := UpdateCommandExecution(execution.ID, map[string]interface{}{
		"status":       execution.Status,
		"completed_at": execution.CompletedAt,
		"output":       execution.Output,
	}); err != nil {
		return nil, err
	}

	return execution, nil
}

// GetCommandExecutions 获取命令执行记录
func GetCommandExecutions(commandID uint) ([]CommandExecution, error) {
	return GetCommandExecutionsByCommandID(commandID)
}
