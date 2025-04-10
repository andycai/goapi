package luban

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type lubanModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&lubanModule{}, enum.ModulePriorityLuban)
}

func (m *lubanModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	return autoMigrate()
}

func (m *lubanModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *lubanModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/luban", app.HasPermission("luban:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/luban", fiber.Map{
			"Title": "游戏配置表管理",
			"Scripts": []string{
				"/static/js/admin/luban.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/luban/projects", app.HasPermission("luban:view"), listProjectsHandler)           // 获取项目列表
	app.RouterApi.Post("/luban/projects", app.HasPermission("luban:create"), createProjectHandler)       // 创建项目
	app.RouterApi.Get("/luban/projects/:id", app.HasPermission("luban:view"), getProjectHandler)         // 获取项目详情
	app.RouterApi.Put("/luban/projects/:id", app.HasPermission("luban:update"), updateProjectHandler)    // 更新项目
	app.RouterApi.Delete("/luban/projects/:id", app.HasPermission("luban:delete"), deleteProjectHandler) // 删除项目

	app.RouterApi.Get("/luban/tables", app.HasPermission("luban:view"), listTablesHandler)           // 获取配置表列表
	app.RouterApi.Post("/luban/tables", app.HasPermission("luban:create"), createTableHandler)       // 创建配置表
	app.RouterApi.Get("/luban/tables/:id", app.HasPermission("luban:view"), getTableHandler)         // 获取配置表详情
	app.RouterApi.Put("/luban/tables/:id", app.HasPermission("luban:update"), updateTableHandler)    // 更新配置表
	app.RouterApi.Delete("/luban/tables/:id", app.HasPermission("luban:delete"), deleteTableHandler) // 删除配置表

	app.RouterApi.Post("/luban/export", app.HasPermission("luban:export"), exportConfigHandler)                 // 导出配置
	app.RouterApi.Get("/luban/exports", app.HasPermission("luban:view"), getExportsHandler)                     // 获取导出记录列表
	app.RouterApi.Get("/luban/exports/:id", app.HasPermission("luban:view"), getExportHandler)                  // 获取导出记录详情
	app.RouterApi.Get("/luban/exports/progress/:id", app.HasPermission("luban:view"), getExportProgressHandler) // 获取导出进度

	return nil
}
