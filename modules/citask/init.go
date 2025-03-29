package citask

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type taskModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&taskModule{}, enum.ModulePriorityCiTask)
}

func (m *taskModule) Awake(a *core.App) error {
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
	app.RouterAdmin.Get("/citask", app.HasPermission("citask:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/citask", fiber.Map{
			"Title": "任务管理",
			"Scripts": []string{
				"/static/js/admin/citask.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/citask", app.HasPermission("citask:list"), listTasksHandler)                       // 获取任务列表
	app.RouterApi.Post("/citask", app.HasPermission("citask:create"), createTaskHandler)                   // 创建任务
	app.RouterApi.Get("/citask/running", app.HasPermission("citask:list"), listRunningTasksHandler)        // 获取正在执行的任务
	app.RouterApi.Get("/citask/next-run", app.HasPermission("citask:list"), getNextRunTimeHandler)         // 计算下次执行时间
	app.RouterApi.Get("/citask/search", app.HasPermission("citask:list"), searchTasksHandler)              // 添加搜索接口
	app.RouterApi.Get("/citask/:id", app.HasPermission("citask:list"), getTaskHandler)                     // 获取任务详情
	app.RouterApi.Put("/citask/:id", app.HasPermission("citask:update"), updateTaskHandler)                // 更新任务
	app.RouterApi.Delete("/citask/:id", app.HasPermission("citask:delete"), deleteTaskHandler)             // 删除任务
	app.RouterApi.Post("/citask/run/:id", app.HasPermission("citask:run"), runTaskHandler)                 // 执行任务
	app.RouterApi.Get("/citask/logs/:id", app.HasPermission("citask:list"), getTaskLogsHandler)            // 获取任务日志
	app.RouterApi.Get("/citask/progress/:logId", app.HasPermission("citask:list"), getTaskProgressHandler) // 获取任务进度
	app.RouterApi.Post("/citask/stop/:logId", app.HasPermission("citask:run"), stopTaskHandler)            // 停止任务

	return nil
}
