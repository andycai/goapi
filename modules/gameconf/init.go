package gameconf

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type gameconfModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&gameconfModule{}, enum.ModulePriorityGameConf)
}

func (m *gameconfModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	// 初始化数据
	return initData()
}

func (m *gameconfModule) Start() error {
	return nil
}

func (m *gameconfModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/gameconf", app.HasPermission("gameconf:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/gameconf", fiber.Map{
			"Title": "游戏配置管理",
			"Scripts": []string{
				"/static/js/admin/gameconf.js",
			},
		}, "admin/layout")
	})

	// api - projects
	app.RouterApi.Get("/gameconf/projects", app.HasPermission("gameconf:list"), getProjects)            // 获取项目列表
	app.RouterApi.Post("/gameconf/projects", app.HasPermission("gameconf:create"), createProject)       // 创建项目
	app.RouterApi.Get("/gameconf/projects/:id", app.HasPermission("gameconf:list"), getProject)         // 获取项目详情
	app.RouterApi.Put("/gameconf/projects/:id", app.HasPermission("gameconf:update"), updateProject)    // 更新项目
	app.RouterApi.Delete("/gameconf/projects/:id", app.HasPermission("gameconf:delete"), deleteProject) // 删除项目

	// api - tables
	app.RouterApi.Get("/gameconf/tables", app.HasPermission("gameconf:list"), getTables)                     // 获取配置表列表
	app.RouterApi.Post("/gameconf/tables", app.HasPermission("gameconf:create"), createTable)                // 创建配置表
	app.RouterApi.Get("/gameconf/tables/:id", app.HasPermission("gameconf:list"), getTable)                  // 获取配置表详情
	app.RouterApi.Put("/gameconf/tables/:id", app.HasPermission("gameconf:update"), updateTable)             // 更新配置表
	app.RouterApi.Delete("/gameconf/tables/:id", app.HasPermission("gameconf:delete"), deleteTable)          // 删除配置表
	app.RouterApi.Post("/gameconf/tables/:id/validate", app.HasPermission("gameconf:update"), validateTable) // 验证配置表

	// api - exports
	app.RouterApi.Get("/gameconf/exports", app.HasPermission("gameconf:list"), getExports)                  // 获取导出记录列表
	app.RouterApi.Post("/gameconf/exports", app.HasPermission("gameconf:create"), createExport)             // 创建导出记录
	app.RouterApi.Get("/gameconf/exports/:id", app.HasPermission("gameconf:list"), getExport)               // 获取导出记录详情
	app.RouterApi.Delete("/gameconf/exports/:id", app.HasPermission("gameconf:delete"), deleteExport)       // 删除导出记录
	app.RouterApi.Get("/gameconf/exports/:id/download", app.HasPermission("gameconf:list"), downloadExport) // 下载导出文件

	return nil
}
