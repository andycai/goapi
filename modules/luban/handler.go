package luban

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/andycai/unitool/models"
	"github.com/andycai/unitool/modules/adminlog"
	"github.com/gofiber/fiber/v2"
)

// ExportProgress 导出进度
type ExportProgress struct {
	ID        uint      `json:"id"`
	ProjectID uint      `json:"project_id"`
	TableID   uint      `json:"table_id"`
	Status    string    `json:"status"`
	Output    string    `json:"output"`
	Error     string    `json:"error"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"`
	Progress  int       `json:"progress"` // 0-100
}

var (
	exportProgressMap = make(map[uint]*ExportProgress)
)

// getProjects 获取项目列表
func getProjects(c *fiber.Ctx) error {
	projects, err := getProjectList()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取项目列表失败: %v", err),
		})
	}
	return c.JSON(projects)
}

// createProject 创建项目
func createProject(c *fiber.Ctx) error {
	var project models.ConfigProject
	if err := c.BodyParser(&project); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 验证目录是否存在
	if _, err := os.Stat(project.RootPath); os.IsNotExist(err) {
		return c.Status(400).JSON(fiber.Map{
			"error": "配置根目录不存在",
		})
	}
	if _, err := os.Stat(project.OutputPath); os.IsNotExist(err) {
		return c.Status(400).JSON(fiber.Map{
			"error": "输出目录不存在",
		})
	}

	if err := app.DB.Create(&project).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建项目失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "create", "luban_project", project.ID, fmt.Sprintf("创建项目：%s", project.Name))

	return c.JSON(project)
}

// getProject 获取项目详情
func getProject(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的项目ID",
		})
	}

	project, err := getProjectByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("项目不存在: %v", err),
		})
	}

	return c.JSON(project)
}

// updateProject 更新项目
func updateProject(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的项目ID",
		})
	}

	var updates models.ConfigProject
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	project, err := getProjectByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("项目不存在: %v", err),
		})
	}

	// 验证目录是否存在
	if updates.RootPath != "" && updates.RootPath != project.RootPath {
		if _, err := os.Stat(updates.RootPath); os.IsNotExist(err) {
			return c.Status(400).JSON(fiber.Map{
				"error": "配置根目录不存在",
			})
		}
	}
	if updates.OutputPath != "" && updates.OutputPath != project.OutputPath {
		if _, err := os.Stat(updates.OutputPath); os.IsNotExist(err) {
			return c.Status(400).JSON(fiber.Map{
				"error": "输出目录不存在",
			})
		}
	}

	if err := app.DB.Model(project).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("更新项目失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "luban_project", project.ID, fmt.Sprintf("更新项目：%s", project.Name))

	return c.JSON(project)
}

// deleteProject 删除项目
func deleteProject(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的项目ID",
		})
	}

	project, err := getProjectByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("项目不存在: %v", err),
		})
	}

	if err := app.DB.Delete(project).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除项目失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "delete", "luban_project", project.ID, fmt.Sprintf("删除项目：%s", project.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// getTables 获取配置表列表
func getTables(c *fiber.Ctx) error {
	projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 32)
	tables, err := getTableList(uint(projectID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取配置表列表失败: %v", err),
		})
	}
	return c.JSON(tables)
}

// createTable 创建配置表
func createTable(c *fiber.Ctx) error {
	var table models.ConfigTable
	if err := c.BodyParser(&table); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 验证项目是否存在
	project, err := getProjectByID(table.ProjectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "项目不存在",
		})
	}

	// 验证文件是否存在
	filePath := filepath.Join(project.RootPath, table.FilePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(400).JSON(fiber.Map{
			"error": "配置文件不存在",
		})
	}

	if err := app.DB.Create(&table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建配置表失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "create", "luban_table", table.ID, fmt.Sprintf("创建配置表：%s", table.Name))

	return c.JSON(table)
}

// getTable 获取配置表详情
func getTable(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的配置表ID",
		})
	}

	table, err := getTableByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("配置表不存在: %v", err),
		})
	}

	return c.JSON(table)
}

// updateTable 更新配置表
func updateTable(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的配置表ID",
		})
	}

	var updates models.ConfigTable
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	table, err := getTableByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("配置表不存在: %v", err),
		})
	}

	// 验证项目是否存在
	project, err := getProjectByID(updates.ProjectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "项目不存在",
		})
	}

	// 验证文件是否存在
	if updates.FilePath != "" && updates.FilePath != table.FilePath {
		filePath := filepath.Join(project.RootPath, updates.FilePath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return c.Status(400).JSON(fiber.Map{
				"error": "配置文件不存在",
			})
		}
	}

	if err := app.DB.Model(table).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("更新配置表失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "luban_table", table.ID, fmt.Sprintf("更新配置表：%s", table.Name))

	return c.JSON(table)
}

// deleteTable 删除配置表
func deleteTable(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的配置表ID",
		})
	}

	table, err := getTableByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("配置表不存在: %v", err),
		})
	}

	if err := app.DB.Delete(table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除配置表失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "delete", "luban_table", table.ID, fmt.Sprintf("删除配置表：%s", table.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// exportConfig 导出配置
func exportConfig(c *fiber.Ctx) error {
	var export models.ConfigExport
	if err := c.BodyParser(&export); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的��求数据: %v", err),
		})
	}

	// 验证项目是否存在
	project, err := getProjectByID(export.ProjectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "项目不存在",
		})
	}

	// 验证配置表是否存在
	table, err := getTableByID(export.TableID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "配置表不存在",
		})
	}

	// 创建导出记录
	export.Status = "running"
	export.StartTime = time.Now()
	if err := app.DB.Create(&export).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建导出记录失败: %v", err),
		})
	}

	// 初始化进度信息
	progress := &ExportProgress{
		ID:        export.ID,
		ProjectID: export.ProjectID,
		TableID:   export.TableID,
		Status:    "running",
		StartTime: export.StartTime,
		Progress:  0,
	}
	exportProgressMap[export.ID] = progress

	// 异步执行导出任务
	go executeExport(&export, project, table)

	// 记录操作日志
	adminlog.CreateAdminLog(c, "export", "luban_export", export.ID, fmt.Sprintf("导出配置：%s", table.Name))

	return c.JSON(export)
}

// getExports 获取导出记录列表
func getExports(c *fiber.Ctx) error {
	projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 32)
	tableID, _ := strconv.ParseUint(c.Query("table_id"), 10, 32)
	exports, err := getExportList(uint(projectID), uint(tableID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取导出记录列表失败: %v", err),
		})
	}
	return c.JSON(exports)
}

// getExport 获取导出记录详情
func getExport(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的导出记录ID",
		})
	}

	export, err := getExportByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("导出记录不存在: %v", err),
		})
	}

	return c.JSON(export)
}

// getExportProgress 获取导出进度
func getExportProgress(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的导出记录ID",
		})
	}

	progress, exists := exportProgressMap[uint(id)]
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "找不到导出进度信息",
		})
	}

	return c.JSON(progress)
}
