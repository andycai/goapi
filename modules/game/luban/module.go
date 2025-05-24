package luban

import (
	"github.com/andycai/goapi/core"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityLuban = 9909 // 游戏-游戏配置管理(Luban)

var app *core.App

type lubanModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&lubanModule{}, ModulePriorityLuban)
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
	app.RouterAdminApi.Get("/luban/projects", app.HasPermission("luban:view"), listProjectsHandler)           // 获取项目列表
	app.RouterAdminApi.Post("/luban/projects", app.HasPermission("luban:create"), createProjectHandler)       // 创建项目
	app.RouterAdminApi.Get("/luban/projects/:id", app.HasPermission("luban:view"), getProjectHandler)         // 获取项目详情
	app.RouterAdminApi.Put("/luban/projects/:id", app.HasPermission("luban:update"), updateProjectHandler)    // 更新项目
	app.RouterAdminApi.Delete("/luban/projects/:id", app.HasPermission("luban:delete"), deleteProjectHandler) // 删除项目

	app.RouterAdminApi.Get("/luban/tables", app.HasPermission("luban:view"), listTablesHandler)           // 获取配置表列表
	app.RouterAdminApi.Post("/luban/tables", app.HasPermission("luban:create"), createTableHandler)       // 创建配置表
	app.RouterAdminApi.Get("/luban/tables/:id", app.HasPermission("luban:view"), getTableHandler)         // 获取配置表详情
	app.RouterAdminApi.Put("/luban/tables/:id", app.HasPermission("luban:update"), updateTableHandler)    // 更新配置表
	app.RouterAdminApi.Delete("/luban/tables/:id", app.HasPermission("luban:delete"), deleteTableHandler) // 删除配置表

	app.RouterAdminApi.Post("/luban/export", app.HasPermission("luban:export"), exportConfigHandler)                 // 导出配置
	app.RouterAdminApi.Get("/luban/exports", app.HasPermission("luban:view"), getExportsHandler)                     // 获取导出记录列表
	app.RouterAdminApi.Get("/luban/exports/:id", app.HasPermission("luban:view"), getExportHandler)                  // 获取导出记录详情
	app.RouterAdminApi.Get("/luban/exports/progress/:id", app.HasPermission("luban:view"), getExportProgressHandler) // 获取导出进度

	return nil
}
