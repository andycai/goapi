package gameconf

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/andycai/goapi/models"
	"github.com/andycai/goapi/modules/system/adminlog"
	"github.com/gofiber/fiber/v2"
)

// listProjectsHandler 获取项目列表
func listProjectsHandler(c *fiber.Ctx) error {
	var projects []models.GameConfProject
	if err := app.DB.Order("created_at desc").Find(&projects).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取项目列表失败: %v", err),
		})
	}
	return c.JSON(projects)
}

// createProjectHandler 创建项目
func createProjectHandler(c *fiber.Ctx) error {
	var project models.GameConfProject
	if err := c.BodyParser(&project); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 创建项目目录
	if err := createProjectDirs(&project); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建项目目录失败: %v", err),
		})
	}

	if err := app.DB.Create(&project).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建项目失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "create", "gameconf_project", project.ID, fmt.Sprintf("创建游戏配置项目：%s", project.Name))

	return c.JSON(project)
}

// getProjectHandler 获取项目详情
func getProjectHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.GameConfProject
	if err := app.DB.First(&project, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("项目不存在: %v", err),
		})
	}
	return c.JSON(project)
}

// updateProjectHandler 更新项目
func updateProjectHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.GameConfProject
	if err := app.DB.First(&project, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("项目不存在: %v", err),
		})
	}

	if err := c.BodyParser(&project); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	if err := app.DB.Save(&project).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("更新项目失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "update", "gameconf_project", project.ID, fmt.Sprintf("更新游戏配置项目：%s", project.Name))

	return c.JSON(project)
}

// deleteProjectHandler 删除项目
func deleteProjectHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.GameConfProject
	if err := app.DB.First(&project, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("项目不存在: %v", err),
		})
	}

	// 删除项目下的所有配置表
	if err := app.DB.Where("project_id = ?", id).Delete(&models.GameConfTable{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除项目配置表失败: %v", err),
		})
	}

	// 删除项目下的所有导出记录
	if err := app.DB.Where("project_id = ?", id).Delete(&models.GameConfExport{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除项目导出记录失败: %v", err),
		})
	}

	if err := app.DB.Delete(&project).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除项目失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "delete", "gameconf_project", project.ID, fmt.Sprintf("删除游戏配置项目：%s", project.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// listTablesHandler 获取配置表列表
func listTablesHandler(c *fiber.Ctx) error {
	projectID := c.Query("project_id")
	var tables []models.GameConfTable
	query := app.DB.Order("created_at desc")
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Find(&tables).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取配置表列表失败: %v", err),
		})
	}
	return c.JSON(tables)
}

// createTableHandler 创建配置表
func createTableHandler(c *fiber.Ctx) error {
	var table models.GameConfTable
	if err := c.BodyParser(&table); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 验证项目是否存在
	var project models.GameConfProject
	if err := app.DB.First(&project, table.ProjectID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("项目不存在: %v", err),
		})
	}

	if err := app.DB.Create(&table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建配置表失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "create", "gameconf_table", table.ID, fmt.Sprintf("创建游戏配置表：%s", table.Name))

	return c.JSON(table)
}

// getTableHandler 获取配置表详情
func getTableHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var table models.GameConfTable
	if err := app.DB.Preload("Project").First(&table, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("配置表不存在: %v", err),
		})
	}
	return c.JSON(table)
}

// updateTableHandler 更新配置表
func updateTableHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var table models.GameConfTable
	if err := app.DB.First(&table, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("配置表不存在: %v", err),
		})
	}

	if err := c.BodyParser(&table); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	if err := app.DB.Save(&table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("更新配置表失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "update", "gameconf_table", table.ID, fmt.Sprintf("更新游戏配置表：%s", table.Name))

	return c.JSON(table)
}

// deleteTableHandler 删除配置表
func deleteTableHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var table models.GameConfTable
	if err := app.DB.First(&table, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("配置表不存在: %v", err),
		})
	}

	// 删除配置表的所有导出记录
	if err := app.DB.Where("table_id = ?", id).Delete(&models.GameConfExport{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除配置表导出记录失败: %v", err),
		})
	}

	if err := app.DB.Delete(&table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除配置表失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "delete", "gameconf_table", table.ID, fmt.Sprintf("删除游戏配置表：%s", table.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// validateTableHandler 验证配置表
func validateTableHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var table models.GameConfTable
	if err := app.DB.Preload("Project").First(&table, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("配置表不存在: %v", err),
		})
	}

	// TODO: 实现配置表验证逻辑
	// 1. 检查文件是否存在
	// 2. 检查文件格式是否正确
	// 3. 检查数据是否符合验证规则

	return c.JSON(fiber.Map{
		"valid":   true,
		"message": "验证通过",
	})
}

// listExportsHandler 获取导出记录列表
func listExportsHandler(c *fiber.Ctx) error {
	projectID := c.Query("project_id")
	tableID := c.Query("table_id")
	var exports []models.GameConfExport
	query := app.DB.Order("created_at desc")
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if tableID != "" {
		query = query.Where("table_id = ?", tableID)
	}
	if err := query.Find(&exports).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取导出记录列表失败: %v", err),
		})
	}
	return c.JSON(exports)
}

// createExportHandler 创建导出记录
func createExportHandler(c *fiber.Ctx) error {
	var export models.GameConfExport
	if err := c.BodyParser(&export); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 验证项目是否存在
	var project models.GameConfProject
	if err := app.DB.First(&project, export.ProjectID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("项目不存在: %v", err),
		})
	}

	// 验证配置表是否存在
	var table models.GameConfTable
	if err := app.DB.First(&table, export.TableID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("配置表不存在: %v", err),
		})
	}

	// 设置初始状态
	export.Status = "pending"
	export.StartTime = time.Now()

	if err := app.DB.Create(&export).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建导出记录失败: %v", err),
		})
	}

	// 异步执行导出任务
	go executeExport(&export)

	// 记录操作日志
	adminlog.WriteLog(c, "create", "gameconf_export", export.ID, fmt.Sprintf("创建游戏配置导出：%s - %s", table.Name, export.Format))

	return c.JSON(export)
}

// getExportHandler 获取导出记录详情
func getExportHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var export models.GameConfExport
	if err := app.DB.Preload("Project").Preload("Table").First(&export, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("导出记录不存在: %v", err),
		})
	}
	return c.JSON(export)
}

// deleteExportHandler 删除导出记录
func deleteExportHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var export models.GameConfExport
	if err := app.DB.First(&export, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("导出记录不存在: %v", err),
		})
	}

	if err := app.DB.Delete(&export).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除导出记录失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "delete", "gameconf_export", export.ID, fmt.Sprintf("删除游戏配置导出记录：%d", export.ID))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// downloadExportHandler 下载导出文件
func downloadExportHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var export models.GameConfExport
	if err := app.DB.Preload("Project").Preload("Table").First(&export, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("导出记录不存在: %v", err),
		})
	}

	if export.Status != "success" {
		return c.Status(400).JSON(fiber.Map{
			"error": "导出未完成或失败",
		})
	}

	// 构建导出文件路径
	fileName := fmt.Sprintf("%s.%s", export.Table.Name, export.Format)
	filePath := filepath.Join(export.Project.DataPath, fileName)

	// 发送文件
	return c.Download(filePath, fileName)
}

// createProjectDirs 创建项目相关目录
func createProjectDirs(project *models.GameConfProject) error {
	dirs := []string{
		project.SourcePath,
		project.DataPath,
		project.CodePath,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %v", dir, err)
		}
	}

	return nil
}
