package command

import (
	"log"
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 数据访问层

func autoMigrate() error {
	return app.DB.AutoMigrate(
		&Command{},
		&CommandExecution{},
	)
}

// 初始化数据
func initData() error {
	if err := initPermissions(); err != nil {
		return err
	}

	return nil
}

func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("command:permission") {
		log.Println("命令模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建命令相关权限
		permissions := []models.Permission{
			{
				Name:        "命令查看",
				Code:        "command:view",
				Description: "查看命令列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "命令管理",
				Code:        "command:manage",
				Description: "管理命令（创建、执行等）",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "command:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetCommands 获取命令列表
func GetCommands(page, limit int) ([]Command, int64, error) {
	var commands []Command
	var total int64

	// 获取总记录数
	if err := app.DB.Model(&Command{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := app.DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&commands).Error; err != nil {
		return nil, 0, err
	}

	return commands, total, nil
}

// CreateCommand 创建命令
func CreateCommand(cmd *Command) error {
	return app.DB.Create(cmd).Error
}

// GetCommandByID 根据ID获取命令
func GetCommandByID(id uint) (*Command, error) {
	var cmd Command
	if err := app.DB.First(&cmd, id).Error; err != nil {
		return nil, err
	}
	return &cmd, nil
}

// UpdateCommand 更新命令
func UpdateCommand(id uint, data map[string]interface{}) error {
	return app.DB.Model(&Command{}).Where("id = ?", id).Updates(data).Error
}

// DeleteCommand 删除命令
func DeleteCommand(id uint) error {
	return app.DB.Delete(&Command{}, id).Error
}

// SaveCommandExecution 保存命令执行记录
func SaveCommandExecution(execution *CommandExecution) error {
	return app.DB.Create(execution).Error
}

// UpdateCommandExecution 更新命令执行记录
func UpdateCommandExecution(id uint, data map[string]interface{}) error {
	return app.DB.Model(&CommandExecution{}).Where("id = ?", id).Updates(data).Error
}

// GetCommandExecutionsByCommandID 根据命令ID获取执行记录
func GetCommandExecutionsByCommandID(cmdID uint) ([]CommandExecution, error) {
	var executions []CommandExecution
	if err := app.DB.Where("command_id = ?", cmdID).Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}
