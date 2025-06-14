package citask

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityCiTask = 4002 // 功能-CI/CD 任务

var app *internal.App

type taskModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&taskModule{}, ModulePriorityCiTask)
}

func (m *taskModule) Awake(a *internal.App) error {
	app = a
	return autoMigrate()
}

func (m *taskModule) Start() error {
	if err := initData(); err != nil {
		return err
	}

	initCron()
	return nil
}

func (m *taskModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/citask", app.HasPermission("citask:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/citask", fiber.Map{
			"Title": "任务管理",
			"Scripts": []string{
				"/static/js/admin/citask.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterAdminApi.Get("/citask", app.HasPermission("citask:view"), listTasksHandler)                       // 获取任务列表
	app.RouterAdminApi.Post("/citask", app.HasPermission("citask:create"), createTaskHandler)                   // 创建任务
	app.RouterAdminApi.Get("/citask/running", app.HasPermission("citask:view"), listRunningTasksHandler)        // 获取正在执行的任务
	app.RouterAdminApi.Get("/citask/next-run", app.HasPermission("citask:view"), getNextRunTimeHandler)         // 计算下次执行时间
	app.RouterAdminApi.Get("/citask/search", app.HasPermission("citask:view"), searchTasksHandler)              // 添加搜索接口
	app.RouterAdminApi.Get("/citask/:id", app.HasPermission("citask:view"), getTaskHandler)                     // 获取任务详情
	app.RouterAdminApi.Put("/citask/:id", app.HasPermission("citask:update"), updateTaskHandler)                // 更新任务
	app.RouterAdminApi.Delete("/citask/:id", app.HasPermission("citask:delete"), deleteTaskHandler)             // 删除任务
	app.RouterAdminApi.Post("/citask/run/:id", app.HasPermission("citask:run"), runTaskHandler)                 // 执行任务
	app.RouterAdminApi.Get("/citask/logs/:id", app.HasPermission("citask:view"), getTaskLogsHandler)            // 获取任务日志
	app.RouterAdminApi.Get("/citask/progress/:logId", app.HasPermission("citask:view"), getTaskProgressHandler) // 获取任务进度
	app.RouterAdminApi.Post("/citask/stop/:logId", app.HasPermission("citask:run"), stopTaskHandler)            // 停止任务

	return nil
}
