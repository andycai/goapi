package gameconf

import (
	"github.com/andycai/goapi/internal"
	"github.com/gofiber/fiber/v2"
)

const ModulePriorityGameConf = 9908 // 游戏-游戏配置管理

var app *internal.App

type gameconfModule struct {
	internal.BaseModule
}

func init() {
	internal.RegisterModule(&gameconfModule{}, ModulePriorityGameConf)
}

func (m *gameconfModule) Awake(a *internal.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	return nil
}

func (m *gameconfModule) Start() error {
	// 初始化数据
	return initData()
}

func (m *gameconfModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/gameconf", app.HasPermission("gameconf:view"), func(c *fiber.Ctx) error {
		return c.Render("admin/gameconf", fiber.Map{
			"Title": "游戏配置管理",
			"Scripts": []string{
				"/static/js/admin/gameconf.js",
			},
		}, "admin/layout")
	})

	// api - projects
	app.RouterAdminApi.Get("/gameconf/projects", app.HasPermission("gameconf:view"), listProjectsHandler)           // 获取项目列表
	app.RouterAdminApi.Post("/gameconf/projects", app.HasPermission("gameconf:create"), createProjectHandler)       // 创建项目
	app.RouterAdminApi.Get("/gameconf/projects/:id", app.HasPermission("gameconf:view"), getProjectHandler)         // 获取项目详情
	app.RouterAdminApi.Put("/gameconf/projects/:id", app.HasPermission("gameconf:update"), updateProjectHandler)    // 更新项目
	app.RouterAdminApi.Delete("/gameconf/projects/:id", app.HasPermission("gameconf:delete"), deleteProjectHandler) // 删除项目

	// api - tables
	app.RouterAdminApi.Get("/gameconf/tables", app.HasPermission("gameconf:view"), listTablesHandler)                    // 获取配置表列表
	app.RouterAdminApi.Post("/gameconf/tables", app.HasPermission("gameconf:create"), createTableHandler)                // 创建配置表
	app.RouterAdminApi.Get("/gameconf/tables/:id", app.HasPermission("gameconf:view"), getTableHandler)                  // 获取配置表详情
	app.RouterAdminApi.Put("/gameconf/tables/:id", app.HasPermission("gameconf:update"), updateTableHandler)             // 更新配置表
	app.RouterAdminApi.Delete("/gameconf/tables/:id", app.HasPermission("gameconf:delete"), deleteTableHandler)          // 删除配置表
	app.RouterAdminApi.Post("/gameconf/tables/:id/validate", app.HasPermission("gameconf:update"), validateTableHandler) // 验证配置表

	// api - exports
	app.RouterAdminApi.Get("/gameconf/exports", app.HasPermission("gameconf:view"), listExportsHandler)                 // 获取导出记录列表
	app.RouterAdminApi.Post("/gameconf/exports", app.HasPermission("gameconf:create"), createExportHandler)             // 创建导出记录
	app.RouterAdminApi.Get("/gameconf/exports/:id", app.HasPermission("gameconf:view"), getExportHandler)               // 获取导出记录详情
	app.RouterAdminApi.Delete("/gameconf/exports/:id", app.HasPermission("gameconf:delete"), deleteExportHandler)       // 删除导出记录
	app.RouterAdminApi.Get("/gameconf/exports/:id/download", app.HasPermission("gameconf:view"), downloadExportHandler) // 下载导出文件

	return nil
}
